package handler

import (
	"log"
	"net/http"
	"os/exec"
	"path"
)

type MP4Handler struct {
	root string
}

func NewMP4Handler(root string) *MP4Handler {
	return &MP4Handler{root}
}

func (s *MP4Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := path.Join(s.root, r.URL.Path)
	w.Header()["Content-Type"] = []string{"video/mp4"}
	w.Header()["Transfer-Encoding"] = []string{"chunked"}
	w.Header()["Accept-Ranges"] = []string{"bytes"}
	log.Printf("Headers : %s", w.Header())
	cmd := exec.Command("/Users/radekdymacz/Downloads/ffmpeg", "-i", path, "-vcodec", "libx264", "-strict", "experimental", "-acodec", "aac", "-s", "480x270 ", "-pix_fmt", "yuv420p", "-r", "25", "-profile:v", "baseline", "-b:v", "2000k", "-maxrate", "2500k", "-movflags", "+faststart", "-f", "mp4", "-")
	if err := ServeCommand(cmd, w); err != nil {
		log.Printf("Error serving MP4: %v", err)
	}
}
