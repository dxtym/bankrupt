package gapi

import (
	"context"
	"time"

	db "github.com/dxtym/bankrupt/db/sqlc"
	"github.com/dxtym/bankrupt/pb"
	"github.com/dxtym/bankrupt/utils"
	"github.com/dxtym/bankrupt/valid"
	"github.com/dxtym/bankrupt/worker"
	"github.com/hibiken/asynq"
	"github.com/lib/pq"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	violations := validateCreateUserRequest(req)
	if len(violations) > 0 {
		return nil, invalidArgumentError(violations)
	}
	hashedPassword, err := utils.HashPassword(req.GetPassword())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "user not authorized: %v", err)
	}

	// run redis queue and create user in one transaction
	arg := db.CreateUserTxParams{
		CreateUserParams: db.CreateUserParams{
			Username:       req.GetUsername(),
			HashedPassword: hashedPassword,
			FullName:       req.GetFullName(),
			Email:          req.GetEmail(),
		},
		AfterCreate: func(user db.User) error {
			taskPayload := &worker.PayloadSendEmail{
				Username: user.Username,
			}
			opts := []asynq.Option{
				asynq.MaxRetry(10),
				asynq.Timeout(10 * time.Second),
				asynq.Queue(worker.QueueCritical),
			}
			return s.taskDistributor.DistributorTaskSendEmail(ctx, *taskPayload, opts...)
		},
	}

	txResult, err := s.store.CreateUserTx(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				return nil, status.Errorf(codes.PermissionDenied, "user already exists: %v", err)
			}
		}
		return nil, status.Errorf(codes.Internal, "cannot create the user: %v", err)
	}

	res := &pb.CreateUserResponse{
		User: convertUser(txResult.User),
	}
	return res, nil
}

func validateCreateUserRequest(req *pb.CreateUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := valid.ValidateUsername(req.GetUsername()); err != nil {
		violations = append(violations, fieldViolation("username", err))
	}
	if err := valid.ValidateFullname(req.GetFullName()); err != nil {
		violations = append(violations, fieldViolation("full_name", err))
	}
	if err := valid.ValidateEmail(req.GetEmail()); err != nil {
		violations = append(violations, fieldViolation("email", err))
	}
	if err := valid.ValidatePassword(req.GetPassword()); err != nil {
		violations = append(violations, fieldViolation("password", err))
	}
	return
}
