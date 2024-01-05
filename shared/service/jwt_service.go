package service

import (
	"fmt"
	"time"

	"catatan-keuangan-app/config"
	"catatan-keuangan-app/entity"
	"catatan-keuangan-app/entity/dto"
	"catatan-keuangan-app/shared/model"

	"github.com/golang-jwt/jwt/v5"
)

type JwtService interface {
	CreateToken(author entity.User) (dto.AuthResponseDto, error)
	ParseToken(tokenHeader string) (jwt.MapClaims, error)
}

type jwtService struct {
	cfg config.TokenConfig
}

func (j *jwtService) CreateToken(user entity.User) (dto.AuthResponseDto, error) {
	claims := model.MyCustomClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    j.cfg.IssuerName,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.cfg.JwtExpiresTime)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		UserId: user.ID,
		Role:   user.Role,
	}

	token := jwt.NewWithClaims(j.cfg.JwtSigningMethod, claims)
	ss, err := token.SignedString(j.cfg.JwtSignatureKy)
	if err != nil {
		return dto.AuthResponseDto{}, fmt.Errorf("oops, failed to create token: %v", err)
	}
	return dto.AuthResponseDto{Token: ss}, nil
}

func (j *jwtService) ParseToken(tokenHeader string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenHeader, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("oops, unexpected signing method: %v", token.Header["alg"])
		}
		return j.cfg.JwtSignatureKy, nil
	})

	if err != nil {
		return nil, fmt.Errorf("oops, failed to verify token: %v", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("oops, failed to parse token claims")
	}
	return claims, nil
}

func NewJwtService(cfg config.TokenConfig) JwtService {
	return &jwtService{cfg: cfg}
}
