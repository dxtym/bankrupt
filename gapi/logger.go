package gapi

import (
	"context"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func GrpcLogger(
	ctx context.Context,
	req any,
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (resp any, err error) {
	start := time.Now()
	resp, err = handler(ctx, req)
	logger := log.Info()
	if err != nil {
		logger = log.Error().Err(err)
	}
	duration := time.Since(start)

	statusCode := codes.Unknown
	if status, ok := status.FromError(err); ok {
		statusCode = status.Code()
	}

	logger.Str("protocol", "grpc").
		Str("method", info.FullMethod).
		Int("status_code", int(statusCode)).
		Str("status_text", statusCode.String()).
		Dur("duration", duration).
		Msg("received a gRPC request")
	return
}

type ResponseRecorder struct {
	http.ResponseWriter
	statusCode int
	body       []byte
}

// rewrite response writer to capture status code and response body
func (r *ResponseRecorder) WriteHeader(statusCode int) {
	r.statusCode = statusCode
	r.ResponseWriter.WriteHeader(statusCode)
}

func (r *ResponseRecorder) Write(b []byte) (int, error) {
	r.body = b
	return r.ResponseWriter.Write(b)
}

func HttpLogger(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		start := time.Now()
		rec := &ResponseRecorder{
			ResponseWriter: res,
			statusCode:     http.StatusOK,
		}
		handler.ServeHTTP(rec, req)
		duration := time.Since(start)

		logger := log.Info()
		if rec.statusCode != http.StatusOK {
			logger = log.Error().Bytes("body", rec.body)
		}

		logger.Str("protocol", "http").
			Str("method", req.Method).
			Str("path", req.RequestURI).
			Dur("duration", duration).
			Msg("received a HTTP request")
	})
}
