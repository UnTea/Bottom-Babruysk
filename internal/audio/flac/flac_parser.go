package flac

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"math"
	"strings"
	"time"
)

var (
	ErrTagsNotFound        = errors.New("tags not found")
	ErrInvalidMetadata     = errors.New("invalid metadata")
	ErrUnsupportedFormat   = errors.New("unsupported format")
	ErrReadFLACMagic       = errors.New("read flac magic")
	ErrReadMetadata        = errors.New("read metadata header")
	ErrBlockTooLarge       = errors.New("block too large")
	ErrReadMetadataPayload = errors.New("read metadata payload")
	ErrBadStreamInfo       = errors.New("bad stream info")
	ErrStreamInfoTooShort  = errors.New("stream info too short")
	ErrPictureTooLarge     = errors.New("picture too large")
)

// RFC 9639 FLAC metadata block types
const (
	flacMagic = "fLaC"

	blockTypeStreamInfo    = 0
	blockTypePadding       = 1
	blockTypeApplication   = 2
	blockTypeSeekTable     = 3
	blockTypeVorbisComment = 4
	blockTypeCueSheet      = 5
	blockTypePicture       = 6
	// 7...126 — зарезервировано
	// 127 — invalid
)

func isFLAC(r io.ReadSeeker) bool {
	_, err := r.Seek(0, io.SeekStart)
	if err != nil {
		return false
	}

	var b [4]byte

	_, err = io.ReadFull(r, b[:])
	if err != nil {
		return false
	}

	return string(b[:]) == flacMagic
}

func decodeFLACMetadata(r io.ReadSeeker) (metadata Metadata, err error) {
	metadata.Format = "FLAC"
	metadata.Tags = make(map[string]string)
	metadata.TagsMulti = make(map[string][]string)

	_, err = r.Seek(0, io.SeekStart) // "fLaC"
	if err != nil {
		return Metadata{}, err
	}

	var magic [4]byte

	_, err = io.ReadFull(r, magic[:])
	if err != nil {
		return Metadata{}, ErrReadFLACMagic
	}

	if string(magic[:]) != flacMagic {
		return Metadata{}, ErrInvalidMetadata
	}

	var haveStreamInfo bool

	for {
		var header [4]byte

		_, err = io.ReadFull(r, header[:])
		if err != nil {
			return Metadata{}, ErrReadMetadata
		}

		isLast := (header[0] & 0x80) != 0
		bType := header[0] & 0x7F
		length := int(binary.BigEndian.Uint32(append([]byte{0}, header[1], header[2], header[3])))

		if length < 0 || length > 32<<20 { // 32MB
			return Metadata{}, ErrBlockTooLarge
		}

		payload, err := readBytes(r, length)
		if err != nil {
			return Metadata{}, ErrReadMetadataPayload
		}

		switch bType {
		case blockTypeStreamInfo:
			err := parseStreamInfo(&metadata, payload)
			if err != nil {
				return Metadata{}, ErrBadStreamInfo
			}

			haveStreamInfo = true

		case blockTypeVorbisComment:
			vendor, tags := parseVorbisComment(payload)
			if vendor != "" {
				metadata.Vendor = vendor
			}

			for i, tag := range tags {
				metadata.TagsMulti[i] = append(metadata.TagsMulti[i], tag...)
			}

		case blockTypeSeekTable:
			points, err := parseSeekTable(payload)
			if err == nil && len(points) > 0 {
				metadata.SeekTable = append(metadata.SeekTable, points...)
			}

		case blockTypePicture:
			picture, err := parsePicture(bytes.NewReader(payload))
			if err == nil {
				metadata.Pictures = append(metadata.Pictures, *picture)
			}

		case blockTypeCueSheet:
			cueSheet, err := parseCueSheet(bytes.NewReader(payload))
			if err == nil {
				metadata.CueSheet = cueSheet
			}

		case blockTypeApplication:
			application, err := parseApplication(bytes.NewReader(payload))
			if err == nil {
				metadata.Applications = append(metadata.Applications, *application)
			}

		case blockTypePadding:
		default:
			// reserved/unknown
		}

		if isLast {
			break
		}
	}

	current, err := r.Seek(0, io.SeekCurrent)
	if err == nil {
		end, err := r.Seek(0, io.SeekEnd)
		if err == nil {
			metadata.Size = end
			_, _ = r.Seek(current, io.SeekStart)
		}
	}

	if haveStreamInfo && metadata.SampleRate > 0 {
		seconds := float64(metadata.TotalSamples) / float64(metadata.SampleRate)
		if !math.IsNaN(seconds) && !math.IsInf(seconds, 0) {
			metadata.Duration = time.Duration(seconds * float64(time.Second))
		}
	}

	for i, tag := range metadata.TagsMulti {
		if len(tag) > 0 {
			metadata.Tags[i] = tag[0]
		}
	}

	if !haveStreamInfo && len(metadata.Tags)+len(metadata.TagsMulti) == 0 {
		return metadata, ErrTagsNotFound
	}

	return metadata, nil
}

func parseStreamInfo(metadata *Metadata, b []byte) error {
	if len(b) < 34 {
		return ErrStreamInfoTooShort
	}

	metadata.MinBlockSize = binary.BigEndian.Uint16(b[0:2])
	metadata.MaxBlockSize = binary.BigEndian.Uint16(b[2:4])
	metadata.MinFrameSize = uint32(b[4])<<16 | uint32(b[5])<<8 | uint32(b[6])
	metadata.MaxFrameSize = uint32(b[7])<<16 | uint32(b[8])<<8 | uint32(b[9])

	packed := binary.BigEndian.Uint64(b[10:18])

	metadata.SampleRate = int((packed >> 44) & 0xFFFFF)       // 20 бит
	metadata.Channels = int(((packed >> 41) & 0x7) + 1)       // 3 бит +1
	metadata.BitsPerSample = int(((packed >> 36) & 0x1F) + 1) // 5 бит +1
	metadata.TotalSamples = packed & 0xFFFFFFFFF              // 36 бит

	copy(metadata.MD5Signature[:], b[18:34])

	return nil
}

func parseSeekTable(b []byte) ([]SeekPoint, error) {
	const entry = 18

	result := make([]SeekPoint, 0, len(b)/entry)

	for offset := 0; offset+entry <= len(b); offset += entry {
		seekPoint := SeekPoint{
			SampleNumber: binary.BigEndian.Uint64(b[offset+0 : offset+8]),
			Offset:       binary.BigEndian.Uint64(b[offset+8 : offset+16]),
			FrameSamples: binary.BigEndian.Uint16(b[offset+16 : offset+18]),
		}

		result = append(result, seekPoint)
	}

	return result, nil
}

func parseVorbisComment(b []byte) (vendor string, result map[string][]string) {
	result = make(map[string][]string)
	buffer := bytes.NewReader(b)

	ven, err := readLPStringLE(buffer)
	if err != nil {
		return "", result
	}

	vendor = ven

	n, err := readU32LE(buffer)
	if err != nil {
		return vendor, result
	}

	for i := uint32(0); i < n; i++ {
		l, err := readU32LE(buffer)
		if err != nil {
			break
		}

		if l == 0 {
			continue
		}

		raw, err := readBytes(buffer, int(l))
		if err != nil {
			break
		}

		kv := strings.SplitN(string(raw), "=", 2)
		if len(kv) != 2 {
			continue
		}

		k := normalizeKey(kv[0])
		v := strings.TrimSpace(kv[1])

		result[k] = append(result[k], v)
	}

	return vendor, result
}

func parsePicture(r io.Reader) (*Picture, error) {
	pictureType, err := readU32BE(r)
	if err != nil {
		return nil, err
	}

	mLen, err := readU32BE(r)
	if err != nil {
		return nil, err
	}

	mimeB, err := readBytes(r, int(mLen))
	if err != nil {
		return nil, err
	}

	dLen, err := readU32BE(r)
	if err != nil {
		return nil, err
	}

	descB, err := readBytes(r, int(dLen))
	if err != nil {
		return nil, err
	}

	width, err := readU32BE(r)
	if err != nil {
		return nil, err
	}

	height, err := readU32BE(r)
	if err != nil {
		return nil, err
	}

	depth, err := readU32BE(r)
	if err != nil {
		return nil, err
	}

	colors, err := readU32BE(r)
	if err != nil {
		return nil, err
	}

	dataLength, err := readU32BE(r)
	if err != nil {
		return nil, err
	}

	const maximumPictureSuze = 20 << 20 // 20Mb

	if dataLength > maximumPictureSuze {
		return nil, ErrPictureTooLarge
	}

	data, err := readBytes(r, int(dataLength))
	if err != nil {
		return nil, err
	}

	return &Picture{
		Type:        pictureType,
		MIME:        string(mimeB),
		Description: string(descB),
		Width:       width,
		Height:      height,
		Depth:       depth,
		Colors:      colors,
		Data:        data,
	}, nil
}

func parseApplication(r io.Reader) (*Application, error) {
	var id [4]byte

	if _, err := io.ReadFull(r, id[:]); err != nil {
		return nil, err
	}

	rest, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	const maximumApplicationSize = 4 << 20

	if len(rest) > maximumApplicationSize {
		return &Application{
			ID:   id,
			Data: rest[:maximumApplicationSize],
		}, nil
	}

	return &Application{
		ID:   id,
		Data: rest,
	}, nil
}

func parseCueSheet(r io.Reader) (*CueSheet, error) {
	// Медиа каталог (128 байт)
	mediaCatalog, err := readBytes(r, 128)
	if err != nil {
		return nil, err
	}

	leadIn, err := readU64BE(r)
	if err != nil {
		return nil, err
	}

	var flags [1]byte

	_, err = io.ReadFull(r, flags[:])
	if err != nil {
		return nil, err
	}

	isCD := (flags[0] & 0x80) != 0

	// Зарезервированные 258 байт.
	_, err = io.CopyN(io.Discard, r, 258)
	if err != nil {
		return nil, err
	}

	// Количество треков.
	var tcount [1]byte

	_, err = io.ReadFull(r, tcount[:])
	if err != nil {
		return nil, err
	}

	n := int(tcount[0])

	tracks := make([]CueTrack, 0, n)

	for i := 0; i < n; i++ {
		offset, err := readU64BE(r)
		if err != nil {
			return nil, err
		}

		var num [1]byte

		_, err = io.ReadFull(r, num[:])
		if err != nil {
			return nil, err
		}

		isr, err := readBytes(r, 12)
		if err != nil {
			return nil, err
		}

		var tflags [1]byte

		_, err = io.ReadFull(r, tflags[:])
		if err != nil {
			return nil, err
		}

		isAudio := (tflags[0] & 0x80) == 0     // audio(0)/non-audio(1)
		preEmphasis := (tflags[0] & 0x40) != 0 // bit6=pre-emphasis

		// Зарезервированные 13 байт.
		_, err = io.CopyN(io.Discard, r, 13)
		if err != nil {
			return nil, err
		}

		var indexesCount [1]byte

		_, err = io.ReadFull(r, indexesCount[:])
		if err != nil {
			return nil, err
		}

		ni := int(indexesCount[0])
		indexEntries := make([]CueIndex, 0, ni)

		for j := 0; j < ni; j++ {
			cueOffset, err := readU64BE(r)
			if err != nil {
				return nil, err
			}

			var cueNumber [1]byte

			_, err = io.ReadFull(r, cueNumber[:])
			if err != nil {
				return nil, err
			}

			// Зарезервированные 3 байта.
			_, err = io.CopyN(io.Discard, r, 3)
			if err != nil {
				return nil, err
			}

			indexEntries = append(indexEntries, CueIndex{
				Offset: cueOffset,
				Number: cueNumber[0],
			})
		}

		tracks = append(tracks, CueTrack{
			Offset:       offset,
			Number:       num[0],
			ISRC:         strings.TrimRight(string(isr), "\x00 "),
			IsAudio:      isAudio,
			PreEmphasis:  preEmphasis,
			IndexEntries: indexEntries,
		})
	}

	return &CueSheet{
		MediaCatalog: strings.TrimRight(string(mediaCatalog), "\x00 "),
		LeadIn:       leadIn,
		IsCD:         isCD,
		Tracks:       tracks,
	}, nil
}
