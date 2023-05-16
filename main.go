package main

import (
	"encoding/base64"
	"errors"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/sigidagi/speech-to-text/internal/process"

	// Packages
	whisper "github.com/ggerganov/whisper.cpp/bindings/go/pkg/whisper"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type WhispParams struct {
	ModelName string
	Args      []string
}

var whParams WhispParams

func transcript(c *gin.Context) {
	id := c.Param("transcript_id")
	log.Printf("transcript id: %s", id)

	// add all other flags in front of the list
	filename := "/tmp/" + id
	//args := []string{filename}

	// TODO skip flags check
	flags, err := process.NewFlags(whParams.ModelName, whParams.Args)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	if flags.GetModel() == "" {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": errors.New("Oh Noooo, No model has been spaciefied!"),
		})
		return
	}

	// Load model
	model, _ := whisper.New(flags.GetModel())
	defer model.Close()

	str, err := process.Process(model, filename, flags)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"text":   str,
		"status": "completed",
	})
}

func main() {
	//
	whParams.ModelName = filepath.Base(os.Args[0])
	whParams.Args = os.Args[1:]

	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"Hello": "OK",
		})
	})

	r.POST("/upload", func(c *gin.Context) {
		file, _ := c.FormFile("file")
		log.Printf("File with name '%s' uploaded", file.Filename)

		b, _ := uuid.New().MarshalBinary()
		unique_filename := base64.RawURLEncoding.EncodeToString(b)
		destination_file_path := "/tmp/" + unique_filename
		c.SaveUploadedFile(file, destination_file_path)

		c.JSON(http.StatusOK, gin.H{
			"transcript_id": unique_filename,
		})
	})

	r.GET("transcript/:transcript_id", transcript)
	r.Run("0.0.0.0:8005") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
