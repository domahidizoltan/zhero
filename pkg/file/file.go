// Package file contains some high level file management functions
package file

import (
	"errors"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/rs/zerolog/log"
	"golang.org/x/image/draw"
)

var (
	ErrFileDownload  = errors.New("file download failed")
	ErrFileOperation = errors.New("file operation failed")
)

var imageExtensions = []string{".jpg", ".jpeg", ".png", ".gif", ".webp", ".svg", ".bmp"}

func IsImage(ext string) bool {
	ext = strings.ToLower(strings.TrimPrefix(ext, "."))
	return slices.Contains(imageExtensions, "."+ext)
}

func IsSupportedType(ext string, supportedFileTypes []string) bool {
	ext = strings.ToLower(strings.TrimPrefix(ext, "."))
	for _, t := range supportedFileTypes {
		if strings.ToLower(strings.TrimPrefix(t, ".")) == ext {
			return true
		}
	}
	return false
}

func UploadFile(dir, filename string, data []byte, overwrite bool) (string, error) {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", fmt.Errorf("%w when creating directory %s: %w", ErrFileOperation, dir, err)
	}

	path := filepath.Join(dir, filename)
	if _, err := os.Stat(path); err == nil && !overwrite {
		log.Debug().Str("file", path).Msg("file already exists")
		return path, nil
	}

	if err := os.WriteFile(path, data, 0o644); err != nil {
		return "", fmt.Errorf("%w when writing %s: %w", ErrFileOperation, path, err)
	}

	return path, nil
}

func DeleteFilesWithPrefix(path string) error {
	files, err := filepath.Glob(path + "_*")
	if err != nil {
		return fmt.Errorf("%w when searching files with prefix %s: %w", ErrFileOperation, path, err)
	}

	for _, f := range files {
		if err := os.Remove(f); err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("%w when deleting %s: %w", ErrFileOperation, f, err)
		}
	}

	return nil
}

func GenerateThumbnail(srcPath, dstPath string, width, height int) error {
	srcFile, err := os.Open(srcPath)
	if err != nil {
		return fmt.Errorf("%w when opening source %s: %w", ErrFileOperation, srcPath, err)
	}
	defer srcFile.Close()

	srcImg, format, err := image.Decode(srcFile)
	if err != nil {
		return fmt.Errorf("failed to decode image %s: %w", srcPath, err)
	}

	// Calculate aspect-ratio-preserving dimensions
	bounds := srcImg.Bounds()
	ratio := float64(bounds.Dx()) / float64(bounds.Dy())
	thumbRatio := float64(width) / float64(height)

	var dstW, dstH int
	if ratio > thumbRatio {
		dstH = height
		dstW = int(float64(height) * ratio)
	} else {
		dstW = width
		dstH = int(float64(width) / ratio)
	}

	dstImg := image.NewRGBA(image.Rect(0, 0, dstW, dstH))
	draw.CatmullRom.Scale(dstImg, dstImg.Bounds(), srcImg, bounds, draw.Over, nil)

	dstFile, err := os.Create(dstPath)
	if err != nil {
		return fmt.Errorf("%w when creating thumbnail %s: %w", ErrFileOperation, dstPath, err)
	}
	defer dstFile.Close()

	switch format {
	case "jpeg", "jpg":
		err = jpeg.Encode(dstFile, dstImg, &jpeg.Options{Quality: 85})
	case "gif":
		err = gif.Encode(dstFile, dstImg, nil)
	default:
		err = png.Encode(dstFile, dstImg)
	}
	if err != nil {
		return fmt.Errorf("failed to encode thumbnail %s: %w", dstPath, err)
	}

	return nil
}

func DownloadToPath(path, url string, overwrite bool) error {
	if _, err := os.Stat(path); err == nil && !overwrite {
		log.Debug().Str("file", path).Msg("file already downloaded")
		return nil
	}

	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("%w from %s: %w", ErrFileDownload, url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("%w from %s: unexpected status code %d", ErrFileDownload, url, resp.StatusCode)
	}

	out, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("%w when creating %s: %w", ErrFileOperation, path, err)
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("%w when writing %s: %w", ErrFileOperation, path, err)
	}

	return nil
}
