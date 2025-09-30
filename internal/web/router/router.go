package router

import (
	"net/http"
	"time"

	"connectrpc.com/grpcreflect"
	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"go.uber.org/zap"

	"github.com/untea/bottom_babruysk/internal/application"
	"github.com/untea/bottom_babruysk/internal/transport/connect"
	"github.com/untea/bottom_babruysk/internal/transport/gen/proto/v1/protov1connect"
	"github.com/untea/bottom_babruysk/internal/web"
	webMiddleware "github.com/untea/bottom_babruysk/internal/web/middleware"
)

type Dependencies struct {
	Logger   *zap.Logger
	Handlers web.HandlerHTTP
	Services application.Services

	EnableCORS       bool
	EnableReflection bool
}

func New(dependencies Dependencies) *chi.Mux {
	r := chi.NewRouter()

	if dependencies.EnableCORS {
		r.Use(cors.Handler(cors.Options{
			AllowedOrigins:   []string{"https://*", "http://*"},
			AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
			AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
			ExposedHeaders:   []string{"Link"},
			AllowCredentials: true,
			MaxAge:           300,
		}))
	}

	r.Use(
		chiMiddleware.RequestID,
		chiMiddleware.RealIP,
		chiMiddleware.Recoverer,
		chiMiddleware.Timeout(60*time.Second),
		webMiddleware.RequestLogger(dependencies.Logger),
	)

	r.Get("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// REST
	r.Route("/api/v1", func(api chi.Router) {
		dependencies.Handlers.MountUsers(api)
		dependencies.Handlers.MountAlbums(api)
		dependencies.Handlers.MountTracks(api)
		dependencies.Handlers.MountPlaylists(api)
		dependencies.Handlers.MountArtists(api)
	})

	// CONNECT RPC
	usersServer := connect.NewUsersServer(dependencies.Services.UsersServices)
	usersPath, usersHandler := protov1connect.NewUsersServiceHandler(usersServer)
	HandleStaticReflect(r, usersPath, usersHandler)

	albumsServer := connect.NewAlbumsServer(dependencies.Services.AlbumServices)
	albumsPath, albumsHandler := protov1connect.NewAlbumsServiceHandler(albumsServer)
	HandleStaticReflect(r, albumsPath, albumsHandler)

	tracksServer := connect.NewTracksServer(dependencies.Services.TacksServices)
	tracksPath, tracksHandler := protov1connect.NewTracksServiceHandler(tracksServer)
	HandleStaticReflect(r, tracksPath, tracksHandler)

	playlistsServer := connect.NewPlaylistsServer(dependencies.Services.PlaylistsService)
	playlistsPath, playlistsHandler := protov1connect.NewPlaylistsServiceHandler(playlistsServer)
	HandleStaticReflect(r, playlistsPath, playlistsHandler)

	artistsServer := connect.NewArtistsServer(dependencies.Services.ArtistsService)
	artistsPath, artistsHandler := protov1connect.NewArtistsServiceHandler(artistsServer)
	HandleStaticReflect(r, artistsPath, artistsHandler)

	if dependencies.EnableReflection {
		reflector := grpcreflect.NewStaticReflector(
			protov1connect.UsersServiceName,
			protov1connect.AlbumsServiceName,
			protov1connect.TracksServiceName,
			protov1connect.PlaylistsServiceName,
			protov1connect.ArtistsServiceName,
		)

		v1Path, v1Handler := grpcreflect.NewHandlerV1(reflector)
		HandleStaticReflect(r, v1Path, v1Handler)

		v1aPath, v1aHandler := grpcreflect.NewHandlerV1Alpha(reflector)
		HandleStaticReflect(r, v1aPath, v1aHandler)
	}

	return r
}
