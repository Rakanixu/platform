package handler

import (
	"log"
	"net/http"
	"os/exec"
	"path"
)

type WebmHandler struct {
	root string
}

func NewWebmHandler(root string) *WebmHandler {
	return &WebmHandler{root}
}

func (s *WebmHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//filePath := path.Join(s.root, r.URL.Path[0:strings.LastIndex(r.URL.Path, "/")])

	filePath := path.Join(s.root, r.URL.Path)
	//idx, _ := strconv.ParseInt(r.URL.Path[strings.LastIndex(r.URL.Path, "/")+1:strings.LastIndex(r.URL.Path, ".")], 0, 64)
	//:startTime := 0
	w.Header()["Content-Type"] = []string{"video/webm"}

	log.Printf("Streaming %v", filePath)
	cmd := exec.Command("/Users/radekdymacz/Downloads/ffmpeg", "-i", filePath, "-ss", "0", "-t", "5", "-c:v", "libvpx", "-profile:v", "baseline", "-b:v", "2000k", "-c:a", "libvorbis", "-an", "-f", "webm", "-dash", "1", "-")

	ServeCommand(cmd, w)
}
