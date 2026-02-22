package server

import (
	"context"
	"strings"

	"github.com/golang-jwt/jwt/v5" // Убедись, что версия совпадает с твоей
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// Ключ для контекста (лучше использовать кастомный тип, чтобы не было коллизий)
type contextKey string

const UserLoginKey contextKey = "user_login"

func AuthInterceptor(secretKey []byte) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		// 1. Пропускаем методы, для которых не нужна авторизация (логин и рега)
		// ВАЖНО: Замени пути на свои реальные названия из .proto файла!
		if info.FullMethod == "/api.monitor.v1.AuthService/Login" ||
			info.FullMethod == "/api.monitor.v1.AuthService/CreateUser" {
			return handler(ctx, req)
		}

		// 2. Достаем метаданные (заголовки) из запроса
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "metadata is not provided")
		}

		// 3. Ищем заголовок Authorization
		values := md["authorization"]
		if len(values) == 0 {
			return nil, status.Error(codes.Unauthenticated, "authorization token is not provided")
		}

		// 4. Проверяем формат "Bearer <token>"
		accessToken := values[0]
		if !strings.HasPrefix(accessToken, "Bearer ") {
			return nil, status.Error(codes.Unauthenticated, "invalid token format")
		}
		tokenString := strings.TrimPrefix(accessToken, "Bearer ")

		// 5. Парсим и валидируем токен
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, status.Errorf(codes.Unauthenticated, "unexpected signing method")
			}
			return secretKey, nil
		})

		if err != nil || !token.Valid {
			return nil, status.Error(codes.Unauthenticated, "invalid or expired token")
		}

		// 6. Достаем логин из токена
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "invalid token claims")
		}

		login, ok := claims["login"].(string)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "login not found in token")
		}

		// 7. Пробрасываем логин в контекст, чтобы методы могли понять, кто к ним пришел
		newCtx := context.WithValue(ctx, UserLoginKey, login)

		// Передаем управление дальше, в сам метод
		return handler(newCtx, req)
	}
}
