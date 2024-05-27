package api

import (
	"database/sql"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type RenewTokenRequst struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type RenewTokenResponse struct {
	AccessToken          string    `json:"access_token"`
	AccessTokenExpiresAt time.Time `json:"access_token_expires_at"`
}

func (s *Server) RenewToken(ctx *gin.Context) {
	var req RenewTokenRequst
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	payload, err := s.token.VerifyToken(req.RefreshToken)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	session, err := s.store.GetSession(ctx, payload.Id)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// check session is blocked
	if session.IsBlocked {
		err := errors.New("session is blocked")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	// check username
	if session.Username != payload.Username {
		err := errors.New("username is not match")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	// check refresh token
	if session.RefreshToken != req.RefreshToken {
		err := errors.New("refresh token is not valid")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	// check refresh token is expired
	if time.Now().After(session.ExpiresAt) {
		err := errors.New("refresh token is expired")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	// create token for user
	accessToken, accessPayload, err := s.token.CreateToken(payload.Username, s.config.TokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// send token and user response
	res := RenewTokenResponse{
		AccessToken:          accessToken,
		AccessTokenExpiresAt: accessPayload.ExpiredAt,
	}
	ctx.JSON(http.StatusOK, res)
}
