package flac

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"strings"
	"time"
)

var (
	// Общие ошибки уровня парсера.
	ErrTagsNotFound      = errors.New("tags not found")
	ErrInvalidMetadata   = errors.New("invalid metadata")
	ErrUnsupportedFormat = errors.New("unsupported format")

	// Ошибки, связанные с конкретными шагами чтения.
	ErrReadFLACMagic       = errors.New("read flac magic")
	ErrReadMetadata        = errors.New("read metadata header")
	ErrBlockTooLarge       = errors.New("block too large")
	ErrReadMetadataPayload = errors.New("read metadata payload")
	ErrBadStreamInfo       = errors.New("bad stream info")
	ErrStreamInfoTooShort  = errors.New("stream info too short")
	ErrPictureTooLarge     = errors.New("picture too large")
)

// RFC 9639 FLAC metadata block types и сигнатура файла.
// Сигнатура "fLaC" строго в начале контейнера.
const (
	flacMagic = "fLaC"

	blockTypeStreamInfo    = 0 // STREAMINFO первый обязательный блок с базовыми параметрами потока
	blockTypePadding       = 1 // PADDING заполнитель, может быть произвольной длины
	blockTypeApplication   = 2 // APPLICATION произвольные данные конкретного приложения/энкодера
	blockTypeSeekTable     = 3 // SEEKTABLE таблица точек позиционирования по сэмплам
	blockTypeVorbisComment = 4 // VORBIS_COMMENT теги (vendor + N пар key=value)
	blockTypeCueSheet      = 5 // CUESHEET структура треков/индексов (CD layout)
	blockTypePicture       = 6 // PICTURE обложки/картинки
	// 7..126 зарезервировано
	// 127 недопустимо (invalid)
)

// decodeFLACMetadata парсит метаданные FLAC последовательно из io.Reader.
//
// Формат заголовка метаблока (RFC 9639):
//   - 1 байт - флаг "last" (старший бит) + тип блока (младшие 7 бит)
//   - 3 байта - длина payload (24-битное BE число)
//
// Требование к вызывающему коду: r должен стоять в начале файла (на сигнатуре).
// Если сигнатура не совпадает, тогда вернётся ErrInvalidMetadata (внешний код может сопоставить это к
// ErrUnsupportedFormat).
func decodeFLACMetadata(r io.Reader) (metadata Metadata, err error) {
	metadata.Format = "FLAC"
	metadata.Tags = make(map[string]string)
	metadata.TagsMulti = make(map[string][]string)

	// Чтение и проверка сигнатуры "fLaC".
	var magic [4]byte

	if _, err = io.ReadFull(r, magic[:]); err != nil {
		return Metadata{}, ErrReadFLACMagic
	}

	if string(magic[:]) != flacMagic {
		return Metadata{}, ErrInvalidMetadata
	}

	var haveStreamInfo bool

	for {
		// Заголовок метаблока (RFC 9639):
		// 1 байт: last-flag (1bit) + block-type (7bit)
		// 3 байта: length (24-bit BE)
		b0, err := readU8(r)
		if err != nil {
			return Metadata{}, ErrReadMetadata
		}

		isLast := (b0 & 0x80) != 0
		bType := b0 & 0x7F

		length, err := readU24BE(r)
		if err != nil {
			return Metadata{}, ErrReadMetadata
		}

		if length > 32<<20 { // 32 MiB sanity limit
			return Metadata{}, ErrBlockTooLarge
		}

		payload, err := readBytes(r, int(length))
		if err != nil {
			return Metadata{}, ErrReadMetadataPayload
		}

		switch bType {
		case blockTypeStreamInfo:
			// STREAMINFO 34 байта (минимум): базовые параметры потока.
			// Если его нет или он некорректен - дальнейшая логика (например длительность) не вычисляется.
			err := parseStreamInfo(&metadata, payload)
			if err != nil {
				return Metadata{}, ErrBadStreamInfo
			}

			haveStreamInfo = true

		case blockTypeVorbisComment:
			// VorbisComment строка vendor + список пар key=value (ключи case-insensitive, традиционно UPPERCASE).
			vendor, tags := parseVorbisComment(payload)
			if vendor != "" {
				metadata.Vendor = vendor
			}

			// TagsMulti копит все значения одного ключа (могут повторяться).
			for key, value := range tags {
				metadata.TagsMulti[key] = append(metadata.TagsMulti[key], value...)
			}

		case blockTypeSeekTable:
			// SEEKTABLE массив структур фиксированной длины (18 байт): сэмпл, байтовое смещение, размер фрейма.
			points, err := parseSeekTable(payload)
			if err == nil && len(points) > 0 {
				metadata.SeekTable = append(metadata.SeekTable, points...)
			}

		case blockTypePicture:
			// PICTURE может быть несколько, содержат MIME, размеры и сам байтовый blob.
			picture, err := parsePicture(bytes.NewReader(payload))
			if err == nil {
				metadata.Pictures = append(metadata.Pictures, *picture)
			}

		case blockTypeCueSheet:
			// CUESHEET CD-подобный трек-лист (лиды, индексы, pre-emphasis и прочие).
			cueSheet, err := parseCueSheet(bytes.NewReader(payload))
			if err == nil {
				metadata.CueSheet = cueSheet
			}

		case blockTypeApplication:
			// APPLICATION специфичный для энкодера блок (ID + произвольные данные).
			application, err := parseApplication(bytes.NewReader(payload))
			if err == nil {
				metadata.Applications = append(metadata.Applications, *application)
			}

		case blockTypePadding:
			// PADDING пустой блок, игнорируем.
		default:
			// Резерв - корректно пропускаем.
		}

		if isLast {
			break
		}
	}

	// Производное поле: длительность = total_samples / sample_rate.
	// Определяется только если STREAMINFO корректен и sample_rate > 0.
	if haveStreamInfo && metadata.SampleRate > 0 {
		seconds := float64(metadata.TotalSamples) / float64(metadata.SampleRate)
		metadata.Duration = time.Duration(seconds * float64(time.Second))
	}

	// Упрощённая мапа тегов: первый элемент по каждому ключу.
	for key, value := range metadata.TagsMulti {
		if len(value) > 0 {
			metadata.Tags[key] = value[0]
		}
	}

	if !haveStreamInfo && len(metadata.Tags)+len(metadata.TagsMulti) == 0 {
		return metadata, ErrTagsNotFound
	}

	return metadata, nil
}

// parseStreamInfo разбирает STREAMINFO (минимум 34 байта).
//
// Поля b[10:18] образуют 64-битное упакованное значение, где
// биты распределяются так (MSB->LSB):
//
//	[20 бит] sample_rate
//	[3  бита] (channels - 1)       физическое число каналов = значение + 1
//	[5  бит] (bits_per_sample - 1) глубина = значение + 1
//	[36 бит] total_samples         общее число сэмплов в потоке
func parseStreamInfo(metadata *Metadata, b []byte) error {
	if len(b) < 34 {
		return ErrStreamInfoTooShort
	}

	metadata.MinBlockSize = binary.BigEndian.Uint16(b[0:2])
	metadata.MaxBlockSize = binary.BigEndian.Uint16(b[2:4])

	// frame sizes — 24-битные BE (RFC).
	metadata.MinFrameSize = uint32(b[4])<<16 | uint32(b[5])<<8 | uint32(b[6])
	metadata.MaxFrameSize = uint32(b[7])<<16 | uint32(b[8])<<8 | uint32(b[9])

	// Упакованное поле с SampleRate/Channels/BitsPerSample/TotalSamples.
	packed := binary.BigEndian.Uint64(b[10:18])

	metadata.SampleRate = int((packed >> 44) & 0xFFFFF)       // 20 бит
	metadata.Channels = int(((packed >> 41) & 0x7) + 1)       // 3 бит (+1)
	metadata.BitsPerSample = int(((packed >> 36) & 0x1F) + 1) // 5 бит (+1)
	metadata.TotalSamples = packed & 0xFFFFFFFFF              // 36 бит

	// MD5 от несжатых PCM-данных (для детерминизма контента).
	copy(metadata.MD5Signature[:], b[18:34])

	return nil
}

// parseSeekTable читает SEEKTABLE из последовательности 18-байтовых записей.
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

// parseVorbisComment парсит VORBIS_COMMENT.
// Формат:
//
//	vendor_string (LPStringLE)
//	user_comment_list_length (u32LE)
//	N раз: comment_length(u32LE) + raw("KEY=VALUE")
//
// Ключи нормализуются (uppercase/trim), значения могут повторяться.
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

	for range n {
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

// parsePicture парсит PICTURE-блок.
// Поля: тип картинки, MIME, описание, размеры, глубина, количество цветов, длина и blob.
// Ограничиваем payload картинок 20 MiB (sane default).
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

// parseApplication читает APPLICATION-блок (фиксированный 4-байтовый ID + данные).
// Ограничение: не больше 4 MiB - отрубаем хвост, чтобы не засирать память.
func parseApplication(r io.Reader) (*Application, error) {
	var id [4]byte

	if _, err := io.ReadFull(r, id[:]); err != nil {
		return nil, err
	}

	rest, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	const maximumApplicationSize = 4 << 20 // 4Mb

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

// parseCueSheet парсит CUESHEET (CD-подобное описание дорожек).
// Поля включают: медиакаталог (ISBN), lead-in, признак CD, массив треков,
// у каждого трека - смещение, номер, ISRC, флаги (audio/non-audio, pre-emphasis),
// а также массив индексов (offset/number).
func parseCueSheet(r io.Reader) (*CueSheet, error) {
	// Медиа-каталог (128 байт, нуль-терминированная строка).
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

	// Зарезервированные 258 байт (RFC) - нужно пропустить.
	_, err = io.CopyN(io.Discard, r, 258)
	if err != nil {
		return nil, err
	}

	// Количество треков (1 байт).
	var tracksCount [1]byte

	_, err = io.ReadFull(r, tracksCount[:])
	if err != nil {
		return nil, err
	}

	n := int(tracksCount[0])
	tracks := make([]CueTrack, 0, n)

	for range n {
		offset, err := readU64BE(r)
		if err != nil {
			return nil, err
		}

		var num [1]byte

		_, err = io.ReadFull(r, num[:])
		if err != nil {
			return nil, err
		}

		isr, err := readBytes(r, 12) // ISRC код
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

		// Количество индексов.
		var indexesCount [1]byte

		_, err = io.ReadFull(r, indexesCount[:])
		if err != nil {
			return nil, err
		}

		ni := int(indexesCount[0])
		indexEntries := make([]CueIndex, 0, ni)

		for range ni {
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
