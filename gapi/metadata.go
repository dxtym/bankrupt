package gapi

import (
	"context"

	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

const (
	grpcUserAgentHeader = "grpcgateway-user-agent"
	userAgentHeader     = "user-agent"
	xForwardForHeader   = "x-forwarded-for"
)

type Metadata struct {
	userAgent string
	clientIP  string
}

func (s *Server) GetMetadata(ctx context.Context) *Metadata {
	meta := &Metadata{}

	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if userAgents := md.Get(grpcUserAgentHeader); len(userAgents) > 0 {
			meta.userAgent = userAgents[0]
		}
		if userAgents := md.Get(userAgentHeader); len(userAgents) > 0 {
			meta.userAgent = userAgents[0]
		}
		if clientIPs := md.Get(xForwardForHeader); len(clientIPs) > 0 {
			meta.clientIP = clientIPs[0]
		}
	}

	if p, ok := peer.FromContext(ctx); ok {
		meta.clientIP = p.Addr.String()
	}

	return meta
}
