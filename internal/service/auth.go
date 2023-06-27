package service

import (
	"context"
	"strings"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/HardDie/mmr_boost_server/internal/config"
	"github.com/HardDie/mmr_boost_server/internal/dto"
	"github.com/HardDie/mmr_boost_server/internal/entity"
	"github.com/HardDie/mmr_boost_server/internal/logger"
	"github.com/HardDie/mmr_boost_server/internal/repository/postgres"
	"github.com/HardDie/mmr_boost_server/internal/repository/smtp"
	"github.com/HardDie/mmr_boost_server/internal/utils"
)

type Auth struct {
	repository     *postgres.Postgres
	smtpRepository *smtp.SMTP

	config config.Config
}

func NewAuth(config config.Config, repository *postgres.Postgres, smtp *smtp.SMTP) *Auth {
	return &Auth{
		config:         config,
		repository:     repository,
		smtpRepository: smtp,
	}
}

func (s *Auth) Register(ctx context.Context, req *dto.AuthRegisterRequest) error {
	var res *entity.User

	err := s.repository.TxManager().ReadWriteTx(ctx, func(ctx context.Context) error {
		// Check if username is not busy
		user, err := s.repository.User.GetByNameOrEmail(ctx, req.Username, req.Email)
		if err != nil {
			logger.Error.Printf("error while trying get user: %v", err.Error())
			return status.Error(codes.Internal, "internal")
		}
		if user != nil {
			return status.Error(codes.InvalidArgument, "username already exist or email is busy")
		}

		// Hashing password
		hashPassword, err := utils.HashBcrypt(req.Password)
		if err != nil {
			logger.Error.Printf("error hash bcrypt: %v", err.Error())
			return status.Error(codes.Internal, "internal")
		}

		// Create a user
		user, err = s.repository.User.Create(ctx, req.Email, req.Username)
		if err != nil {
			logger.Error.Printf("error writing user into DB: %v", err.Error())
			return status.Error(codes.Internal, "internal")
		}

		// Create a password
		_, err = s.repository.Password.Create(ctx, user.ID, hashPassword)
		if err != nil {
			logger.Error.Printf("error writing password into DB: %v", err.Error())
			return status.Error(codes.Internal, "internal")
		}

		res = user
		return nil
	})
	if err != nil {
		return err
	}

	// Create record in history
	err = s.repository.History.NewEvent(ctx, res.ID, "user was created")
	if err != nil {
		logger.Error.Println("error writing history message: user was created")
	}

	err = s.generateAndSendEmail(ctx, res.ID, res.Email)
	if err != nil {
		return err
	}
	return nil
}
func (s *Auth) Login(ctx context.Context, req *dto.AuthLoginRequest) (*entity.User, error) {
	var res *entity.User

	err := s.repository.TxManager().ReadWriteTx(ctx, func(ctx context.Context) error {
		// Check if such user exist
		user, err := s.repository.User.GetByName(ctx, req.Username)
		if err != nil {
			logger.Error.Printf("error while trying get user: %v", err.Error())
			return status.Error(codes.Internal, "internal")
		}
		if user == nil {
			return status.Error(codes.InvalidArgument, "username or password is invalid")
		}

		if !user.IsActivated {
			return status.Error(codes.InvalidArgument, "account is not activated")
		}

		// Get password from DB
		password, err := s.repository.Password.GetByUserID(ctx, user.ID)
		if err != nil {
			logger.Error.Printf("error while trying get password: %v", err.Error())
			return status.Error(codes.Internal, "internal")
		}
		if password == nil {
			logger.Error.Printf("password for user %d not found", user.ID)
			return status.Error(codes.Internal, "internal")
		}

		// Check if the password is locked after failed attempts
		if password.FailedAttempts >= int32(s.config.Password.MaxAttempts) {
			// Check if the password block time has expired
			if time.Since(password.UpdatedAt) <= time.Hour*time.Duration(s.config.Password.BlockTime) {
				return status.Error(codes.PermissionDenied, "too many invalid requests")
			}
			// If the blocking time has expired, reset the counter of failed attempts
			password, err = s.repository.Password.ResetFailedAttempts(ctx, password.ID)
			if err != nil {
				logger.Error.Printf("error resetting the counter of failed attempts: %v", err)
				return status.Error(codes.Internal, "internal")
			}
		}

		// Check if password is correct
		if !utils.HashBcryptCompare(req.Password, password.PasswordHash) {
			// Increased number of failed attempts
			_, err = s.repository.Password.IncreaseFailedAttempts(ctx, password.ID)
			if err != nil {
				logger.Error.Printf("Error increasing failed attempts: %v", err.Error())
			}
			return status.Error(codes.InvalidArgument, "username or password is invalid")
		}

		// Reset the failed attempts counter after the first successful attempt
		if password.FailedAttempts > 0 {
			_, err = s.repository.Password.ResetFailedAttempts(ctx, password.ID)
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
func (s *Auth) Logout(ctx context.Context, sessionID int32) error {
	err := s.repository.AccessToken.DeleteByID(ctx, sessionID)
	if err != nil {
		logger.Error.Printf("error deleting session: %v", err.Error())
		return status.Error(codes.Internal, "internal")
	}
	return nil
}
func (s *Auth) GenerateCookie(ctx context.Context, userID int32) (*entity.AccessToken, error) {
	// Generate session key
	sessionKey, err := utils.GenerateSessionKey()
	if err != nil {
		logger.Error.Printf("error generate session key: %v", err)
		return nil, status.Error(codes.Internal, "internal")
	}

	// Calculate expired at
	expiredAt := time.Now().Add(time.Minute * time.Duration(s.config.Session.AccessToken))

	// Write session to DB
	resp, err := s.repository.AccessToken.CreateOrUpdate(ctx, userID, utils.HashSha256(sessionKey), expiredAt)
	if err != nil {
		logger.Error.Printf("write access token to DB: %v", err)
		return nil, status.Error(codes.Internal, "internal")
	}
	resp.TokenHash = sessionKey

	return resp, nil
}
func (s *Auth) ValidateCookie(ctx context.Context, sessionKey string) (*entity.User, *entity.AccessToken, error) {
	// Check if access token exist
	tokenHash := utils.HashSha256(sessionKey)
	accessToken, err := s.repository.AccessToken.GetByUserID(ctx, tokenHash)
	if err != nil {
		logger.Error.Printf("error read access token from db: %v", err.Error())
		return nil, nil, status.Error(codes.Internal, "internal")
	}
	if accessToken == nil {
		return nil, nil, status.Error(codes.PermissionDenied, "access token expired or not exist")
	}

	// Check if session is not expired
	if time.Now().After(accessToken.ExpiredAt) {
		return nil, nil, status.Error(codes.PermissionDenied, "access token expired or not exist")
	}

	user, err := s.repository.User.GetByID(ctx, accessToken.UserID)
	if err != nil {
		logger.Error.Println("can't found user from access token")
		return nil, nil, status.Error(codes.Internal, "internal")
	}

	return user, accessToken, nil
}
func (s *Auth) GetUserInfo(ctx context.Context, userID int32) (*entity.User, error) {
	user, err := s.repository.User.GetByID(ctx, userID)
	if err != nil {
		logger.Error.Printf("error get user info: %v", err.Error())
		return nil, status.Error(codes.Internal, "internal")
	}
	return user, nil
}
func (s *Auth) ValidateEmail(ctx context.Context, code string) error {
	codeHash := utils.HashSha256(code)
	emailValidation, err := s.repository.EmailValidation.GetByCode(ctx, codeHash)
	if err != nil {
		logger.Error.Printf("error finding email validation record: %v", err.Error())
		return status.Error(codes.Internal, "internal")
	}

	// Check if validation code exist
	if emailValidation == nil {
		return status.Error(codes.InvalidArgument, "validation code expired, invalid or not exist")
	}

	// Check if validation code expired
	if time.Now().After(emailValidation.ExpiredAt) {
		err = s.repository.EmailValidation.DeleteByID(ctx, emailValidation.ID)
		if err != nil {
			logger.Error.Printf("error deleting email validation expired record: %v", err.Error())
		}
		return status.Error(codes.InvalidArgument, "validation code expired, invalid or not exist")
	}

	// Activate user
	_, err = s.repository.User.ActivateRecord(ctx, emailValidation.UserID)
	if err != nil {
		logger.Error.Printf("error activating user with email code: %v", err.Error())
		return status.Error(codes.Internal, "internal")
	}

	// Delete activation code from DB
	err = s.repository.EmailValidation.DeleteByID(ctx, emailValidation.ID)
	if err != nil {
		logger.Error.Printf("error deleting email validation record after validating: %v", err.Error())
	}

	// Write history record
	err = s.repository.History.NewEvent(ctx, emailValidation.UserID, "account was activated")
	if err != nil {
		logger.Error.Println("error writing history message: account was activated")
	}

	return nil
}
func (s *Auth) SendValidationEmail(ctx context.Context, name string) error {
	u, err := s.repository.User.GetByName(ctx, name)
	if err != nil {
		logger.Error.Println("error get user by name:", err.Error())
		return status.Error(codes.Internal, "internal")
	}
	if u == nil || u.IsActivated {
		return nil
	}

	err = s.generateAndSendEmail(ctx, u.ID, u.Email)
	if err != nil {
		return err
	}
	return nil
}

func (s *Auth) generateAndSendEmail(ctx context.Context, userID int32, email string) error {
	// Generate new email validation code
	emailCode, err := utils.UUIDGenerate()
	if err != nil {
		logger.Error.Println("error generating email code:", err.Error())
		return status.Error(codes.Internal, "internal")
	}
	emailCode = strings.ToLower(emailCode)
	codeHash := utils.HashSha256(emailCode)

	// Calculate expired at
	expiredAt := time.Now().Add(time.Hour * time.Duration(s.config.EmailValidation.Expiration))

	// Create record of email validation in DB
	_, err = s.repository.EmailValidation.CreateOrUpdate(ctx, userID, codeHash, expiredAt)
	if err != nil {
		logger.Error.Println("error writing email validation token to DB:", err.Error())
		return status.Error(codes.Internal, "internal")
	}

	// Send code to email
	err = s.smtpRepository.SendEmailVerification(email, emailCode)
	if err != nil {
		logger.Error.Println("error sending email with verification code:", err.Error())
		return status.Error(codes.Internal, "internal")
	}
	return nil
}
