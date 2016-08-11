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
	w.Header()["X-Content-Duration"] = []string{"5"}
	w.Header()["Content-Duration"] = []string{"750"}
	w.Header()["Accept-Ranges"] = []string{"bytes"}
	w.Header().Set("Transfer-Encoding", "identity")
	log.Printf("Streaming %v", filePath)
	cmd := exec.Command(
		"/Users/radekdymacz/Downloads/ffmpeg",
		"-ss", "0", //starting time offset
		"-t", "5",
		"-c:v", "libvpx", //video using vpx (webm) codec
		"-b:v", "1M", //1Mb/s bitrate for the video
		"-crf", "10",
		"-cpu-used", "2", //total # of cpus used
		"-threads", "4", //number of threads shared between all cpu-used
		"-deadline", "realtime", //speeds up transcode time (necessary unless you want frames dropped)
		"-strict", "-2", //ffmpeg complains about using vorbis, and wanted this param
		"-c:a", "libvorbis", //audio using the vorbis (ogg) codec
		"-s", "1280x800", //size
		"-f", "webm", //filetype for the pipe
		"-", //send output to stdout
	)

	log.Printf("Streaming %v", filePath)

	ServeCommand(cmd, w)
}
