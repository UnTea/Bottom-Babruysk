package flac

import (
	"io"
	"strings"
	"time"
)

// SeekPoint из SEEKTABLE (18 байт на точку).
type SeekPoint struct {
	SampleNumber uint64 // 0xFFFFFFFFFFFFFFFF placeholder
	Offset       uint64 // абсолютное смещение байтов относительно начала аудиопотока
	FrameSamples uint16 // количество сэмплов в аудиофрейме, на который указывает точка
}

// Picture соответствует метаблоку PICTURE.
type Picture struct {
	Type        uint32
	MIME        string
	Description string
	Width       uint32
	Height      uint32
	Depth       uint32
	Colors      uint32
	Data        []byte // Само изображение
}

// Application APPLICATION block (ID + raw payload).
type Application struct {
	ID   [4]byte
	Data []byte
}

// CueIndex CUESHEET entities (сокращённый, но совместимый набор).
type CueIndex struct {
	Offset uint64
	Number uint8
}

type CueTrack struct {
	Offset       uint64
	Number       uint8
	ISRC         string
	IsAudio      bool
	PreEmphasis  bool
	IndexEntries []CueIndex
}

type CueSheet struct {
	MediaCatalog string
	LeadIn       uint64
	IsCD         bool
	Tracks       []CueTrack
}

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
	MD5Signature  [16]byte

	// Derived
	Duration time.Duration

	// VORBIS_COMMENT
	Vendor    string
	Tags      map[string]string   // first-value convenience
	TagsMulti map[string][]string // full list

	// Additional blocks
	SeekTable    []SeekPoint
	Pictures     []Picture
	CueSheet     *CueSheet
	Applications []Application
}

// DecodeMetadata auto-detects and decodes metadata; now supports only FLAC.
func DecodeMetadata(r io.ReadSeeker) (Metadata, error) {
	if isFLAC(r) {
		if _, err := r.Seek(0, io.SeekStart); err != nil {
			return Metadata{}, err
		}

		return decodeFLACMetadata(r)
	}

	return Metadata{}, ErrUnsupportedFormat
}

func normalizeKey(k string) string {
	return strings.ToUpper(strings.TrimSpace(k))
}
