package server

import (
	"context"
	"fmt"
	"log"
	"time"

	"main/internal/storage"
	desc "main/pkg/api/monitor/v1"

	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthGRPCServer struct {
	// Обязательно встраиваем, чтобы сервер не падал, если вызовут нереализованный метод
	desc.UnimplementedAuthServiceServer

	// Используем указатель на структуру хранилища (или интерфейс, если он есть)
	storage *storage.InMemoryUser
	// WARNING
	secretKey []byte
}

func NewAuthGRPCServer(store *storage.InMemoryUser) *AuthGRPCServer {
	return &AuthGRPCServer{
		storage:   store,
		secretKey: []byte("sjkehgjikg2378456jksgjkh234yghb4h278"),
	}
}

// ======================================

// user_login -- in context login user req

// ======================================
func (s *AuthGRPCServer) Login(ctx context.Context, req *desc.LoginRequest) (*desc.LoginResponse, error) {
	log.Println("== Вызван Login ==")
	ok, err := s.storage.Login(req.GetLogin(), req.GetPassword())
	if err != nil {
		// Логируем реальную ошибку для дебага на сервере
		// log.Printf("Ошибка БД при логине: %v", err)
		log.Println(err)
		return nil, status.Error(codes.Unauthenticated, "invalid credentials")
	}

	// 2. Если пароль не подошел или юзера нет
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "invalid credentials")
	}

	// 3. Генерируем токен
	myClaim := jwt.MapClaims{}
	myClaim["login"] = req.GetLogin()
	myClaim["exp"] = time.Now().Add(24 * time.Hour).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, myClaim)
	result, err := token.SignedString(s.secretKey)
	if err != nil {
		log.Println(err)
		// Если токен не подписался — это наша вина, отдаем Internal
		// log.Printf("Ошибка подписи JWT: %v", err)
		return nil, status.Error(codes.Internal, "failed to generate token")
	}

	return &desc.LoginResponse{
		Jwt: result,
	}, nil
}

func (s *AuthGRPCServer) CreateUser(ctx context.Context, req *desc.CreateUserRequest) (*desc.CreateUserResponse, error) {
	log.Println("==CreateUSER==")
	err := s.storage.Create(req.GetLogin(), req.GetPassword())
	if err != nil {
		log.Println("Ошибка при создании! ERR = ", err)
		return nil, status.Error(codes.AlreadyExists, fmt.Sprintf("ERR = %v", err))
	}
	log.Println("Пользователь успешно создан!")
	return &desc.CreateUserResponse{}, nil
}

/*
func (s *AuthGRPCServer) GetUser(ctx context.Context, req *desc.GetUserRequest) (*desc.GetUserResponse, error) {
}

func (s *AuthGRPCServer) DeleteUser(ctx context.Context, req *desc.DeleteUserRequest) (*desc.DeleteUserResponse, error) {
}*/
