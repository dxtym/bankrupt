package gapi

import (
	"context"
	"fmt"
	"strings"

	"github.com/dxtym/bankrupt/token"
	"google.golang.org/grpc/metadata"
)

const (
	authHeader = "authorization"
	authType   = "bearer"
)

func (s *Server) authorizeUser(ctx context.Context) (*token.Payload, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, fmt.Errorf("metadata is not provided")
	}
	authValues := md.Get(authHeader)
	if len(authValues) == 0 {
		return nil, fmt.Errorf("authorization token is not provided")
	}

	authFields := strings.Fields(authValues[0])
	if len(authFields) < 2 {
		return nil, fmt.Errorf("authorization token is not provided")
	}

	authFieldType := strings.ToLower(authFields[0])
	if authFieldType != authType {
		return nil, fmt.Errorf("unsupported authorization type")
	}

	accessToken := authFields[1]
	payload, err := s.token.VerifyToken(accessToken)
	if err != nil {
		return nil, fmt.Errorf("invalid access token")
	}

	return payload, nil
}
