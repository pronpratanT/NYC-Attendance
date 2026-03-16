package service

import (
	"context"
	"encoding/json"
	"errors"
	"hr-program/internal/user-service/dto"
	"hr-program/shared/auth"
	"hr-program/shared/config"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

func (s *UserService) Login(employeeID, password, ip, userAgent string, redisClient *redis.Client) (*dto.LoginResponse, error) {
	user, err := s.AppRepo.GetUserByEmployeeID(employeeID)
	if err != nil {
		return nil, err
	}

	if user.Password != password {
		return nil, errors.New("invalid credentials")
	}

	// production ควรใช้ bcrypt แทน
	// if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
	// 	return nil, errors.New("invalid credentials")
	// }

	jti := uuid.NewString()
	ttlMinutes, _ := strconv.Atoi(config.AppConfig.JWTAccessTTLMinutes)
	ttl := time.Duration(ttlMinutes) * time.Minute

	token, expiresAt, err := auth.GenerateAccessToken(
		config.AppConfig.JWTSecret,
		user.ID,
		user.EmployeeID,
		jti,
		ttl,
	)
	if err != nil {
		return nil, err
	}

	session := auth.SessionData{
		UserID:     user.ID,
		EmployeeID: user.EmployeeID,
		JTI:        jti,
		IssuedAt:   time.Now().Unix(),
		ExpiresAt:  expiresAt,
		IP:         ip,
		UserAgent:  userAgent,
		Revoked:    false,
	}

	raw, err := json.Marshal(session)
	if err != nil {
		return nil, err
	}

	key := "auth:session:" + jti
	if err := redisClient.Set(context.Background(), key, raw, ttl).Err(); err != nil {
		return nil, err
	}

	return &dto.LoginResponse{
		UserID:       user.ID,
		EmployeeID:   user.EmployeeID,
		DepartmentID: user.DepartmentID,
		FName:        user.FName,
		LName:        user.LName,
		IsActive:     user.IsActive,
		AccessToken:  token,
		TokenType:    "Bearer",
		ExpiresAt:    expiresAt,
	}, nil
}

func (s *UserService) LogOut(jti string, redisClient *redis.Client) error {
	key := "auth:session:" + jti
	return redisClient.Del(context.Background(), key).Err()
}
