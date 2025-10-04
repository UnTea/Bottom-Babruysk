package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
	"unicode"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/untea/bottom_babruysk/internal/audio/flac"
	"github.com/untea/bottom_babruysk/internal/domain"
	"github.com/untea/bottom_babruysk/internal/logger"
	"github.com/untea/bottom_babruysk/internal/repository/postgres"
)

const (
	mimeFLAC             = "audio/flac"
	formatFLAC           = "flac"
	codecFLAC            = "flac"
	defaultVisibility    = "private"
	defaultUploaderEmail = "imposter@sosal.cock"
)

var (
	tagAlbum       = "ALBUM"
	tagArtist      = "ARTIST"
	tagAlbumArtist = "ALBUMARTIST"
	tagTitle       = "TITLE"
	tagTrackNumber = "TRACKNUMBER"
	tagDate        = "DATE"
)

type (
	artistKey string
	albumKey  string
)

type Importer struct {
	ctx        context.Context
	db         *postgres.Client
	uploaderID *uuid.UUID
	logger     *zap.Logger

	artists map[artistKey]*uuid.UUID
	albums  map[albumKey]*uuid.UUID
}

func main() {
	var (
		flagPath          string
		flagDSN           string
		flagUploaderEmail string
	)

	root := &cobra.Command{
		Use:   "importer",
		Short: "Import FLAC files from a directory into DB",
		RunE: func(cmd *cobra.Command, args []string) error {
			l, err := logger.New()
			if err != nil {
				return err
			}

			defer func() {
				_ = l.Sync()
			}()

			if flagPath == "" || flagDSN == "" {
				return fmt.Errorf("path and dsn are required")
			}

			dbCfg := postgres.Configuration{
				ConnectionString: flagDSN,
				Timeout:          30 * time.Second,
			}

			dbClient, err := postgres.New(context.Background(), dbCfg)
			if err != nil {
				l.Fatal("failed to initialize db client", zap.Error(err))
			}

			defer dbClient.Close()

			importer := &Importer{
				ctx:     context.Background(),
				db:      dbClient,
				logger:  l,
				artists: make(map[artistKey]*uuid.UUID),
				albums:  make(map[albumKey]*uuid.UUID),
			}

			uploaderID, err := importer.ensureUploaderUser(flagUploaderEmail)
			if err != nil {
				l.Fatal("failed to ensure uploader user", zap.Error(err))
			}

			importer.uploaderID = uploaderID

			files, err := collectFlacFiles(flagPath)
			if err != nil {
				return err
			}

			if len(files) == 0 {
				l.Info("no .flac files found", zap.String("path", flagPath))
				return nil
			}

			l.Info("found FLAC files", zap.Int("count", len(files)))

			var imported, skipped int

			for _, file := range files {
				skip, err := importer.ingestOneFile(file)
				if err != nil {
					l.Warn("ingest failed", zap.String("file", file), zap.Error(err))
					continue
				}

				if skip {
					skipped++
					continue
				}

				imported++
			}

			l.Info("done", zap.Int("imported", imported), zap.Int("skipped", skipped))

			return nil
		},
	}

	root.Flags().StringVarP(&flagPath, "path", "p", "", "directory to scan for .flac files (required)")
	root.Flags().StringVarP(&flagDSN, "dsn", "d", "", "Postgres DSN (required)")
	root.Flags().StringVar(&flagUploaderEmail, "uploader-email", defaultUploaderEmail, "email for uploader user")

	_ = root.MarkFlagRequired("path")
	_ = root.MarkFlagRequired("dsn")

	if err := root.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func collectFlacFiles(root string) ([]string, error) {
	var out []string

	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		if strings.ToLower(filepath.Ext(d.Name())) == ".flac" {
			out = append(out, path)
		}

		return nil
	})

	return out, err
}

func (i *Importer) ensureUploaderUser(email string) (*uuid.UUID, error) {
	const displayName = "FLAC importer"
	const role = "user"

	const getUserSQL = `
		select id from users where email = $1 limit 1;
	`

	getUserArguments := []any{
		email,
	}

	user, err := postgres.FetchOne[domain.User](i.ctx, i.db, getUserSQL, getUserArguments...)
	if err == nil && user != nil {
		return user.ID, nil
	}

	const createUserSQL = `
		insert into users (email, 
		                   password_hash, 
		                   display_name, 
		                   role)
		values ($1, 
		        '', 
		        $2, 
		        $3)
		returning id;
	`

	createUserArguments := []any{
		email,
		displayName,
		role,
	}

	user, err = postgres.FetchOne[domain.User](i.ctx, i.db, createUserSQL, createUserArguments...)
	if err != nil {
		return &uuid.Nil, err
	}

	return user.ID, nil
}

func (i *Importer) ingestOneFile(path string) (bool, error) {
	files, err := os.Open(path)
	if err != nil {
		i.logger.Error("failed to open file", zap.String("path", path), zap.Error(err))
		return false, err
	}

	defer files.Close()

	meta, err := flac.DecodeMetadata(files)
	if err != nil {
		i.logger.Error("failed to decode file metadata", zap.String("path", path), zap.Error(err))
		return false, err
	}

	exists, err := i.checkChecksumExists(meta.MD5Signature[:])
	if err != nil {
		return false, err
	}

	if exists {
		return true, nil
	}

	title, album, artist, yearOfRecording := i.extractTags(meta, path)

	artistID, err := i.getOrCreateArtist(artist)
	if err != nil {
		i.logger.Error("failed to get or create artist", zap.String("artist", artist), zap.Error(err))
		return false, err
	}

	_ = artistID // TODO: связать, когда появится M2M трека с артистом.

	albumID, err := i.getOrCreateAlbum(i.uploaderID, album, yearOfRecording)
	if err != nil {
		i.logger.Error("failed to get or create album", zap.String("album", album), zap.Error(err))
		return false, err
	}

	_ = albumID // TODO: связать с треком, когда будет явная схема/связь.

	trackID, err := i.ensureTrack(artistID, title, meta.Duration, yearOfRecording, fileMimeTime(path))
	if err != nil {
		return false, err
	}

	filename := filepath.Base(path)

	_, err = i.insertTrackFile(trackID, meta, filename, fileMimeTime(path))
	if err != nil && !isUniqueViolation(err) {
		return false, err
	}

	i.logger.Info("successfully created track file", zap.Any("file", filename))

	return false, nil
}

func (i *Importer) getOrCreateArtist(displayName string) (*uuid.UUID, error) {
	normalizedDisplayName := normalizeName(displayName)

	key := artistKey(normalizedDisplayName)
	if id, ok := i.artists[key]; ok {
		return id, nil
	}

	const checkArtistSQL = `
		select id from artists where lower(normalized_display_name) = lower($1);
	`

	checkArtistArguments := []any{
		normalizedDisplayName,
	}

	artists, err := postgres.FetchOne[domain.Artist](i.ctx, i.db, checkArtistSQL, checkArtistArguments...)
	if err == nil && artists != nil {
		i.artists[key] = artists.ID

		return artists.ID, err
	}

	const createArtistSQL = `
		insert into artists (display_name, normalized_display_name, bio) values ($1, $2, $3) returning id
	`

	createArtistArguments := []any{
		displayName,
		normalizedDisplayName,
		"",
	}

	artist, err := postgres.FetchOne[domain.Artist](i.ctx, i.db, createArtistSQL, createArtistArguments...)
	if err != nil {
		return nil, err
	}

	i.artists[key] = artist.ID

	return artist.ID, nil
}

func (i *Importer) getOrCreateAlbum(ownerID *uuid.UUID, title string, createdAt time.Time) (*uuid.UUID, error) {
	key := albumKey(strings.TrimSpace(strings.ToLower(title)))
	if id, ok := i.albums[key]; ok {
		return id, nil
	}

	const checkAlbumSQL = `
		select id from albums where owner_id = $1 and lower(title) = lower($2);
	`

	checkAlbumArguments := []any{
		ownerID,
		title,
	}

	album, err := postgres.FetchOne[domain.Album](i.ctx, i.db, checkAlbumSQL, checkAlbumArguments...)
	if err == nil && album != nil {
		i.albums[key] = album.ID

		return album.ID, err
	}

	const createAlbumSQL = `
		insert into albums (owner_id, 
		                    title, 
		                    description, 
		                    visibility,
		                    release_date) 
		values ($1, 
		        $2, 
		        $3, 
		        $4::visibility,
		        $5) 
		returning id;
	`

	createAlbumArguments := []any{
		ownerID,
		title,
		"",
		defaultVisibility,
		createdAt,
	}

	album, err = postgres.FetchOne[domain.Album](i.ctx, i.db, createAlbumSQL, createAlbumArguments...)
	if err != nil {
		return nil, err
	}

	i.albums[key] = album.ID

	return album.ID, nil
}

func (i *Importer) checkChecksumExists(checksum []byte) (bool, error) {
	const checkTrackFilesSQL = `select id from track_files where checksum = $1;`

	trackFile, err := postgres.FetchOne[domain.TrackFile](i.ctx, i.db, checkTrackFilesSQL, checksum)
	if err == nil && trackFile != nil && trackFile.ID != nil && *trackFile.ID != uuid.Nil {
		return true, nil
	}

	if errors.Is(err, postgres.ErrNotFound) {
		return false, nil
	}

	return false, err
}

func (i *Importer) ensureTrack(artistID *uuid.UUID, title string, duration time.Duration, createdAt, uploadedAt time.Time) (*uuid.UUID, error) {
	const checkTrackSQL = `
		select t.id
		from tracks t
		join track_artists ta on ta.track_id = t.id
		where ta.artist_id = $1
		  and lower(t.title) = lower($2)
	`

	existing, err := postgres.FetchOne[domain.Track](i.ctx, i.db, checkTrackSQL, artistID, title)
	if err == nil && existing != nil && existing.ID != nil && *existing.ID != uuid.Nil {
		return existing.ID, nil
	}

	if err != nil && !errors.Is(err, postgres.ErrNotFound) {
		return nil, err
	}

	const createTrackSQL = `
		insert into tracks (uploader_id, 
		                    title, 
		                    subtitle, 
		                    description,
		                    duration, 
		                    visibility, 
		                    created_at, 
		                    uploaded_at
		)
		values ($1,
		        $2, 
		        $3, 
		        $4, 
		        $5::interval, 
		        $6::visibility, 
		        $7, 
		        $8)
		returning id;
	`

	createTrackArguments := []any{
		i.uploaderID,
		title,
		"",
		"",
		duration,
		defaultVisibility,
		createdAt,
		uploadedAt,
	}

	track, err := postgres.FetchOne[domain.Track](i.ctx, i.db, createTrackSQL, createTrackArguments...)
	if err != nil {
		return nil, err
	}

	const linkArtistToTrackSQL = `
		insert into track_artists (track_id, artist_id)
		values ($1, $2)
		on conflict (track_id, artist_id) do nothing;
	`

	linkArtistToTrackArguments := []any{
		track.ID,
		artistID,
	}

	_, err = postgres.ExecAffected(i.ctx, i.db, linkArtistToTrackSQL, linkArtistToTrackArguments...)
	if err != nil {
		return nil, err
	}

	return track.ID, nil
}

func (i *Importer) insertTrackFile(trackID *uuid.UUID, meta flac.Metadata, filename string, uploadedAt time.Time) (*uuid.UUID, error) {
	bitrate := computeKbps(meta.SampleRate, meta.BitsPerSample, meta.Channels)

	const insertTrackFile = `
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
		                         uploaded_at, 
		                         created_at
		)
		values ($1, 
		        $2, 
		        null, 
		        $3, 
		        $4::format, 
		        $5::codec,
		        $6, 
		        $7, 
		        $8, 
		        $9, 
		        $10::interval,
		        $11, 
		        $12, 
		        now()
		)
		
		returning id;
	`

	tf, err := postgres.FetchOne[domain.TrackFile](
		i.ctx, i.db, insertTrackFile,
		trackID,
		filename,
		mimeFLAC,
		formatFLAC,
		codecFLAC,
		bitrate,
		meta.SampleRate,
		meta.Channels,
		meta.Size,
		meta.Duration,
		meta.MD5Signature[:],
		uploadedAt,
	)
	if err != nil {
		return nil, err
	}

	return tf.ID, nil
}

func (i *Importer) extractTags(meta flac.Metadata, path string) (title, album, artist string, year time.Time) {
	title = firstTag(meta, tagTitle)
	album = firstTag(meta, tagAlbum)
	artist = firstNonEmpty(firstTag(meta, tagArtist), firstTag(meta, tagAlbumArtist))

	if strings.TrimSpace(title) == "" {
		title = strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
	}

	if strings.TrimSpace(album) == "" {
		album = "Unknown Album"
	}

	if strings.TrimSpace(artist) == "" {
		artist = "Unknown Artist"
	}

	year = asYearOrZero(firstTag(meta, tagDate))

	return
}

func normalizeName(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))

	var b strings.Builder

	for _, r := range s {
		if unicode.IsLetter(r) || unicode.IsNumber(r) {
			b.WriteRune(r)
		}
	}

	return b.String()
}

func asYearOrZero(s string) time.Time {
	s = strings.TrimSpace(s)
	if len(s) < 4 {
		return time.Time{}
	}

	t, err := time.Parse("2006", s[:4])
	if err != nil {
		return time.Time{}
	}

	return t
}

func nullableDate(t time.Time) any {
	if t.IsZero() {
		return nil
	}

	return t
}

func isUniqueViolation(err error) bool {
	if err == nil {
		return false
	}

	msg := strings.ToLower(err.Error())

	return strings.Contains(msg, "duplicate key value") || strings.Contains(msg, "unique constraint")
}

func firstTag(metadata flac.Metadata, key string) string {
	if v, ok := metadata.Tags[key]; ok && strings.TrimSpace(v) != "" {
		return strings.TrimSpace(v)
	}

	if vv, ok := metadata.TagsMulti[key]; ok && len(vv) > 0 {
		for _, s := range vv {
			s = strings.TrimSpace(s)
			if s != "" {
				return s
			}
		}
	}

	return ""
}

func firstNonEmpty(values ...string) string {
	for _, s := range values {
		if strings.TrimSpace(s) != "" {
			return strings.TrimSpace(s)
		}
	}

	return ""
}

func fileMimeTime(path string) time.Time {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return time.Now().UTC()
	}

	return fileInfo.ModTime().UTC()
}

func computeKbps(sampleRate, bitsPerSample, channels int) int {
	if sampleRate <= 0 || bitsPerSample <= 0 || channels <= 0 {
		return 0
	}

	return (sampleRate * bitsPerSample * channels) / 1000
}
