package gapi

import (
	"context"
	"database/sql"

	db "github.com/dxtym/bankrupt/db/sqlc"
	"github.com/dxtym/bankrupt/pb"
	"github.com/dxtym/bankrupt/utils"
	"github.com/dxtym/bankrupt/valid"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *Server) LoginUser(ctx context.Context, req *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {
	violations := validateLoginUserRequest(req)
	if len(violations) > 0 {
		return nil, invalidArgumentError(violations)
	}

	user, err := s.store.GetUser(ctx, req.GetUsername())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "no such user: %v", err)
		}

		return nil, status.Errorf(codes.Internal, "cannot get user: %v", err)
	}

	err = utils.CheckPassword(req.GetPassword(), user.HashedPassword)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "password is incorrect: %v", err)
	}

	accessToken, accessPayload, err := s.token.CreateToken(user.Username, s.config.TokenDuration)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot create access token: %v", err)
	}

	refreshToken, refreshPayload, err := s.token.CreateToken(user.Username, s.config.RefreshDuration)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot create refresh token: %v", err)
	}

	meta := s.GetMetadata(ctx)

	session, err := s.store.CreateSession(ctx, db.CreateSessionParams{
		ID:           refreshPayload.Id,
		Username:     user.Username,
		RefreshToken: refreshToken,
		UserAgent:    meta.userAgent,
		ClientIp:     meta.clientIP,
		IsBlocked:    false,
		ExpiresAt:    refreshPayload.ExpiredAt,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot create session: %v", err)
	}

	res := &pb.LoginUserResponse{
		SessionId:             convertUUID(session.ID),
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  timestamppb.New(accessPayload.ExpiredAt),
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: timestamppb.New(refreshPayload.ExpiredAt),
		User:                  convertUser(user),
	}
	return res, nil
}

func validateLoginUserRequest(req *pb.LoginUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := valid.ValidateUsername(req.GetUsername()); err != nil {
		violations = append(violations, fieldViolation("username", err))
	}
	if err := valid.ValidatePassword(req.GetPassword()); err != nil {
		violations = append(violations, fieldViolation("password", err))
	}
	return
}
