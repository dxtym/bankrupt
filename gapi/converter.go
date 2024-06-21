package gapi

import (
	db "github.com/dxtym/bankrupt/db/sqlc"
	"github.com/dxtym/bankrupt/pb"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func convertUser(user db.User) *pb.User {
	return &pb.User{
		Username:          user.Username,
		FullName:          user.FullName,
		Email:             user.Email,
		PasswordChangedAt: timestamppb.New(user.PasswordChangedAt),
		CreatedAt:         timestamppb.New(user.CreatedAt),
	}
}

func convertUUID(uuid uuid.UUID) *pb.UUID {
	return &pb.UUID{
		Value: uuid.String(),
	}
}
