package file

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/domahidizoltan/zhero/config"
	"github.com/domahidizoltan/zhero/controller"
	file_pkg "github.com/domahidizoltan/zhero/pkg/file"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type Controller struct {
	uploadsCfg config.UploadsConfig
}

func NewController(uploadsCfg config.UploadsConfig) Controller {
	return Controller{
		uploadsCfg: uploadsCfg,
	}
}

func (pc *Controller) Upload(c *gin.Context) {
	fileDir := c.GetHeader("X-File-Dir")
	fileNamePrefix := c.GetHeader("X-File-Name-Prefix")
	formFileName := c.GetHeader("X-Form-File-Name")

	file, header, err := c.Request.FormFile(formFileName)
	if err != nil {
		c.Data(http.StatusBadRequest, gin.MIMEHTML, []byte(`No file selected`))
		return
	}
	defer file.Close()

	ext := strings.ToLower(filepath.Ext(header.Filename))
	if ext == "" {
		c.Data(http.StatusBadRequest, gin.MIMEHTML, []byte(`File has no extension`))
		return
	}

	if !file_pkg.IsSupportedType(ext, append(pc.uploadsCfg.SupportedFileTypes, pc.uploadsCfg.SupportedImageTypes...)) {
		c.Data(http.StatusBadRequest, gin.MIMEHTML, []byte("Unsupported file type: "+ext))
		return
	}

	data := make([]byte, header.Size)
	if _, err := file.Read(data); err != nil {
		controller.InternalServerError(c, "failed to read file", err)
		return
	}

	timestamp := time.Now().Unix()
	fileName := fmt.Sprintf("%s_%d.%s", fileNamePrefix, timestamp, ext[1:])
	dir := filepath.Join(".", UploadsPath, fileDir)

	if err := pc.delete(fileDir, fileNamePrefix); err != nil {
		controller.InternalServerError(c, "failed to delete file", err)
		return
	}

	filePath, err := file_pkg.UploadFile(dir, fileName, data, false)
	if err != nil {
		controller.InternalServerError(c, "failed to save file", err)
		return
	}

	var thumbPath string
	if file_pkg.IsImage(ext) {
		if size, err := pc.uploadsCfg.GetThumbnailSize(); err != nil {
		} else {
			thumbFileName := fmt.Sprintf("%s_%d_thumb.%s", fileNamePrefix, timestamp, ext[1:])
			thumbPath = filepath.Join(dir, thumbFileName)
			if therr := file_pkg.GenerateThumbnail(filePath, thumbPath, size.Width, size.Height); therr != nil {
				log.Warn().Err(therr).Str("file", filePath).Msg("thumbnail generation failed")
			}
		}
	}

	relPath := filepath.ToSlash(filePath)
	preview := ""
	if thumbPath != "" {
		preview = fmt.Sprintf(`<img class="w-20 h-20 object-cover rounded" src="%s" width="100" height="100"/><br/>`, "/"+thumbPath)
	}
	preview += filepath.Base(filePath)
	html := fmt.Sprintf(`<a href="/%[1]s" target="_blank" class="text-info link">%[2]s</a>
	<button type="button" class="btn btn-xs btn-error ml-2"
		hx-trigger="click"
		hx-delete="/admin/file"
	  hx-headers='{"X-File-Dir":"%[3]s", "X-File-Name-Prefix":"%[4]s"}'
		hx-confirm="Delete '%[4]s' file '%[5]s'?"
		hx-target="#file-result-%[4]s"
		hx-swap="innerHTML">
		<i class="fa-solid fa-trash"></i>
	</button>`, relPath, preview, fileDir, fileNamePrefix, fileName)

	trigger, _ := json.Marshal(map[string]map[string]string{"fileUploaded": {"field": fileNamePrefix, "fileName": fileName}})
	c.Header("HX-Trigger", string(trigger))
	c.Data(http.StatusOK, gin.MIMEHTML, []byte(html))
}

func (pc *Controller) Delete(c *gin.Context) {
	fileDir := c.GetHeader("X-File-Dir")
	fileNamePrefix := c.GetHeader("X-File-Name-Prefix")

	if err := pc.delete(fileDir, fileNamePrefix); err != nil {
		controller.InternalServerError(c, "failed to delete file", err)
		return
	}

	c.Data(http.StatusOK, gin.MIMEHTML, []byte(""))
}

func (pc *Controller) delete(fileDir, fileNamePrefix string) error {
	dir := filepath.Join(".", UploadsPath, fileDir)
	filePath := fmt.Sprintf("%s/%s", dir, fileNamePrefix)
	return file_pkg.DeleteFilesWithPrefix(filePath)
}
