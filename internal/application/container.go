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
	AlbumServices *service.AlbumsService
	TacksServices *service.TracksService
}

type Repositories struct {
	UsersRepository  repository.Users
	AlbumsRepository repository.Albums
	TracksRepository repository.Tracks
}

type Container struct {
	Logger        *zap.Logger
	Configuration *configuration.Configuration
	Repositories  Repositories
	Services      Services
}

func BuildContainer(configuration *configuration.Configuration, logger *zap.Logger, dbClient *repository.Client) (*Container, error) {
	usersRepository := postgres.NewUsersRepository(dbClient)
	albumsRepository := postgres.NewAlbumsRepository(dbClient)
	tracksRepository := postgres.NewTracksRepository(dbClient)

	repositories := Repositories{
		UsersRepository:  usersRepository,
		AlbumsRepository: albumsRepository,
		TracksRepository: tracksRepository,
	}

	usersServices := service.NewUsersService(usersRepository)
	albumsServices := service.NewAlbumsService(albumsRepository)
	tracksServices := service.NewTracksService(tracksRepository)

	services := Services{
		UsersServices: usersServices,
		AlbumServices: albumsServices,
		TacksServices: tracksServices,
	}

	container := &Container{
		Logger:        logger,
		Configuration: configuration,
		Repositories:  repositories,
		Services:      services,
	}

	return container, nil
}
