package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
)

func ServeCommand(cmd *exec.Cmd, w io.Writer) error {
	stdout, err := cmd.StdoutPipe()
	defer stdout.Close()
	if err != nil {
		log.Printf("Error opening stdout of command: %v", err)
		return err
	}
	err = cmd.Start()
	if err != nil {
		log.Printf("Error starting command: %v", err)
		return err
	}
	_, err = io.Copy(w, stdout)

	if err != nil {
		log.Printf("Error copying data to client: %v", err)
		// Ask the process to exit
		cmd.Process.Signal(syscall.SIGKILL)
		cmd.Process.Wait()
		return err
	}
	cmd.Wait()
	return nil
}

var videoSuffixes = []string{".mp4", ".avi", ".mkv", ".flv", ".wmv", ".mov", ".mpg"}

type VideoInfo struct {
	Duration float64 `json:"duration"`
}

func FilenameLooksLikeVideo(name string) bool {
	for _, suffix := range videoSuffixes {
		if strings.HasSuffix(name, suffix) {
			return true
		}
	}
	return false
}

func GetRawFFMPEGInfo(path string) ([]byte, error) {
	log.Printf("Executing ffprobe for %v", path)
	cmd := exec.Command("/Users/radekdymacz/Downloads/ffprobe", "-v", "quiet", "-print_format", "json", "-show_format", ""+path+"")
	data, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("Error executing ffprobe for file '%v':", path, err)
	}
	return data, nil
}

func GetFFMPEGJson(path string) (map[string]interface{}, error) {
	data, cmderr := GetRawFFMPEGInfo(path)
	if cmderr != nil {
		return nil, cmderr
	}
	var info map[string]interface{}
	err := json.Unmarshal(data, &info)
	if err != nil {
		return nil, fmt.Errorf("Error unmarshalling JSON from ffprobe output for file '%v':", path, err)
	}
	return info, nil
}

func GetVideoInformation(path string) (*VideoInfo, error) {
	info, jsonerr := GetFFMPEGJson(path)
	if jsonerr != nil {
		return nil, jsonerr
	}
	log.Printf("ffprobe for %v returned", path, info)
	if _, ok := info["format"]; !ok {
		return nil, fmt.Errorf("ffprobe data for '%v' does not contain format info", path)
	}
	format := info["format"].(map[string]interface{})
	if _, ok := format["duration"]; !ok {
		return nil, fmt.Errorf("ffprobe format data for '%v' does not contain duration", path)
	}
	duration, perr := strconv.ParseFloat(format["duration"].(string), 64)
	if perr != nil {
		return nil, fmt.Errorf("Could not parse duration (%v) of '%v': ", format["duration"].(string), path, perr)
	}
	return &VideoInfo{duration}, nil
}
