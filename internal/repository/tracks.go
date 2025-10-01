package repository

import (
	"context"

	"github.com/untea/bottom_babruysk/internal/domain"
	"github.com/untea/bottom_babruysk/internal/repository/postgres"
)

type TracksRepository struct {
	db *postgres.Client
}

func NewTracksRepository(db *postgres.Client) *TracksRepository {
	return &TracksRepository{db: db}
}

func (r *TracksRepository) CreateTrack(ctx context.Context, request domain.CreateTrackRequest) (*domain.CreateTrackResponse, error) {
	const createTracksQL = `
		insert into tracks (uploader_id, 
		                    title, 
		                    subtitle, 
		                    description, 
		                    duration, 
							visibility, 
		                    uploaded_at)
		values ($1, 
		        $2, 
		        $3, 
		        $4, 
		        $5, 
		        coalesce($6::visibility, 'private'::visibility), 
		        $7)
		returning id;
	`

	arguments := []any{
		request.UploaderID,
		request.Title,
		request.Subtitle,
		request.Description,
		request.Duration,
		request.Visibility,
		request.UploadedAt,
	}

	track, err := postgres.FetchOne[domain.Track](ctx, r.db, createTracksQL, arguments...)
	if err != nil {
		return nil, err
	}

	return &domain.CreateTrackResponse{
		ID: track.ID,
	}, nil
}

func (r *TracksRepository) GetTrack(ctx context.Context, request domain.GetTrackRequest) (*domain.GetTrackResponse, error) {
	const getTracksQL = `
		select 
			id, 
			uploader_id, 
			title, 
			subtitle, 
			description, 
			duration, 
			visibility, 
			created_at, 
			updated_at, 
			uploaded_at
		from tracks
		where id = $1;
	`

	arguments := []any{
		request.ID,
	}

	track, err := postgres.FetchOne[domain.Track](ctx, r.db, getTracksQL, arguments...)
	if err != nil {
		return nil, err
	}

	return &domain.GetTrackResponse{
		Track: track,
	}, nil
}

func (r *TracksRepository) ListTracks(ctx context.Context, request domain.ListTracksRequest) (*domain.ListTracksResponse, error) {
	const getListTracksSQL = `
		with params as (
			select
				$1::uuid                                      as uploader_filter,
				$2::visibility                                as visibility_filter,
				coalesce(nullif(lower($3), ''), 'created_at') as sort_field,
				coalesce(nullif(lower($4), ''), 'desc')       as sort_order,
				greatest(coalesce($5, 50), 1)                 as limit_val,
				greatest(coalesce($6, 0), 0)                  as offset_val
		)
		select
			t.id, 
			t.uploader_id, 
			t.title, 
			t.subtitle, 
			t.description, 
			t.duration,
			t.visibility,
			t.created_at, 
			t.updated_at, 
			t.uploaded_at
		from tracks as t, params as p
		where
			(p.uploader_filter is null or t.uploader_id = p.uploader_filter)
			and (p.visibility_filter is null or t.visibility = p.visibility_filter)
		order by
			case when p.sort_field = 'title'       and p.sort_order = 'asc'  then t.title       end nulls last,
			case when p.sort_field = 'title'       and p.sort_order = 'desc' then t.title       end desc nulls last,

			case when p.sort_field = 'uploaded_at' and p.sort_order = 'asc'  then t.uploaded_at end nulls last,
			case when p.sort_field = 'uploaded_at' and p.sort_order = 'desc' then t.uploaded_at end desc nulls last,

			case when p.sort_field = 'created_at'  and p.sort_order = 'asc'  then t.created_at  end nulls last,
			case when p.sort_field = 'created_at'  and p.sort_order = 'desc' then t.created_at  end desc nulls last,

			case when p.sort_field = 'updated_at'  and p.sort_order = 'asc'  then t.updated_at  end nulls last,
			case when p.sort_field = 'updated_at'  and p.sort_order = 'desc' then t.updated_at  end desc nulls last,

			t.created_at desc
		limit (select limit_val from params)
		offset (select offset_val from params);
	`

	arguments := []any{
		request.UploaderID,
		request.Visibility,
		request.SortField,
		request.SortOrder,
		request.Limit,
		request.Offset,
	}

	tracks, err := postgres.FetchMany[domain.Track](ctx, r.db, getListTracksSQL, arguments...)
	if err != nil {
		return nil, err
	}

	return &domain.ListTracksResponse{
		Tracks: tracks,
	}, nil
}

func (r *TracksRepository) UpdateTrack(ctx context.Context, request domain.UpdateTrackRequest) error {
	const updateTracksQL = `
		update tracks
		set
			title       = coalesce($2, title),
			subtitle    = coalesce($3, subtitle),
			description = coalesce($4, description),
			visibility  = coalesce($5::visibility, visibility),
			updated_at  = now()
		where id = $1;
	`

	arguments := []any{
		request.ID,
		request.Title,
		request.Subtitle,
		request.Description,
		request.Visibility,
	}

	_, err := postgres.FetchOne[domain.Track](ctx, r.db, updateTracksQL, arguments...)
	if err != nil {
		return err
	}

	return nil
}

func (r *TracksRepository) DeleteTrack(ctx context.Context, request domain.DeleteTrackRequest) error {
	const deleteTracksQL = `
		delete from tracks where id = $1;;
	`

	arguments := []any{
		request.ID,
	}

	affected, err := postgres.ExecAffected(ctx, r.db, deleteTracksQL, arguments...)
	if err != nil {
		return err
	}

	if affected == 0 {
		return postgres.ErrNotFound
	}

	return nil
}
