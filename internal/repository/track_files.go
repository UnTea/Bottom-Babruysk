package repository

import (
	"context"

	"github.com/untea/bottom_babruysk/internal/domain"
	"github.com/untea/bottom_babruysk/internal/repository/postgres"
)

type TrackFilesRepository struct {
	db *postgres.Client
}

func NewTrackFilesRepository(db *postgres.Client) *TrackFilesRepository {
	return &TrackFilesRepository{db: db}
}

func (r *TrackFilesRepository) CreateTrackFile(ctx context.Context, request domain.CreateTrackFileRequest) (*domain.CreateTrackFileResponse, error) {
	const createTrackFileSQL = `
		insert into track_files (track_id, 
		                         filename, 
		                         s3_key, 
		                         mime, 
		                         format, 
		                         codec, 
		                         bitrate, 
		                         sample_rate, 
		                         channels, 
		                         size, 
		                         duration, 
		                         checksum, 
		                         uploaded_at) 
		values ($1, 
		        $2, 
		        $3, 
		        $4, 
		        $5, 
		        $6, 
		        $7, 
		        $8, 
		        $9, 
		        $10, 
		        $11,
		        $12, 
		        $13)
		
		returning id;
	`

	arguments := []any{
		request.TrackID,
		request.Filename,
		request.S3Key,
		request.Mime,
		request.Format,
		request.Codec,
		request.Bitrate,
		request.SampleRate,
		request.Channels,
		request.Size,
		request.Duration,
		request.Checksum,
		request.UploadedAt,
	}

	trackFile, err := postgres.FetchOne[domain.TrackFile](ctx, r.db, createTrackFileSQL, arguments...)
	if err != nil {

		return nil, err

	}

	return &domain.CreateTrackFileResponse{
		ID: trackFile.ID,
	}, nil
}

func (r *TrackFilesRepository) GetTrackFile(ctx context.Context, request domain.GetTrackFileRequest) (*domain.GetTrackFileResponse, error) {
	const getTrackFileSQL = `
		select 
		    id, 
		    track_id, 
		    filename, 
		    s3_key, 
		    mime, 
		    format, 
		    codec,
		    bitrate, 
		    sample_rate, 
		    channels, 
		    size, 
		    duration, 
		    checksum,
		    created_at, 
		    updated_at, 
		    uploaded_at
		from track_files
		where id = $1 and track_id = $2;
	`

	arguments := []any{
		request.TrackID,
	}

	trackFile, err := postgres.FetchOne[domain.TrackFile](ctx, r.db, getTrackFileSQL, arguments...)
	if err != nil {
		return nil, err
	}

	return &domain.GetTrackFileResponse{
		TrackFile: trackFile,
	}, nil
}

func (r *TrackFilesRepository) ListTrackFiles(ctx context.Context, request domain.ListTrackFilesRequest) (*domain.ListTrackFilesResponse, error) {
	const listTrackFilesSQL = `
		with params as (
			select
				$1::uuid                                      as track_filter,
				coalesce(nullif(lower($2), ''), 'created_at') as sort_field,
				coalesce(nullif(lower($3), ''), 'desc')       as sort_order,
				greatest(coalesce($4, 50), 1)                 as limit_val,
				greatest(coalesce($5, 0), 0)                  as offset_val
		)
		select
			tf.id, 
			tf.track_id, 
			tf.filename, 
			tf.s3_key, 
			tf.mime, 
			tf.format, 
			tf.codec,
			tf.bitrate, 
			tf.sample_rate, 
			tf.channels, 
			tf.size, 
			tf.duration, 
			tf.checksum,
			tf.created_at, 
			tf.updated_at, 
			tf.uploaded_at
		from track_files as tf, params as p
		where tf.track_id = p.track_filter
		order by
			case when p.sort_field = 'filename'    and p.sort_order = 'asc'  then tf.filename    end nulls last,
			case when p.sort_field = 'filename'    and p.sort_order = 'desc' then tf.filename    end desc nulls last,

			case when p.sort_field = 'format'      and p.sort_order = 'asc'  then tf.format      end nulls last,
			case when p.sort_field = 'format'      and p.sort_order = 'desc' then tf.format      end desc nulls last,

			case when p.sort_field = 'codec'       and p.sort_order = 'asc'  then tf.codec       end nulls last,
			case when p.sort_field = 'codec'       and p.sort_order = 'desc' then tf.codec       end desc nulls last,

			case when p.sort_field = 'bitrate'     and p.sort_order = 'asc'  then tf.bitrate     end nulls last,
			case when p.sort_field = 'bitrate'     and p.sort_order = 'desc' then tf.bitrate     end desc nulls last,

			case when p.sort_field = 'sample_rate' and p.sort_order = 'asc'  then tf.sample_rate end nulls last,
			case when p.sort_field = 'sample_rate' and p.sort_order = 'desc' then tf.sample_rate end desc nulls last,

			case when p.sort_field = 'channels'    and p.sort_order = 'asc'  then tf.channels    end nulls last,
			case when p.sort_field = 'channels'    and p.sort_order = 'desc' then tf.channels    end desc nulls last,

			case when p.sort_field = 'size'        and p.sort_order = 'asc'  then tf.size        end  nulls last,
			case when p.sort_field = 'size'        and p.sort_order = 'desc' then tf.size        end desc nulls last,

			case when p.sort_field = 'duration'    and p.sort_order = 'asc'  then tf.duration    end nulls last,
			case when p.sort_field = 'duration'    and p.sort_order = 'desc' then tf.duration    end desc nulls last,

			case when p.sort_field = 'uploaded_at' and p.sort_order = 'asc'  then tf.uploaded_at end nulls last,
			case when p.sort_field = 'uploaded_at' and p.sort_order = 'desc' then tf.uploaded_at end desc nulls last,

			case when p.sort_field = 'created_at'  and p.sort_order = 'asc'  then tf.created_at  end nulls last,
			case when p.sort_field = 'created_at'  and p.sort_order = 'desc' then tf.created_at  end desc nulls last,

			case when p.sort_field = 'updated_at'  and p.sort_order = 'asc'  then tf.updated_at  end nulls last,
			case when p.sort_field = 'updated_at'  and p.sort_order = 'desc' then tf.updated_at  end desc nulls last,

			tf.created_at desc
		limit (select limit_val from params)
		offset (select offset_val from params);
	`

	arguments := []any{
		request.TrackID,
		request.SortField,
		request.SortOrder,
		request.Limit,
		request.Offset,
	}

	trackFiles, err := postgres.FetchMany[domain.TrackFile](ctx, r.db, listTrackFilesSQL, arguments...)
	if err != nil {
		return nil, err
	}

	return &domain.ListTrackFilesResponse{
		TrackFiles: trackFiles,
	}, nil
}

func (r *TrackFilesRepository) UpdateTrackFile(ctx context.Context, request domain.UpdateTrackFileRequest) error {
	const updateTrackFilesSQL = `
		update track_files
		set
			filename    = coalesce($3, filename),
			s3_key      = coalesce($4, s3_key),
			mime        = coalesce($5, mime),
			format      = coalesce($6::format, format),
			codec       = coalesce($7::codec, codec),
			bitrate     = coalesce($8, bitrate),
			sample_rate = coalesce($9, sample_rate),
			channels    = coalesce($10, channels),
			size        = coalesce($11, size),
			duration    = coalesce($12, duration),
			checksum    = coalesce($13, checksum),
			updated_at  = now()
		where id = $1 and track_id = $2;
	`

	arguments := []any{
		request.ID,
		request.TrackID,
		request.Filename,
		request.S3Key,
		request.Mime,
		request.Format,
		request.Codec,
		request.Bitrate,
		request.SampleRate,
		request.Channels,
		request.Size,
		request.Duration,
		request.Checksum,
	}

	_, err := postgres.FetchOne[domain.TrackFile](ctx, r.db, updateTrackFilesSQL, arguments...)
	if err != nil {
		return err
	}

	return err
}

func (r *TrackFilesRepository) DeleteTrackFile(ctx context.Context, request domain.DeleteTrackFileRequest) error {
	const deleteTrackFilesSQL = `
		delete from track_files where id = $1 and track_id = $2;
	`

	arguments := []any{
		request.ID,
		request.TrackID,
	}

	affected, err := postgres.ExecAffected(ctx, r.db, deleteTrackFilesSQL, arguments...)
	if err != nil {
		return err
	}

	if affected == 0 {
		return postgres.ErrNotFound
	}

	return nil
}
