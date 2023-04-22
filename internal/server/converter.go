package server

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/HardDie/mmr_boost_server/internal/entity"
	"github.com/HardDie/mmr_boost_server/internal/utils"
	pb "github.com/HardDie/mmr_boost_server/pkg/proto/server"
)

func UserToPb(u *entity.User) *pb.UserObject {
	return &pb.UserObject{
		Id:          u.ID,
		Email:       u.Email,
		Username:    u.Username,
		RoleId:      u.RoleID,
		SteamId:     u.SteamID,
		IsActivated: u.IsActivated,
		CreatedAt:   timestamppb.New(u.CreatedAt),
		UpdatedAt:   timestamppb.New(u.UpdatedAt),
		DeletedAt:   utils.TimetamppbFromTime(u.DeletedAt),
	}
}

func ApplicationPublicToPb(a *entity.ApplicationPublic) *pb.PublicApplicationObject {
	return &pb.PublicApplicationObject{
		Id:         a.ID,
		UserId:     a.UserID,
		StatusId:   a.StatusID,
		TypeId:     a.TypeID,
		CurrentMmr: a.CurrentMMR,
		TargetMmr:  a.TargetMMR,
		TgContact:  a.TgContact,
		CreatedAt:  timestamppb.New(a.CreatedAt),
		UpdatedAt:  timestamppb.New(a.UpdatedAt),
		DeletedAt:  utils.TimetamppbFromTime(a.DeletedAt),
	}
}