package handler

import (
	"fmt"
	"github.com/getlantern/byteexec"
	"github.com/kazoup/platform/media/web/ffmpeg"
	"log"
	"net/http"
	"path"
)

type FrameHandler struct {
	root string
	be   *byteexec.Exec
}

func NewFrameHandler(root string) *FrameHandler {

	pb, err := ffmpeg.Asset("ffmpeg")
	if err != nil {
		log.Print(err.Error())
	}
	be, err := byteexec.New(pb, "ffmpeg")

	if err != nil {
		log.Print(err.Error())
	}
	return &FrameHandler{
		root: root,
		be:   be,
	}
}

func (s *FrameHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := path.Join(s.root, r.URL.Path)
	cmd := s.be.Command("-loglevel", "error", "-ss", "00:00:01", "-i", path, "-vf", "scale=420:-1", "-frames:v", "1", "-f", "image2", "-")
	if err := ServeCommand(cmd, w); err != nil {
		fmt.Fprintf(w, "Error serving screenshot: %v", err.Error())
	}
}
