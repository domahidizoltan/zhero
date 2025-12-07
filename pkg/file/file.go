// Package file contains some high level file managemen functions
package file

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/rs/zerolog/log"
)

func DownloadToPath(path, url string, overwrite bool) error {
	if _, err := os.Stat(path); err == nil && !overwrite {
		log.Debug().Str("file", path).Msg("file already downloaded")
		return nil
	}

	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to download file from %s: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download file from %s: unexpected status code %d", url, resp.StatusCode)
	}

	out, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create file at %s: %w", path, err)
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to write file to %s: %w", path, err)
	}

	return nil
}
