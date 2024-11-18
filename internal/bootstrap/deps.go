package bootstrap

import (
	"github.com/bjorndonald/lasgcce/internal/repository"
	"github.com/bjorndonald/lasgcce/internal/service"
	"github.com/bjorndonald/lasgcce/internal/service/streaming"
	"gorm.io/gorm"
)

type AppDependencies struct {
	EmailService  service.EmailServicer
	UserRepo      repository.UserRepositoryInterface
	LocationRepo  repository.LocationRepositoryInterface
	AgentRepo     repository.AgentRepositoryInterface
	EventProducer streaming.EventProducer
	// EventProducer   sarama.SyncProducer
	DatabaseService *gorm.DB
}

func InitializeDependencies(db *gorm.DB) *AppDependencies {
	return &AppDependencies{
		UserRepo:        repository.NewUserRepository(db),
		LocationRepo:    repository.NewLocationRepository(db),
		AgentRepo:       repository.NewAgentRepository(db),
		EmailService:    service.NewEmailService(),
		DatabaseService: db,
	}
}
