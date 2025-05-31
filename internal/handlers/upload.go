package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/redis/go-redis/v9"
	"golang.org/x/net/context"
)

var rdb *redis.Client

func init() {
	rdb = redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "",
		DB:       0,
	})
}

type Job struct {
	FilePath string `json:"file_path"`
	Uploaded int64  `json:"uploaded"`
}

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10 << 20) // 10MB
	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Cannot read file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	timestamp := time.Now().UnixNano()
	filename := fmt.Sprintf("%d_%s", timestamp, handler.Filename)
	path := filepath.Join("/tmp", filename)

	f, err := os.Create(path)
	if err != nil {
		http.Error(w, "Cannot save file", http.StatusInternalServerError)
		return
	}
	defer f.Close()

	_, err = io.Copy(f, file)
	if err != nil {
		http.Error(w, "Failed to save file", http.StatusInternalServerError)
		return
	}

	job := Job{
		FilePath: path,
		Uploaded: timestamp,
	}
	jobData, err := json.Marshal(job)
	if err != nil {
		http.Error(w, "Failed to serialize job", http.StatusInternalServerError)
		return
	}

	err = rdb.LPush(context.Background(), "import_jobs", jobData).Err()
	if err != nil {
		http.Error(w, "Failed to queue job", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	fmt.Fprintf(w, "File %s uploaded successfully", handler.Filename)
}
