package application

import (
	"go.uber.org/zap"

	"github.com/untea/bottom_babruysk/internal/configuration"
	"github.com/untea/bottom_babruysk/internal/repository"
	"github.com/untea/bottom_babruysk/internal/repository/postgres"
	"github.com/untea/bottom_babruysk/internal/service"
)

type Services struct {
	UsersServices *service.UsersService
}

type Repositories struct {
	UsersRepository repository.Users
}

type Container struct {
	Logger        *zap.Logger
	Configuration *configuration.Configuration
	Repositories  Repositories
	Services      Services
}

func BuildContainer(configuration *configuration.Configuration, logger *zap.Logger, dbClient *repository.Client) (*Container, error) {
	usersRepository := postgres.NewUsersRepository(dbClient)

	repositories := Repositories{
		UsersRepository: usersRepository,
	}

	usersServices := service.NewUsersService(usersRepository)

	services := Services{
		UsersServices: usersServices,
	}

	container := &Container{
		Logger:        logger,
		Configuration: configuration,
		Repositories:  repositories,
		Services:      services,
	}

	return container, nil
}
