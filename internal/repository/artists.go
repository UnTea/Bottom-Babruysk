package repository

import (
	"context"

	"github.com/untea/bottom_babruysk/internal/domain"
	"github.com/untea/bottom_babruysk/internal/repository/postgres"
)

type ArtistsRepository struct {
	db *postgres.Client
}

func NewArtistsRepository(db *postgres.Client) *ArtistsRepository {
	return &ArtistsRepository{db: db}
}

func (r *ArtistsRepository) CreateArtist(ctx context.Context, request domain.CreateArtistRequest) (*domain.CreateArtistResponse, error) {
	const createArtistSQL = `
		insert into artists (name, 
		                     bio)
		values ($1, 
		        $2)
		returning id;
	`

	arguments := []any{
		request.Name,
		request.Bio,
	}

	artist, err := postgres.FetchOne[domain.Artist](ctx, r.db, createArtistSQL, arguments...)
	if err != nil {
		return nil, err
	}

	return &domain.CreateArtistResponse{
		ID: artist.ID,
	}, nil
}

func (r *ArtistsRepository) GetArtist(ctx context.Context, request domain.GetArtistRequest) (*domain.GetArtistResponse, error) {
	const getArtistSQL = `
		select 
			id, 
			name, 
			bio, 
			created_at, 
			updated_at
		from artists
		where id = $1;
	`

	arguments := []any{
		request.ID,
	}

	artist, err := postgres.FetchOne[domain.Artist](ctx, r.db, getArtistSQL, arguments...)
	if err != nil {
		return nil, err
	}

	return &domain.GetArtistResponse{
		Artist: artist,
	}, nil
}

func (r *ArtistsRepository) ListArtists(ctx context.Context, request domain.ListArtistsRequest) (*domain.ListArtistsResponse, error) {
	const listArtistsSQL = `
		with params as (
			select
				nullif($1, '')                                as name_q,
				nullif($2, '')                                as bio_q,
				coalesce(nullif(lower($3), ''), 'created_at') as sort_field,
				coalesce(nullif(lower($4), ''), 'desc')       as sort_order,
				greatest(coalesce($5, 50), 1)                 as limit_val,
				greatest(coalesce($6, 0), 0)                  as offset_val
		)
		select 
			a.id,
			a.name,
			a.bio,
			a.created_at,
			a.updated_at
		from artists as a, params as p
		where
			(p.name_q is null or a.name ilike ('%' || p.name_q || '%'))
			and (p.bio_q  is null or a.bio  ilike ('%' || p.bio_q  || '%'))
		order by
			case when p.sort_field = 'name'       and p.sort_order = 'asc'  then a.name       end nulls last,
			case when p.sort_field = 'name'       and p.sort_order = 'desc' then a.name       end desc nulls last,

			case when p.sort_field = 'created_at' and p.sort_order = 'asc'  then a.created_at end nulls last,
			case when p.sort_field = 'created_at' and p.sort_order = 'desc' then a.created_at end desc nulls last,

			case when p.sort_field = 'updated_at' and p.sort_order = 'asc'  then a.updated_at end nulls last,
			case when p.sort_field = 'updated_at' and p.sort_order = 'desc' then a.updated_at end desc nulls last,

			a.created_at desc
		limit (select limit_val  from params)
		offset (select offset_val from params);
	`

	arguments := []any{
		request.Name,
		request.Bio,
		request.SortField,
		request.SortOrder,
		request.Limit,
		request.Offset,
	}

	artists, err := postgres.FetchMany[domain.Artist](ctx, r.db, listArtistsSQL, arguments...)
	if err != nil {
		return nil, err
	}

	return &domain.ListArtistsResponse{
		Artists: artists,
	}, nil
}

func (r *ArtistsRepository) UpdateArtist(ctx context.Context, request domain.UpdateArtistRequest) error {
	const updateArtistSQL = `
		update artists
		set
			name       = coalesce($2, name),
			bio        = coalesce($3, bio),
			updated_at = now()
		where id = $1;
	`

	arguments := []any{
		request.ID,
		request.Name,
		request.Bio,
	}

	_, err := postgres.FetchOne[domain.Artist](ctx, r.db, updateArtistSQL, arguments...)
	if err != nil {
		return err
	}

	return nil
}

func (r *ArtistsRepository) DeleteArtist(ctx context.Context, request domain.DeleteArtistRequest) error {
	const deleteArtistSQL = `
		delete from artists where id = $1;
	`

	affected, err := postgres.ExecAffected(ctx, r.db, deleteArtistSQL, *request.ID)
	if err != nil {
		return err
	}

	if affected == 0 {
		return postgres.ErrNotFound
	}

	return nil
}
