package grpc

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/mesameen/micro-app/src/api/gen"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// SecretProvider defines a provider of secrets for our handler
type SecretProvider func() []byte

// Handler defines an auth gRPC handler
type Handler struct {
	gen.UnimplementedAuthServiceServer
	secretProvider SecretProvider
}

// New creates a new auth gRPC service
func New(secretProvider SecretProvider) *Handler {
	return &Handler{
		secretProvider: secretProvider,
	}
}

// GetToken performs verfication of user credentails and returns a JWT token in case of success
func (h *Handler) GetToken(ctx context.Context, req *gen.GetTokenRequest) (*gen.GetTokenResponse, error) {
	if !validateCredentials(req.Username, req.Password) {
		return nil, status.Error(codes.Unauthenticated, "invalid credentials")
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": req.Username,
		"iat":      time.Now().Unix(),
	})
	tokenString, err := token.SignedString(h.secretProvider())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &gen.GetTokenResponse{
		Token: tokenString,
	}, nil
}

func validateCredentials(username, password string) bool {
	if username == "" || password == "" {
		return false
	}
	return true
}

// ValidateToken performs JWT token validation
func (h *Handler) ValidateToken(ctx context.Context, req *gen.ValidateTokenRequest) (*gen.ValidateTokenResponse, error) {
	token, err := jwt.Parse(req.Token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signning method: %v", t.Header["alg"])
		}
		return h.secretProvider(), nil
	})
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid token")
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, status.Error(codes.Unauthenticated, "invalid token")
	}
	var username string
	if v, ok := claims["username"]; ok {
		if u, ok := v.(string); ok {
			username = u
		}
	}
	return &gen.ValidateTokenResponse{Username: username}, nil
}
