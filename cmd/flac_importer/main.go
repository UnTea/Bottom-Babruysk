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
	"github.com/untea/bottom_babruysk/internal/repository/pgerrors"
	"github.com/untea/bottom_babruysk/internal/repository/postgres"
)

const (
	// Жёстко известные свойства FLAC-файла.
	mimeFLAC   = "audio/flac"
	formatFLAC = "flac"
	codecFLAC  = "flac"
	flacExt    = ".flac"

	// Значения по умолчанию.
	defaultVisibility    = "private"
	defaultUploaderEmail = "imposter@sosal.cock"
)

var (
	// Ключи VorbisComment, которые нас интересуют в импорте.
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

// Importer здесь кэшируются уже найденные/созданные сущности (артисты/альбомы), чтобы уменьшить количество одинаковых
// select/insert при обходе каталога.
type Importer struct {
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
			ctx := context.Background()

			l, err := logger.New()
			if err != nil {
				return err
			}

			defer func() { _ = l.Sync() }()

			if flagPath == "" || flagDSN == "" {
				return fmt.Errorf("path and dsn are required")
			}

			dbConfiguration := postgres.Configuration{
				ConnectionString: flagDSN,
				Timeout:          30 * time.Second,
			}

			db, err := postgres.New(context.Background(), dbConfiguration)
			if err != nil {
				l.Fatal("failed to initialize db client", zap.Error(err))
			}

			defer db.Close()

			importer := &Importer{
				db:      db,
				logger:  l,
				artists: make(map[artistKey]*uuid.UUID),
				albums:  make(map[albumKey]*uuid.UUID),
			}

			uploaderID, err := importer.ensureUploaderUser(ctx, flagUploaderEmail)
			if err != nil {
				l.Fatal("failed to ensure uploader user", zap.Error(err))
			}

			importer.uploaderID = uploaderID

			return importer.RunAll(ctx, flagPath)
		},
	}

	root.Flags().StringVarP(&flagPath, "path", "p", "", "directory to scan for .flac files (required)")
	root.Flags().StringVarP(&flagDSN, "dsn", "d", "", "Postgres DSN (required)")
	root.Flags().StringVar(&flagUploaderEmail, "uploader-email", defaultUploaderEmail, "email for uploader user")

	_ = root.MarkFlagRequired("path")
	_ = root.MarkFlagRequired("dsn")

	err := root.Execute()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// RunAll основной цикл импорта. Ищем ".flac" файлы и обрабатываем каждый.
func (i *Importer) RunAll(ctx context.Context, scanPath string) error {
	files, err := collectFlacFiles(scanPath)
	if err != nil {
		return err
	}

	if len(files) == 0 {
		i.logger.Info("no .flac files found", zap.String("path", scanPath))
		return nil
	}

	i.logger.Info("found FLAC files", zap.Int("count", len(files)))

	var imported, skipped int

	for _, file := range files {
		skip, err := i.ingestOneFile(ctx, file)
		if err != nil {
			i.logger.Warn("ingest failed", zap.String("file", file), zap.Error(err))
			continue
		}

		if skip {
			skipped++
			continue
		}

		imported++
	}

	i.logger.Info("done", zap.Int("imported", imported), zap.Int("skipped", skipped))

	return nil
}

// collectFlacFiles рекурсивно находит все файлы с расширением ".flac".
func collectFlacFiles(root string) ([]string, error) {
	var out []string

	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		if strings.ToLower(filepath.Ext(d.Name())) == flacExt {
			out = append(out, path)
		}

		return nil
	})

	return out, err
}

// ensureUploaderUser гарантирует, что в таблице users есть запись с указанным email.
// Возвращает её id (существующий или только что созданный).
func (i *Importer) ensureUploaderUser(ctx context.Context, email string) (*uuid.UUID, error) {
	const (
		displayName = "FLAC importer"
		role        = "user"
	)

	const getUserSQL = `
		select id from users where email = $1;
	`

	getUserArguments := []any{
		email,
	}

	userID, err := postgres.FetchScalar[uuid.UUID](ctx, i.db, getUserSQL, getUserArguments...)
	if err == nil {
		return userID, nil
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

	userID, err = postgres.FetchScalar[uuid.UUID](ctx, i.db, createUserSQL, createUserArguments...)
	if err != nil {
		return &uuid.Nil, err
	}

	return userID, nil
}

// ingestOneFile бизнес-операция импорта одного FLAC-файла:
//
//  1. Читаем метаданные (FLAC),
//  2. По MD5 (из STREAMINFO) проверяем, не импортирован ли файл ранее,
//  3. Достаём теги (title/album/artist/year),
//  4. Обеспечиваем наличие артиста/альбома,
//  5. Ищем/создаём трек (по normalized_title + artist),
//  6. Создаём track_file (уникально по checksum).
//
// Возвращает skip=true, если файл уже существовал (по checksum).
func (i *Importer) ingestOneFile(ctx context.Context, path string) (bool, error) {
	file, err := os.Open(path)
	if err != nil {
		i.logger.Error("failed to open file", zap.String("path", path), zap.Error(err))
		return false, err
	}

	defer file.Close()

	metadata, err := flac.DecodeMetadata(file)
	if err != nil {
		i.logger.Error("failed to decode file metadata", zap.String("path", path), zap.Error(err))
		return false, err
	}

	// Проверяем уникальность файла по MD5 из STREAMINFO (md5 несжатых PCM-данных).
	exists, err := i.checkTrackFileExists(ctx, metadata.MD5Signature[:])
	if err != nil {
		return false, err
	}

	if exists {
		return true, nil
	}

	// Достаём человекочитаемые поля из тегов/пути.
	title, album, artist, yearOfRecording := i.extractTags(metadata, path)

	// Обеспечиваем артиста и альбом (upsert по normalized_*).
	artistID, err := i.getOrCreateArtist(ctx, artist)
	if err != nil {
		i.logger.Error("failed to get or create artist", zap.String("artist", artist), zap.Error(err))
		return false, err
	}

	_ = artistID // TODO: связать, когда появится M2M трека с артистом.

	albumID, err := i.getOrCreateAlbum(ctx, i.uploaderID, album, yearOfRecording)
	if err != nil {
		i.logger.Error("failed to get or create album", zap.String("album", album), zap.Error(err))
		return false, err
	}

	_ = albumID // TODO: связать с треком, когда будет явная схема/связь.

	// Убеждаемся, что логический трек (tracks) не дублируется:
	//	условие — normalized_title + конкретный artist.
	trackID, err := i.ensureTrack(ctx, artistID, title, metadata.Duration, yearOfRecording, ModTime(path))
	if err != nil {
		return false, err
	}

	// Создаём track_files
	filename := filepath.Base(path)

	_, err = i.insertTrackFile(ctx, trackID, metadata, filename, ModTime(path))
	if err != nil && !pgerrors.IsUniqueViolation(err) {
		return false, err
	}

	i.logger.Info("successfully created track file", zap.Any("file", filename))

	return false, nil
}

// getOrCreateArtist upsert артиста по normalized_display_name.
func (i *Importer) getOrCreateArtist(ctx context.Context, displayName string) (*uuid.UUID, error) {
	normalizedDisplayName := normalizeString(displayName)

	key := artistKey(normalizedDisplayName)
	if id, ok := i.artists[key]; ok {
		return id, nil
	}

	const upsertArtisSQL = `
        insert into artists (display_name, 
                             normalized_display_name, 
                             bio)
        values ($1, 
                $2, 
                $3)
        on conflict (normalized_display_name) do update
            set display_name = excluded.display_name
        returning id;
	`

	checkArtistArguments := []any{
		displayName,
		normalizedDisplayName,
		"",
	}

	artistID, err := postgres.FetchScalar[uuid.UUID](ctx, i.db, upsertArtisSQL, checkArtistArguments...)
	if err != nil {
		return nil, err
	}

	i.artists[key] = artistID

	return artistID, nil
}

// getOrCreateAlbum upsert альбома по normalized_title.
func (i *Importer) getOrCreateAlbum(ctx context.Context, ownerID *uuid.UUID, title string, releaseDate time.Time) (*uuid.UUID, error) {
	normalizedTitle := normalizeString(title)

	key := albumKey(strings.TrimSpace(strings.ToLower(title)))
	if id, ok := i.albums[key]; ok {
		return id, nil
	}

	const upsertAlbumSQL = `
        insert into albums (owner_id, 
							title, 
                            normalized_title, 
                            description, 
                            visibility, 
                            release_date)
        values ($1, 
                $2, 
                $3, 
                $4, 
                $5::visibility, 
                $6)
        on conflict (normalized_title) do update
            set title = excluded.title
        returning id;
    `

	upsertAlbumArguments := []any{
		ownerID,
		title,
		normalizedTitle,
		"",
		defaultVisibility,
		releaseDate,
	}

	albumID, err := postgres.FetchScalar[uuid.UUID](ctx, i.db, upsertAlbumSQL, upsertAlbumArguments...)
	if err != nil {
		return nil, err
	}

	i.albums[key] = albumID

	return albumID, nil
}

// checkTrackFileExists быстрая проверка по уникальному checksum (bytea unique).
func (i *Importer) checkTrackFileExists(ctx context.Context, checksum []byte) (bool, error) {
	const checkTrackFilesSQL = `
		select exists(select 1 from track_files where checksum = $1);
	`

	exists, err := postgres.FetchScalar[bool](ctx, i.db, checkTrackFilesSQL, checksum)
	if err != nil {
		return false, err
	}

	return *exists, err
}

// ensureTrack возвращает id трека, если он уже есть (artist + normalized_title), иначе создаёт новый и связывает его
// с артистом в track_artists.
func (i *Importer) ensureTrack(ctx context.Context, artistID *uuid.UUID, title string, duration time.Duration, createdAt, uploadedAt time.Time) (*uuid.UUID, error) {
	normalizedTitle := normalizeString(title)

	const checkTrackSQL = `
        select t.id
        from tracks as t
        join track_artists as ta on ta.track_id = t.id
        where ta.artist_id = $1
          and t.normalized_title = $2
        limit 1;
    `

	checkTrackArguments := []any{
		artistID,
		normalizedTitle,
	}

	existing, err := postgres.FetchOne[domain.Track](ctx, i.db, checkTrackSQL, checkTrackArguments...)
	if err == nil && existing != nil && existing.ID != nil && *existing.ID != uuid.Nil {
		return existing.ID, nil
	}

	if err != nil && !errors.Is(err, postgres.ErrNotFound) {
		return nil, err
	}

	const createTrackSQL = `
        insert into tracks (uploader_id, 
                            title, 
                            normalized_title,
                            subtitle, 
                            description, 
                            duration,
                            visibility, 
                            created_at, 
                            uploaded_at )
        values ($1, 
                $2, 
                $3, 
                $4, 
                $5, 
                $6::interval, 
                $7::visibility, 
                $8, 
                $9)
        returning id;
    `

	createTrackArguments := []any{
		i.uploaderID,
		title,
		normalizedTitle,
		"",
		"",
		duration,
		defaultVisibility,
		createdAt,
		uploadedAt,
	}

	trackID, err := postgres.FetchScalar[uuid.UUID](ctx, i.db, createTrackSQL, createTrackArguments...)
	if err != nil {
		return nil, err
	}

	const linkArtistToTrackSQL = `
		insert into track_artists (track_id, artist_id)
		values ($1, $2)
		on conflict (track_id, artist_id) do nothing;
	`

	linkArtistToTrackArguments := []any{
		trackID,
		artistID,
	}

	_, err = postgres.ExecAffected(ctx, i.db, linkArtistToTrackSQL, linkArtistToTrackArguments...)
	if err != nil {
		return nil, err
	}

	return trackID, nil
}

// insertTrackFile создаёт запись в track_files для данного трека.
func (i *Importer) insertTrackFile(ctx context.Context, trackID *uuid.UUID, meta flac.Metadata, filename string, uploadedAt time.Time) (*uuid.UUID, error) {
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
		                         uploaded_at
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
		        $12
		)
		returning id;
	`

	insertTrackFileArguments := []any{
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
	}

	trackFileID, err := postgres.FetchScalar[uuid.UUID](ctx, i.db, insertTrackFile, insertTrackFileArguments...)
	if err != nil {
		return nil, err
	}

	return trackFileID, nil
}

// extractTags достаёт основные поля из VorbisComment, с запасными значениями, если тегов нет (имя файла -> title,
// Unknown Album/Artist -> плейсхолдеры).
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

// normalizeString нормализация для полей, по которым хотим устойчивые уникальные ключи: trim -> lower -> оставить
// только буквы/цифры.
func normalizeString(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))

	var b strings.Builder

	for _, r := range s {
		if unicode.IsLetter(r) || unicode.IsNumber(r) {
			b.WriteRune(r)
		}
	}

	return b.String()
}

// asYearOrZero пытается распарсить год "YYYY" из свободного текста (DATE), возвращает time.Time{} (zero), если не
// удалось.
func asYearOrZero(s string) time.Time {
	s = strings.TrimSpace(s)
	if len(s) < 4 {
		return time.Time{}
	}

	parsedTime, err := time.Parse("2006", s[:4])
	if err != nil {
		return time.Time{}
	}

	return parsedTime
}

// firstTag берёт первый непустой элемент по ключу из Tags/TagsMulti.
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

// firstNonEmpty возвращает первый непустой строковый аргумент.
func firstNonEmpty(values ...string) string {
	for _, s := range values {
		if strings.TrimSpace(s) != "" {
			return strings.TrimSpace(s)
		}
	}

	return ""
}

// ModTime время последней модификации файла (UTC). Используем как uploaded_at/created_at при импорте.
func ModTime(path string) time.Time {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return time.Now().UTC()
	}

	return fileInfo.ModTime().UTC()
}

// computeKbps оценка битрейта по формуле: rate * bits * channels / 1000.
// Для lossless это теоретический пиковый битрейт PCM. (У FLAC фактический битрейт меньше — зависит от сжатия.)
func computeKbps(sampleRate, bitsPerSample, channels int) int {
	if sampleRate <= 0 || bitsPerSample <= 0 || channels <= 0 {
		return 0
	}

	return (sampleRate * bitsPerSample * channels) / 1000
}
