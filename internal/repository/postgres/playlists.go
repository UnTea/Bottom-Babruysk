package postgres

import (
	"context"

	"github.com/untea/bottom_babruysk/internal/domain"
	"github.com/untea/bottom_babruysk/internal/repository"
)

type PlaylistsRepository struct {
	db *repository.Client // TODO реализовать интерфейс для fetch и прокидывать просто r.db
}

func NewPlaylistsRepository(db *repository.Client) *PlaylistsRepository {
	return &PlaylistsRepository{db: db}
}

func (r *PlaylistsRepository) CreatePlaylist(ctx context.Context, request domain.CreatePlaylistRequest) (*domain.CreatePlaylistResponse, error) {
	const createPlaylistSQL = `
		insert into playlists (owner_id, title, description, visibility)
		values ($1, $2, $3, coalesce($4::visibility, 'private'::visibility))
		returning id;
	`

	arguments := []any{
		request.OwnerID,
		request.Title,
		request.Description,
		request.Visibility,
	}

	playlist, err := repository.FetchOne[domain.Playlist](ctx, r.db.Driver(), createPlaylistSQL, arguments...) // TODO реализовать интерфейс для fetch и прокидывать просто r.db
	if err != nil {
		return nil, err
	}

	return &domain.CreatePlaylistResponse{ID: playlist.ID}, nil
}

func (r *PlaylistsRepository) GetPlaylist(ctx context.Context, request domain.GetPlaylistRequest) (*domain.GetPlaylistResponse, error) {
	const getPlaylistSQL = `
		select id, owner_id, title, description, visibility::text as visibility, created_at, updated_at
		from playlists
		where id = $1;
	`

	arguments := []any{
		request.ID,
	}

	playlist, err := repository.FetchOne[domain.Playlist](ctx, r.db.Driver(), getPlaylistSQL, arguments...) // TODO реализовать интерфейс для fetch и прокидывать просто r.db
	if err != nil {
		return nil, err
	}

	return &domain.GetPlaylistResponse{Playlist: playlist}, nil
}

func (r *PlaylistsRepository) ListPlaylists(ctx context.Context, request domain.ListPlaylistsRequest) (*domain.ListPlaylistsResponse, error) {
	const listPlaylistSQL = `
		with params as (
			select
				$1::uuid                                      as owner_filter,
				$2::visibility                                as visibility_filter,
				coalesce(nullif(lower($3), ''), 'created_at') as sort_field,
				coalesce(nullif(lower($4), ''), 'desc')       as sort_order,
				greatest(coalesce($5, 50), 1)                 as limit_val,
				greatest(coalesce($6, 0), 0)                  as offset_val
		)
		select
			p.id, p.owner_id, p.title, p.description, p.visibility::text as visibility, p.created_at, p.updated_at
		from playlists p, params p2
		where
			(p2.owner_filter is null or p.owner_id = p2.owner_filter)
			and (p2.visibility_filter is null or p.visibility = p2.visibility_filter)
		order by
			case when p2.sort_field = 'title'      and p2.sort_order = 'asc'  then p.title      end  nulls last,
			case when p2.sort_field = 'title'      and p2.sort_order = 'desc' then p.title      end desc nulls last,

			case when p2.sort_field = 'created_at' and p2.sort_order = 'asc'  then p.created_at end  nulls last,
			case when p2.sort_field = 'created_at' and p2.sort_order = 'desc' then p.created_at end desc nulls last,

			case when p2.sort_field = 'updated_at' and p2.sort_order = 'asc'  then p.updated_at end  nulls last,
			case when p2.sort_field = 'updated_at' and p2.sort_order = 'desc' then p.updated_at end desc nulls last,

			p.created_at desc
		limit (select limit_val from params)
		offset (select offset_val from params);
	`

	arguments := []any{
		request.OwnerID,
		request.Visibility,
		request.SortField,
		request.SortOrder,
		request.Limit,
		request.Offset,
	}

	playlists, err := repository.FetchMany[domain.Playlist](ctx, r.db.Driver(), listPlaylistSQL, arguments...) // TODO реализовать интерфейс для fetch и прокидывать просто r.db
	if err != nil {
		return nil, err
	}

	return &domain.ListPlaylistsResponse{Playlists: playlists}, nil
}

func (r *PlaylistsRepository) UpdatePlaylist(ctx context.Context, request domain.UpdatePlaylistRequest) error {
	const updatePlaylistSQL = `
		update playlists
		set
			title       = coalesce($2, title),
			description = coalesce($3, description),
			visibility  = coalesce($4::visibility, visibility),
			updated_at  = now()
		where id = $1;
	`

	arguments := []any{
		request.ID,
		request.Title,
		request.Description,
		request.Visibility,
	}

	affected, err := repository.ExecAffected(ctx, r.db.Driver(), updatePlaylistSQL, arguments...)
	if err != nil {
		return err
	}

	if affected == 0 {
		return repository.ErrNotFound
	}

	return nil
}

func (r *PlaylistsRepository) DeletePlaylist(ctx context.Context, request domain.DeletePlaylistRequest) error {
	const deletePlaylistSQL = `
		delete from playlists where id = $1;
	`

	arguments := []any{
		request.ID,
	}

	affected, err := repository.ExecAffected(ctx, r.db.Driver(), deletePlaylistSQL, arguments...) // TODO реализовать интерфейс для fetch и прокидывать просто r.db
	if err != nil {
		return err
	}

	if affected == 0 {
		return repository.ErrNotFound
	}

	return nil
}
