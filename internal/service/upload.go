package service

import (
	"encoding/base64"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/friedrichad/golang_web_api_demo/internal/common"
	"github.com/friedrichad/golang_web_api_demo/internal/configs/upload"
	"github.com/friedrichad/golang_web_api_demo/internal/model"
	"github.com/friedrichad/golang_web_api_demo/internal/model/constants"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type IUploadService interface {
	UploadBase64(c *gin.Context) (*model.UploadResponse, *common.Error)
	UploadMultipart(c *gin.Context) (*model.UploadResponse, *common.Error)
	UploadMultiple(c *gin.Context) ([]model.UploadResponse, int, *common.Error)
}

type UploadService struct {
	BasePath  string
	ReturnUrl string
}

func NewUploadService() IUploadService {
	cfg := upload.LoadUploadConfig()
	return &UploadService{
		BasePath:  cfg.BasePath,
		ReturnUrl: cfg.ReturnUrl,
	}
}

func (s *UploadService) UploadMultipart(c *gin.Context) (*model.UploadResponse, *common.Error) {

	fileHeader, err := c.FormFile("file")
	if err != nil {
		log.Println(err)
		return nil, common.FileError
	}

	if fileHeader == nil {
		return nil, common.FileEmpty
	}

	file, err := fileHeader.Open()
	if err != nil {
		log.Println(err)
		return nil, common.SystemError
	}

	defer file.Close()

	buffer, err := io.ReadAll(file)
	if err != nil {
		log.Println(err)
		return nil, common.SystemError
	}

	contentType := http.DetectContentType(buffer)

	if err := constants.ValidateFileType(contentType); err != nil {
		return nil, common.FileError
	}

	fileName := uuid.New().String() + filepath.Ext(fileHeader.Filename)

	path := filepath.Join(s.BasePath, fileName)

	err = os.WriteFile(path, buffer, 0644)
	if err != nil {
		log.Println(err)
		return nil, common.SystemError
	}

	url := s.ReturnUrl + "/" + fileName

	return &model.UploadResponse{
		FileName: fileName,
		FileUrl:  url,
		Size:     fileHeader.Size,
	}, nil
}

func (s *UploadService) UploadMultiple(c *gin.Context) ([]model.UploadResponse, int, *common.Error) {
	form, err := c.MultipartForm()
	count := 0
	if err != nil {
		log.Println(err)
		return nil, 0, common.RequestInvalid
	}

	files := form.File["files"] // "files" là tên field trong form
	if len(files) == 0 {
		return nil, 0, common.FileEmpty
	}

	var responses []model.UploadResponse

	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			log.Println(err)
			return nil, 0, common.SystemError
		}
		defer file.Close()

		buffer, err := io.ReadAll(file)
		if err != nil {
			log.Println(err)
			return nil, 0, common.SystemError
		}

		contentType := http.DetectContentType(buffer)
		if err := constants.ValidateFileType(contentType); err != nil {
			return nil, 0, common.FileError
		}

		fileName := uuid.New().String() + filepath.Ext(fileHeader.Filename)
		path := filepath.Join(s.BasePath, fileName)

		if err := os.WriteFile(path, buffer, 0644); err != nil {
			log.Println(err)
			return nil, 0, common.SystemError
		}

		url := s.ReturnUrl + "/" + fileName

		responses = append(responses, model.UploadResponse{
			FileName: fileName,
			FileUrl:  url,
			Size:     fileHeader.Size,
		})
		count++
	}

	return responses, count, nil
}

func (s *UploadService) UploadBase64(c *gin.Context) (*model.UploadResponse, *common.Error) {

	var req model.UploadBase64Request

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println(err)
		return nil, common.RequestInvalid
	}
	log.Printf("Received base64 data of length: %d", len(req.Data))
	data := removeBase64Prefix(req.Data)

	decoded, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		log.Println(err)
		return nil, common.FileError
	}

	if err := constants.ValidateFileSize(int64(len(decoded))); err != nil {
		log.Println(err)
		return nil, common.FileError
	}

	contentType := http.DetectContentType(decoded)

	if err := constants.ValidateFileType(contentType); err != nil {
		log.Println(err)
		return nil, common.FileError
	}

	fileName := uuid.New().String() + filepath.Ext(req.FileName)

	path := filepath.Join(s.BasePath, fileName)

	err = os.WriteFile(path, decoded, 0644)
	if err != nil {
		log.Println(err)
		return nil, common.SystemError
	}

	url := s.ReturnUrl + "/" + fileName

	return &model.UploadResponse{
		FileName: fileName,
		FileUrl:  url,
		Size:     int64(len(decoded)),
	}, nil
}

func removeBase64Prefix(data string) string {
	idx := strings.Index(data, ",")
	if idx != -1 {
		return data[idx+1:]
	}
	return data
}
