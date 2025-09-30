package application

import (
	"go.uber.org/zap"

	"github.com/untea/bottom_babruysk/internal/configuration"
	"github.com/untea/bottom_babruysk/internal/repository"
	"github.com/untea/bottom_babruysk/internal/repository/postgres"
	"github.com/untea/bottom_babruysk/internal/service"
)

type Services struct {
	UsersServices    *service.UsersService
	AlbumServices    *service.AlbumsService
	TacksServices    *service.TracksService
	PlaylistsService *service.PlaylistsService
}

type Repositories struct {
	UsersRepository     repository.Users
	AlbumsRepository    repository.Albums
	TracksRepository    repository.Tracks
	PlaylistsRepository repository.Playlists
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
	playlistsRepository := postgres.NewPlaylistsRepository(dbClient)

	repositories := Repositories{
		UsersRepository:     usersRepository,
		AlbumsRepository:    albumsRepository,
		TracksRepository:    tracksRepository,
		PlaylistsRepository: playlistsRepository,
	}

	usersServices := service.NewUsersService(usersRepository)
	albumsServices := service.NewAlbumsService(albumsRepository)
	tracksServices := service.NewTracksService(tracksRepository)
	playlistsServices := service.NewPlaylistsService(playlistsRepository)

	services := Services{
		UsersServices:    usersServices,
		AlbumServices:    albumsServices,
		TacksServices:    tracksServices,
		PlaylistsService: playlistsServices,
	}

	container := &Container{
		Logger:        logger,
		Configuration: configuration,
		Repositories:  repositories,
		Services:      services,
	}

	return container, nil
}
