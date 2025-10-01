package application

import (
	"go.uber.org/zap"

	"github.com/untea/bottom_babruysk/internal/configuration"
	"github.com/untea/bottom_babruysk/internal/repository"
	"github.com/untea/bottom_babruysk/internal/repository/postgres"
	"github.com/untea/bottom_babruysk/internal/service"
)

type Services struct {
	UsersServices     *service.UsersService
	AlbumServices     *service.AlbumsService
	TacksServices     *service.TracksService
	PlaylistsService  *service.PlaylistsService
	ArtistsService    *service.ArtistsService
	TrackFilesService *service.TrackFilesService
}

type Repositories struct {
	UsersRepository      service.Users
	AlbumsRepository     service.Albums
	TracksRepository     service.Tracks
	PlaylistsRepository  service.Playlists
	ArtistsRepository    service.Artists
	TrackFilesRepository service.TrackFiles
}

type Container struct {
	Logger        *zap.Logger
	Configuration *configuration.Configuration
	Repositories  Repositories
	Services      Services
}

func BuildContainer(configuration *configuration.Configuration, logger *zap.Logger, dbClient *postgres.Client) (*Container, error) {
	usersRepository := repository.NewUsersRepository(dbClient)
	albumsRepository := repository.NewAlbumsRepository(dbClient)
	tracksRepository := repository.NewTracksRepository(dbClient)
	playlistsRepository := repository.NewPlaylistsRepository(dbClient)
	artistsRepository := repository.NewArtistsRepository(dbClient)
	trackFilesRepository := repository.NewTrackFilesRepository(dbClient)

	repositories := Repositories{
		UsersRepository:      usersRepository,
		AlbumsRepository:     albumsRepository,
		TracksRepository:     tracksRepository,
		PlaylistsRepository:  playlistsRepository,
		ArtistsRepository:    artistsRepository,
		TrackFilesRepository: trackFilesRepository,
	}

	usersServices := service.NewUsersService(usersRepository)
	albumsServices := service.NewAlbumsService(albumsRepository)
	tracksServices := service.NewTracksService(tracksRepository)
	playlistsServices := service.NewPlaylistsService(playlistsRepository)
	artistsServices := service.NewArtistsService(artistsRepository)
	trackFilesService := service.NewTrackFilesService(trackFilesRepository)

	services := Services{
		UsersServices:     usersServices,
		AlbumServices:     albumsServices,
		TacksServices:     tracksServices,
		PlaylistsService:  playlistsServices,
		ArtistsService:    artistsServices,
		TrackFilesService: trackFilesService,
	}

	container := &Container{
		Logger:        logger,
		Configuration: configuration,
		Repositories:  repositories,
		Services:      services,
	}

	return container, nil
}
