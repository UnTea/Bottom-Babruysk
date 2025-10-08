package flac

import (
	"errors"
	"os"
	"strings"
	"time"
)

// SeekPoint — запись из SEEKTABLE (18 байт на точку).
// Позволяет быстро искать по абсолютным сэмплам без полного декодирования.
type SeekPoint struct {
	SampleNumber uint64 // 0xFFFFFFFFFFFFFFFF placeholder (нет явной точки)
	Offset       uint64 // Абсолютное байтовое смещение от начала аудиопотока
	FrameSamples uint16 // Количество сэмплов в целевом аудиофрейме
}

// Picture соответствует метаблоку PICTURE (обложки, картинки и т.п.).
type Picture struct {
	Type        uint32
	MIME        string
	Description string
	Width       uint32
	Height      uint32
	Depth       uint32
	Colors      uint32
	Data        []byte // Само изображение (байтовый blob)
}

// Application блок APPLICATION (4-байтовый ID + произвольный payload).
type Application struct {
	ID   [4]byte
	Data []byte
}

// CueIndex индекс в пределах трека CUESHEET (offset/number).
type CueIndex struct {
	Offset uint64
	Number uint8
}

// CueTrack запись трека внутри CUESHEET (см. RFC).
type CueTrack struct {
	Offset       uint64
	Number       uint8
	ISRC         string
	IsAudio      bool
	PreEmphasis  bool
	IndexEntries []CueIndex
}

// CueSheet целиком блок CUESHEET.
type CueSheet struct {
	MediaCatalog string
	LeadIn       uint64
	IsCD         bool
	Tracks       []CueTrack
}

// Metadata собранные метаданные FLAC файла.
// Поле Size не гарантировано для произвольных io.Reader, но заполняется, если вход — файл (*os.File).
type Metadata struct {
	Format string
	Size   int64

	// STREAMINFO
	MinBlockSize  uint16
	MaxBlockSize  uint16
	MinFrameSize  uint32 // 24 бита
	MaxFrameSize  uint32 // 24 бита
	SampleRate    int
	Channels      int
	BitsPerSample int
	TotalSamples  uint64
	MD5Signature  [16]byte // MD5 от несжатых PCM-данных

	// Производное
	Duration time.Duration

	// VORBIS_COMMENT
	Vendor    string
	Tags      map[string]string   // Первое значение для каждого ключа
	TagsMulti map[string][]string // Все значения по ключу

	// Дополнительные блоки
	SeekTable    []SeekPoint
	Pictures     []Picture
	CueSheet     *CueSheet
	Applications []Application
}

// DecodeMetadata декодирует метаданные FLAC из файла.
// Если сигнатура не FLAC, тогда возвращается ErrUnsupportedFormat.
func DecodeMetadata(file *os.File) (Metadata, error) {
	metadata, err := decodeFLACMetadata(file)
	if errors.Is(err, ErrInvalidMetadata) {
		return Metadata{}, ErrUnsupportedFormat
	}

	fi, err := file.Stat()
	if err == nil {
		metadata.Size = fi.Size()
	} else {
		return metadata, err
	}

	return metadata, err
}

func normalizeKey(k string) string {
	return strings.ToUpper(strings.TrimSpace(k))
}
