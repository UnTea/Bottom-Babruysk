package audio

import (
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/untea/bottom_babruysk/internal/domain"
)

// TrackFileMetadata сжатое представление ключевой метаинформации, извлечённой из FLAC-файла. Эта структура
// приблизительно отражает подмножество полей domain.TrackFile, но намеренно отделена, чтобы не связывать parser жёстко
// с деталями хранения.
type TrackFileMetadata struct {
	// Path абсолютный или относительный путь к файлу, из которого извлечена метаинформация.
	Path string
	// Filename только базовое имя файла без директорий. Удобно иметь его под рукой; при этом вызывающая сторона всё
	// равно может отдельно сохранять полный путь или ключ в хранилище объектов.
	Filename string
	// Format всегда domain.FormatFLAC для FLAC-файлов. Держим его здесь, чтобы при сохранении не вычислять значение
	// повторно.
	Format domain.Format
	// Codec всегда domain.CodecFLAC для FLAC-файлов.
	Codec domain.Codec
	// SampleRate частота дискретизации в Гц (например 44100).
	SampleRate int
	// Channels количество каналов в аудиопотоке.
	Channels int
	// BitsPerSample разрядность PCM. Частые значения: 16 или 24.
	BitsPerSample int
	// Bitrate упрощённый расчёт сырого аудиобитрейта (кбит/с):
	// SampleRate * BitsPerSample * Channels / 1000.
	Bitrate int
	// TotalSamples общее число PCM-сэмплов. В паре с SampleRate позволяет вычислить длительность.
	TotalSamples uint64
	// Duration длительность трека (time.Duration). Вычисляется из TotalSamples и SampleRate.
	Duration time.Duration
	// Size размер файла на диске в байтах.
	Size int64
	// MD5Signature 16-байтная MD5-подпись не кодированных PCM-данных из блока STREAMINFO FLAC.
	MD5Signature [16]byte
}

// ParseFlacFile открывает файл по пути и пытается разобрать метаблок STREAMINFO, чтобы извлечь базовые параметры
// потока FLAC. Если файл не начинается с магической последовательности "fLaC" либо не содержит валидного STREAMINFO,
// будет возвращена ошибка. Некритичные проблемы (например, неожиданный порядок метаблоков) допускаются, пока удаётся
// найти и разобрать STREAMINFO.
func ParseFlacFile(filePath string) (*TrackFileMetadata, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("open flac file: %w", err)
	}

	defer f.Close()

	// Читаем и валидируем магическую сигнатуру FLAC.
	// Файл обязан начинаться с ASCII-последовательности "fLaC".
	var header [4]byte
	if _, err := io.ReadFull(f, header[:]); err != nil {
		return nil, fmt.Errorf("reading flac header: %w", err)
	}

	if string(header[:]) != "fLaC" {
		return nil, fmt.Errorf("invalid flac signature: %q", string(header[:]))
	}

	var streamInfoData []byte
	for {
		var h [4]byte

		if _, err := io.ReadFull(f, h[:]); err != nil {
			return nil, fmt.Errorf("reading metadata header: %w", err)
		}

		isLast := (h[0] & 0x80) != 0
		blockType := h[0] & 0x7F

		// Длина метаблока беззнаковое 24-битное число (big-endian),
		// записанное в следующих трёх байтах.
		blockLen := (uint32(h[1]) << 16) | (uint32(h[2]) << 8) | uint32(h[3])

		data := make([]byte, blockLen)
		if _, err := io.ReadFull(f, data); err != nil {
			return nil, fmt.Errorf("reading metadata block: %w", err)
		}

		if blockType == 0 { // STREAMINFO
			streamInfoData = data
			break
		}

		if isLast {
			break
		}
	}

	if len(streamInfoData) < 34 {
		return nil, fmt.Errorf("invalid STREAMINFO block length: %d", len(streamInfoData))
	}

	// Согласно спецификации FLAC (https://xiph.org/flac/documentation.html), блок STREAMINFO содержит (в байтах):
	//  0:1  минимальный размер блока (uint16)
	//  2:3  максимальный размер блока (uint16)
	//  4:6  минимальный размер фрейма (24 бита)
	//  7:9  максимальный размер фрейма (24 бита)
	// 10:17 sample rate, channels, bits per sample, total samples (упаковано)
	// 18:33 MD5-подпись
	// Мы игнорируем размеры блоков/фреймов и читаем 8 байт с offset 10.
	packed := binary.BigEndian.Uint64(streamInfoData[10:18])
	sampleRate := uint32((packed >> 44) & 0xFFFFF) // 20 бит на sample rate
	channels := uint8((packed>>41)&0x7) + 1        // 3 бита + 1
	bitsPerSample := uint8((packed>>36)&0x1F) + 1  // 5 бит + 1
	totalSamples := packed & 0xFFFFFFFFF           // младшие 36 бит

	var md5 [16]byte

	copy(md5[:], streamInfoData[18:34])

	// Узнаём размер файла.
	fi, err := f.Stat()
	if err != nil {
		return nil, fmt.Errorf("stat flac file: %w", err)
	}

	// Считаем простой битрейт. Важно: FLAC VBR кодек, поэтому эта оценка предполагает как если бы PCM без сжатия и
	// может не совпасть с реальным закодированным битрейтом.
	var bitrate int

	if sampleRate != 0 {
		bitrate = int((uint64(sampleRate) * uint64(bitsPerSample) * uint64(channels)) / 1000)
	}

	// Длительность: через float64, чтобы избежать переполнения при конвертации в time.Duration.
	var duration time.Duration

	if sampleRate != 0 {
		seconds := float64(totalSamples) / float64(sampleRate)
		duration = time.Duration(seconds * float64(time.Second))
	}

	return &TrackFileMetadata{
		Path:          filePath,
		Filename:      filepath.Base(filePath),
		Format:        domain.FormatFLAC,
		Codec:         domain.CodecFLAC,
		SampleRate:    int(sampleRate),
		Channels:      int(channels),
		BitsPerSample: int(bitsPerSample),
		Bitrate:       bitrate,
		TotalSamples:  totalSamples,
		Duration:      duration,
		Size:          fi.Size(),
		MD5Signature:  md5,
	}, nil
}

// WalkAndParseFlac рекурсивно обходит каталог root, проверяя каждый обычный файл и пытаясь распарсить его как FLAC.
// Файлы с иным расширением (case-intensive) пропускаются. Ошибки парсинга отдельных файлов игнорируются, чтобы один
// битый файл не обрывал обход целиком. Контекст позволяет досрочно отменить обход.
func WalkAndParseFlac(ctx context.Context, root string) ([]*TrackFileMetadata, error) {
	var results []*TrackFileMetadata

	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			// Если ОС вернула ошибку при чтении пути — прекращаем обход.
			return err
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		if d.IsDir() {
			return nil
		}

		if ext := strings.ToLower(filepath.Ext(d.Name())); ext != ".flac" {
			return nil
		}

		meta, err := ParseFlacFile(path)
		if err != nil {
			// Пропускаем файл с ошибкой парсинга.
			return nil
		}

		results = append(results, meta)

		return nil
	})

	return results, err
}
