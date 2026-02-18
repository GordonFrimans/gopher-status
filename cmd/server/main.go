package main

import (
	"context"
	"log"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"

	"main/internal/server"
	"main/internal/storage"
	"main/internal/worker"
	desc "main/pkg/api/monitor/v1"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

func main() {
	store := storage.NewInMemoryStorageMonitors()
	// Воркеры для пинговки сайтов (управление в хандлерах )
	workers := worker.NewWorkerPool(5, worker.SimpleCheck)
	err := workers.Start()
	if err != nil {
		log.Println("[ERR]Ошибка запуска ворверов = ", err)
	}

	shelduler := worker.NewScheduler(workers, store)
	shelduler.Start()

	myServer := server.NewGRPCServer(store, workers)
	// TEST WARNING

	// 2. Запуск gRPC сервера (в отдельной горутине)
	go func() {
		lis, err := net.Listen("tcp", ":50051")
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}
		s := grpc.NewServer()
		desc.RegisterMonitorServiceServer(s, myServer)
		reflection.Register(s)

		log.Println("Serving gRPC on 0.0.0.0:50051")
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()
	runGateway()
}

func runGateway() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Создаем мультиплексор (роутер) для Gateway
	mux := runtime.NewServeMux()

	// Опции подключения к gRPC серверу (без TLS пока что)
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	// Регистрируем хендлер, который будет транслировать HTTP в gRPC
	// Обрати внимание: мы подключаемся к НАШЕМУ ЖЕ gRPC серверу на 50051
	err := desc.RegisterMonitorServiceHandlerFromEndpoint(ctx, mux, "localhost:50051", opts)
	if err != nil {
		log.Fatalf("failed to register gateway: %v", err)
	}

	httpMux := http.NewServeMux()

	// 1. Все запросы на /v1/ отправляем в gRPC Gateway
	httpMux.Handle("/v1/", mux)

	// 2. Все остальные запросы отправляем в папку dist (React)
	// Важно: http.FileServer умеет отдавать файлы
	fileServer := http.FileServer(http.Dir("/app/web/dist"))
	httpMux.Handle("/", fileServer)

	log.Println("Serving API & UI on 0.0.0.0:8080")
	// Запускаем сервер на одном порту!
	if err := http.ListenAndServe(":8080", httpMux); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
