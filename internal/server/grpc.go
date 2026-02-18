package server

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"main/internal/storage"
	"main/internal/worker"
	desc "main/pkg/api/monitor/v1"
)

type GRPCServer struct {
	// Обязательно встраиваем, чтобы сервер не падал, если вызовут нереализованный метод
	desc.UnimplementedMonitorServiceServer

	// Используем указатель на структуру хранилища (или интерфейс, если он есть)
	storage *storage.InMemoryStorageMonitors
	workers *worker.WorkerPool
}

func NewGRPCServer(store *storage.InMemoryStorageMonitors, workers *worker.WorkerPool) *GRPCServer {
	return &GRPCServer{
		storage: store,
		workers: workers,
	}
}

func (s *GRPCServer) CreateMonitor(ctx context.Context, req *desc.CreateMonitorRequest) (*desc.CreateMonitorResponse, error) {
	log.Println("CreateMonitor trigg")
	newMonitor := storage.Monitor{
		ID:        s.storage.GetLastID(),
		URL:       req.GetUrl(),
		Name:      req.GetName(),
		Interval:  req.GetInterval(),
		Status:    "PENDING",
		LastCheck: "", // Дефолт значения времени при инициализации!
	}
	newID, err := s.storage.Create(newMonitor)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("ERR=%v", err))
	}
	return &desc.CreateMonitorResponse{
		Id: newID,
	}, nil
}

func (s *GRPCServer) ListMonitors(ctx context.Context, req *desc.ListMonitorsRequest) (*desc.ListMonitorsResponse, error) {
	log.Println("ListMonitors trigg")
	monitors, err := s.storage.List()
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("ERR=%v", err))
	}

	// Конвертируем []Monitor в []*desc.Monitor
	protoMonitors := make([]*desc.Monitor, 0, len(monitors))
	for _, m := range monitors {
		protoMonitors = append(protoMonitors, &desc.Monitor{
			Id:        m.ID,
			Url:       m.URL,
			Name:      m.Name,
			Interval:  m.Interval,
			Status:    m.Status,
			LastCheck: m.LastCheck,
		})
	}

	return &desc.ListMonitorsResponse{
		Monitors: protoMonitors,
	}, nil
}

func (s *GRPCServer) DeleteMonitor(ctx context.Context, req *desc.DeleteMonitorRequest) (*desc.DeleteMonitorResponse, error) {
	log.Println("DeleteMonitor trigg")
	err := s.storage.Delete(req.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("ERR=%v", err))
	}
	return &desc.DeleteMonitorResponse{}, nil
}
