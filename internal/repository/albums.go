package repository

import (
	"context"

	"github.com/untea/bottom_babruysk/internal/domain"
	"github.com/untea/bottom_babruysk/internal/repository/postgres"
)

type AlbumsRepository struct {
	db *postgres.Client
}

func NewAlbumsRepository(db *postgres.Client) *AlbumsRepository {
	return &AlbumsRepository{db: db}
}

func (r *AlbumsRepository) CreateAlbum(ctx context.Context, request domain.CreateAlbumRequest) (*domain.CreateAlbumResponse, error) {
	const createAlbumSQL = `
		insert into albums (owner_id, title, description, release_date)
		values ($1, $2, $3, $4)
		returning id;
	`

	arguments := []any{
		request.OwnerID,
		request.Title,
		request.Description,
		request.ReleaseDate,
	}

	album, err := postgres.FetchOne[domain.Album](ctx, r.db, createAlbumSQL, arguments...)
	if err != nil {
		return nil, err
	}

	return &domain.CreateAlbumResponse{
		ID: album.ID,
	}, nil
}

func (r *AlbumsRepository) GetAlbum(ctx context.Context, request domain.GetAlbumRequest) (*domain.GetAlbumResponse, error) {
	const getAlbumSQL = `
		select id, owner_id, title, description, release_date, created_at, updated_at
		from albums
		where id = $1;
	`

	arguments := []any{
		request.ID,
	}

	album, err := postgres.FetchOne[domain.Album](ctx, r.db, getAlbumSQL, arguments...)
	if err != nil {
		return nil, err
	}

	return &domain.GetAlbumResponse{
		Album: album,
	}, nil
}

func (r *AlbumsRepository) ListAlbums(ctx context.Context, request domain.ListAlbumsRequest) (*domain.ListAlbumsResponse, error) {
	const getListAlbumsSQL = `
		with params as (
			select
				coalesce(nullif(lower($1), ''), 'created_at') as sort_field,
				coalesce(nullif(lower($2), ''), 'desc')       as sort_order,
				greatest(coalesce($3, 50), 1)                 as limit_val,
				greatest(coalesce($4, 0), 0)                  as offset_val
		)
		select
			a.id, a.owner_id, a.title, a.description, a.release_date, a.created_at, a.updated_at
		from albums a, params p
		order by
			case when p.sort_field = 'title'        and p.sort_order = 'asc'  then a.title        end nulls last,
			case when p.sort_field = 'title'        and p.sort_order = 'desc' then a.title        end desc nulls last,
		
			case when p.sort_field = 'release_date' and p.sort_order = 'asc'  then a.release_date end nulls last,
			case when p.sort_field = 'release_date' and p.sort_order = 'desc' then a.release_date end desc nulls last,
		
			case when p.sort_field = 'created_at'   and p.sort_order = 'asc'  then a.created_at   end nulls last,
			case when p.sort_field = 'created_at'   and p.sort_order = 'desc' then a.created_at   end desc nulls last,
		
			case when p.sort_field = 'updated_at'   and p.sort_order = 'asc'  then a.updated_at   end nulls last,
			case when p.sort_field = 'updated_at'   and p.sort_order = 'desc' then a.updated_at   end desc nulls last,
		
			a.created_at desc
		limit (select limit_val from params)
		offset (select offset_val from params);
	`

	arguments := []any{
		request.SortField,
		request.SortOrder,
		request.Limit,
		request.Offset,
	}

	albums, err := postgres.FetchMany[domain.Album](ctx, r.db, getListAlbumsSQL, arguments...)
	if err != nil {
		return nil, err
	}

	return &domain.ListAlbumsResponse{
		Albums: albums,
	}, nil
}

func (r *AlbumsRepository) UpdateAlbum(ctx context.Context, request domain.UpdateAlbumRequest) error {
	const updateAlbumSQL = `
		update albums
		set
			title        = coalesce($2, title),
			description  = coalesce($3, description),
			release_date = coalesce($4, release_date),
			updated_at   = now()
		where id = $1;
	`

	arguments := []any{
		request.ID,
		request.Title,
		request.Description,
		request.ReleaseDate,
	}

	_, err := postgres.FetchOne[domain.Album](ctx, r.db, updateAlbumSQL, arguments...)
	if err != nil {
		return err
	}

	return nil
}

func (r *AlbumsRepository) DeleteAlbum(ctx context.Context, request domain.DeleteAlbumRequest) error {
	const deleteAlbumSQL = `
		delete from albums where id = $1;
	`

	arguments := []any{
		request.ID,
	}

	affected, err := postgres.ExecAffected(ctx, r.db, deleteAlbumSQL, arguments...)
	if err != nil {
		return err
	}

	if affected == 0 {
		return postgres.ErrNotFound
	}

	return nil
}
