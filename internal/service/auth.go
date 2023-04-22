package service

import (
	"context"
	"strings"
	"time"

	"github.com/HardDie/mmr_boost_server/internal/config"
	"github.com/HardDie/mmr_boost_server/internal/dto"
	"github.com/HardDie/mmr_boost_server/internal/entity"
	"github.com/HardDie/mmr_boost_server/internal/errs"
	"github.com/HardDie/mmr_boost_server/internal/logger"
	"github.com/HardDie/mmr_boost_server/internal/repository/postgres"
	"github.com/HardDie/mmr_boost_server/internal/repository/smtp"
	"github.com/HardDie/mmr_boost_server/internal/utils"
)

type auth struct {
	repository     *postgres.Postgres
	smtpRepository *smtp.SMTP

	config config.Config
}

func newAuth(config config.Config, repository *postgres.Postgres, smtp *smtp.SMTP) auth {
	return auth{
		config:         config,
		repository:     repository,
		smtpRepository: smtp,
	}
}

func (s *auth) AuthRegister(ctx context.Context, req *dto.AuthRegisterRequest) error {
	var res *entity.User

	err := s.repository.TxManager().ReadWriteTx(ctx, func(ctx context.Context) error {
		// Check if username is not busy
		user, err := s.repository.UserGetByName(ctx, req.Username)
		if err != nil {
			logger.Error.Printf("error while trying get user: %v", err.Error())
			return errs.InternalError
		}
		if user != nil {
			return errs.BadRequest.AddMessage("username already exist")
		}

		// Hashing password
		hashPassword, err := utils.HashBcrypt(req.Password)
		if err != nil {
			logger.Error.Printf("error hash bcrypt: %v", err.Error())
			return errs.InternalError
		}

		// Create a user
		user, err = s.repository.UserCreate(ctx, req.Email, req.Username)
		if err != nil {
			logger.Error.Printf("error writing user into DB: %v", err.Error())
			return errs.InternalError
		}

		// Create a password
		_, err = s.repository.PasswordCreate(ctx, user.ID, hashPassword)
		if err != nil {
			logger.Error.Printf("error writing password into DB: %v", err.Error())
			return errs.InternalError
		}

		res = user
		return nil
	})
	if err != nil {
		return err
	}

	// Create record in history
	err = s.repository.HistoryNewEvent(ctx, res.ID, "user was created")
	if err != nil {
		logger.Error.Println("error writing history message: user was created")
	}

	// Generate new email validation code
	emailCode, err := utils.UUIDGenerate()
	if err != nil {
		logger.Error.Println("error generating email code:", err.Error())
		return errs.InternalError
	}
	emailCode = strings.ToLower(emailCode)
	codeHash := utils.HashSha256(emailCode)

	// Calculate expired at
	expiredAt := time.Now().Add(time.Hour * time.Duration(s.config.EmailValidation.Expiration))

	// Create record of email validation in DB
	_, err = s.repository.EmailValidationCreateOrUpdate(ctx, res.ID, codeHash, expiredAt)
	if err != nil {
		logger.Error.Println("error writing email validation token to DB:", err.Error())
		return errs.InternalError
	}

	// Send code to email
	err = s.smtpRepository.SendEmailVerification(res.Email, emailCode)
	if err != nil {
		logger.Error.Println("error sending email with verification code:", err.Error())
	}

	return nil
}
func (s *auth) AuthLogin(ctx context.Context, req *dto.AuthLoginRequest) (*entity.User, error) {
	var res *entity.User

	err := s.repository.TxManager().ReadWriteTx(ctx, func(ctx context.Context) error {
		// Check if such user exist
		user, err := s.repository.UserGetByName(ctx, req.Username)
		if err != nil {
			logger.Error.Printf("error while trying get user: %v", err.Error())
			return errs.InternalError
		}
		if user == nil {
			return errs.BadRequest.AddMessage("username or password is invalid")
		}

		if !user.IsActivated {
			return errs.BadRequest.AddMessage("account is not activated")
		}

		// Get password from DB
		password, err := s.repository.PasswordGetByUserID(ctx, user.ID)
		if err != nil {
			logger.Error.Printf("error while trying get password: %v", err.Error())
			return errs.InternalError
		}
		if password == nil {
			logger.Error.Printf("password for user %d not found", user.ID)
			return errs.InternalError
		}

		// Check if the password is locked after failed attempts
		if password.FailedAttempts >= int32(s.config.Password.MaxAttempts) {
			// Check if the password block time has expired
			if time.Now().Sub(password.UpdatedAt) <= time.Hour*time.Duration(s.config.Password.BlockTime) {
				return errs.UserBlocked.AddMessage("too many invalid requests")
			}
			// If the blocking time has expired, reset the counter of failed attempts
			password, err = s.repository.PasswordResetFailedAttempts(ctx, password.ID)
			if err != nil {
				logger.Error.Printf("error resetting the counter of failed attempts: %v", err)
				return errs.InternalError
			}
		}

		// Check if password is correct
		if !utils.HashBcryptCompare(req.Password, password.PasswordHash) {
			// Increased number of failed attempts
			_, err = s.repository.PasswordIncreaseFailedAttempts(ctx, password.ID)
			if err != nil {
				logger.Error.Printf("Error increasing failed attempts: %v", err.Error())
			}
			return errs.BadRequest.AddMessage("username or password is invalid")
		}

		// Reset the failed attempts counter after the first successful attempt
		if password.FailedAttempts > 0 {
			_, err = s.repository.PasswordResetFailedAttempts(ctx, password.ID)
			if err != nil {
				logger.Error.Printf("Error flushing failed attempts: %v", err.Error())
			}
		}

		res = user
		return nil
	})
	if err != nil {
		return nil, err
	}

	return res, nil
}
func (s *auth) AuthLogout(ctx context.Context, sessionID int32) error {
	err := s.repository.AccessTokenDeleteByID(ctx, sessionID)
	if err != nil {
		logger.Error.Printf("error deleting session: %v", err.Error())
		return errs.InternalError
	}
	return nil
}
func (s *auth) AuthGenerateCookie(ctx context.Context, userID int32) (*entity.AccessToken, error) {
	// Generate session key
	sessionKey, err := utils.GenerateSessionKey()
	if err != nil {
		logger.Error.Printf("error generate session key: %v", err)
		return nil, errs.InternalError
	}

	// Calculate expired at
	expiredAt := time.Now().Add(time.Minute * time.Duration(s.config.Session.AccessToken))

	// Write session to DB
	resp, err := s.repository.AccessTokenCreateOrUpdate(ctx, userID, utils.HashSha256(sessionKey), expiredAt)
	if err != nil {
		logger.Error.Printf("write access token to DB: %v", err)
		return nil, errs.InternalError
	}
	resp.TokenHash = sessionKey

	return resp, nil
}
func (s *auth) AuthValidateCookie(ctx context.Context, sessionKey string) (*entity.AccessToken, error) {
	// Check if access token exist
	tokenHash := utils.HashSha256(sessionKey)
	accessToken, err := s.repository.AccessTokenGetByUserID(ctx, tokenHash)
	if err != nil {
		logger.Error.Printf("error read access token from db: %v", err.Error())
		return nil, errs.InternalError
	}
	if accessToken == nil {
		return nil, errs.SessionInvalid.AddMessage("access token not exist")
	}

	// Check if session is not expired
	if time.Now().After(accessToken.ExpiredAt) {
		return nil, errs.SessionInvalid.AddMessage("access token has expired")
	}
	return accessToken, nil
}
func (s *auth) AuthGetUserInfo(ctx context.Context, userID int32) (*entity.User, error) {
	user, err := s.repository.UserGetByID(ctx, userID)
	if err != nil {
		logger.Error.Printf("error get user info: %v", err.Error())
		return nil, errs.InternalError
	}
	return user, nil
}
func (s *auth) AuthValidateEmail(ctx context.Context, code string) error {
	codeHash := utils.HashSha256(code)
	emailValidation, err := s.repository.EmailValidationGetByCode(ctx, codeHash)
	if err != nil {
		logger.Error.Printf("error finding email validation record: %v", err.Error())
		return errs.InternalError
	}

	// Check if validation code exist
	if emailValidation == nil {
		return errs.EmailValidationCodeNotExist
	}

	// Check if validation code expired
	if time.Now().After(emailValidation.ExpiredAt) {
		err = s.repository.EmailValidationDeleteByID(ctx, emailValidation.ID)
		if err != nil {
			logger.Error.Printf("error deleting email validation expired record: %v", err.Error())
		}
		return errs.EmailValidationCodeExpired
	}

	// Activate user
	_, err = s.repository.UserActivateRecord(ctx, emailValidation.UserID)
	if err != nil {
		logger.Error.Printf("error activating user with email code: %v", err.Error())
		return errs.InternalError
	}

	// Delete activation code from DB
	err = s.repository.EmailValidationDeleteByID(ctx, emailValidation.ID)
	if err != nil {
		logger.Error.Printf("error deleting email validation record after validating: %v", err.Error())
	}

	// Write history record
	err = s.repository.HistoryNewEvent(ctx, emailValidation.UserID, "account was activated")
	if err != nil {
		logger.Error.Println("error writing history message: account was activated")
	}

	return nil
}
func (s *auth) AuthSendValidationEmail(ctx context.Context, name string) error {
	u, err := s.repository.UserGetByName(ctx, name)
	if err != nil {
		logger.Error.Println("error get user by name:", err.Error())
		return errs.InternalError
	}
	if u == nil || u.IsActivated {
		return nil
	}

	// Generate new email validation code
	emailCode, err := utils.UUIDGenerate()
	if err != nil {
		logger.Error.Println("error generating email code:", err.Error())
		return errs.InternalError
	}
	emailCode = strings.ToLower(emailCode)
	codeHash := utils.HashSha256(emailCode)

	// Calculate expired at
	expiredAt := time.Now().Add(time.Hour * time.Duration(s.config.EmailValidation.Expiration))

	// Create record of email validation in DB
	_, err = s.repository.EmailValidationCreateOrUpdate(ctx, u.ID, codeHash, expiredAt)
	if err != nil {
		logger.Error.Println("error writing email validation token to DB:", err.Error())
		return errs.InternalError
	}

	// Send code to email
	err = s.smtpRepository.SendEmailVerification(u.Email, emailCode)
	if err != nil {
		logger.Error.Println("error sending email with verification code:", err.Error())
		return errs.InternalError
	}

	return nil
}
