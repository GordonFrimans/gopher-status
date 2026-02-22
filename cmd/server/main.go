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

	MonitorServer := server.NewMonitorGRPCServer(store, workers)
	AuthServer := server.NewAuthGRPCServer(storage.NewInMemoryUser())

	// TEST WARNING

	// 2. Запуск gRPC сервера (в отдельной горутине)
	go func() {
		lis, err := net.Listen("tcp", ":50051")
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}
		secretKey := []byte("sjkehgjikg2378456jksgjkh234yghb4h278")
		s := grpc.NewServer(
			// Добавляем наш интерцептор ко всем входящим запросам
			grpc.UnaryInterceptor(server.AuthInterceptor(secretKey)),
		)
		//====  Servers ====

		// --Montior server
		desc.RegisterMonitorServiceServer(s, MonitorServer)
		// --Auth sever
		desc.RegisterAuthServiceServer(s, AuthServer)

		//====  Servers ====
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

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	// 1. Регистрируем MonitorService
	err := desc.RegisterMonitorServiceHandlerFromEndpoint(ctx, mux, "localhost:50051", opts)
	if err != nil {
		log.Fatalf("failed to register monitor gateway: %v", err)
	}

	// 2. ДОБАВЛЯЕМ ЭТО: Регистрируем AuthService!
	err = desc.RegisterAuthServiceHandlerFromEndpoint(ctx, mux, "localhost:50051", opts)
	if err != nil {
		log.Fatalf("failed to register auth gateway: %v", err)
	}

	httpMux := http.NewServeMux()
	httpMux.Handle("/v1/", mux)

	fileServer := http.FileServer(http.Dir("/app/web/dist"))
	httpMux.Handle("/", fileServer)

	log.Println("Serving API & UI on 0.0.0.0:8080")
	if err := http.ListenAndServe(":8080", httpMux); err != nil {
		log.Fatalf("failed to serve HTTP: %v", err)
	}
}
