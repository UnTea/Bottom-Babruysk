package repository

import (
	"context"

	"github.com/untea/bottom_babruysk/internal/domain"
	"github.com/untea/bottom_babruysk/internal/repository/postgres"
)

type PlaylistsRepository struct {
	db *postgres.Client
}

func NewPlaylistsRepository(db *postgres.Client) *PlaylistsRepository {
	return &PlaylistsRepository{db: db}
}

func (r *PlaylistsRepository) CreatePlaylist(ctx context.Context, request domain.CreatePlaylistRequest) (*domain.CreatePlaylistResponse, error) {
	const createPlaylistSQL = `
		insert into playlists (owner_id, 
		                       title, 
		                       description, 
		                       visibility)
		values ($1, 
		        $2, 
		        $3, 
		        coalesce($4::visibility, 'private'::visibility))
		returning id;
	`

	arguments := []any{
		request.OwnerID,
		request.Title,
		request.Description,
		request.Visibility,
	}

	playlist, err := postgres.FetchOne[domain.Playlist](ctx, r.db, createPlaylistSQL, arguments...)
	if err != nil {
		return nil, err
	}

	return &domain.CreatePlaylistResponse{
		ID: playlist.ID,
	}, nil
}

func (r *PlaylistsRepository) GetPlaylist(ctx context.Context, request domain.GetPlaylistRequest) (*domain.GetPlaylistResponse, error) {
	const getPlaylistSQL = `
		select 
			id, 
			owner_id, 
			title, 
			description, 
			visibility, 
			created_at, 
			updated_at
		from playlists
		where id = $1;
	`

	arguments := []any{
		request.ID,
	}

	playlist, err := postgres.FetchOne[domain.Playlist](ctx, r.db, getPlaylistSQL, arguments...)
	if err != nil {
		return nil, err
	}

	return &domain.GetPlaylistResponse{
		Playlist: playlist,
	}, nil
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
			p.id, 
			p.owner_id, 
			p.title, 
			p.description, 
			p.visibility, 
			p.created_at, 
			p.updated_at
		from playlists as p, params as par
		where
			(par.owner_filter is null or p.owner_id = par.owner_filter)
			and (par.visibility_filter is null or p.visibility = par.visibility_filter)
		order by
			case when par.sort_field = 'title'      and par.sort_order = 'asc'  then p.title      end nulls last,
			case when par.sort_field = 'title'      and par.sort_order = 'desc' then p.title      end desc nulls last,

			case when par.sort_field = 'created_at' and par.sort_order = 'asc'  then p.created_at end nulls last,
			case when par.sort_field = 'created_at' and par.sort_order = 'desc' then p.created_at end desc nulls last,

			case when par.sort_field = 'updated_at' and par.sort_order = 'asc'  then p.updated_at end nulls last,
			case when par.sort_field = 'updated_at' and par.sort_order = 'desc' then p.updated_at end desc nulls last,

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

	playlists, err := postgres.FetchMany[domain.Playlist](ctx, r.db, listPlaylistSQL, arguments...)
	if err != nil {
		return nil, err
	}

	return &domain.ListPlaylistsResponse{
		Playlists: playlists,
	}, nil
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

	_, err := postgres.FetchOne[domain.Playlist](ctx, r.db, updatePlaylistSQL, arguments...)
	if err != nil {
		return err
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

	affected, err := postgres.ExecAffected(ctx, r.db, deletePlaylistSQL, arguments...)
	if err != nil {
		return err
	}

	if affected == 0 {
		return postgres.ErrNotFound
	}

	return nil
}
