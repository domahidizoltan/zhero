// Package file contains some high level file managemen functions
package file

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/rs/zerolog/log"
)

var (
	ErrFileDownload  = errors.New("file download failed")
	ErrFileOperation = errors.New("file operation failed")
)

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
