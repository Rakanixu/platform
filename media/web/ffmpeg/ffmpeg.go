//go:generate go-bindata -pkg ffmpeg -o ffmpeg_linux.go -prefix ../bin/linux/amd64/ffmpeg-3.1.3-32bit-static/ ../bin/linux/amd64/ffmpeg-3.1.3-32bit-static/ffmpeg
//go:generate go-bindata -pkg ffmpeg -o ffmpeg_darwin.go -prefix ../bin/darwin/amd64/ ../bin/darwin/amd64/ffmpeg
//go:generate go-bindata -pkg ffmpeg -o ffmpeg_windows.go -prefix ../bin/windows/amd64/ffmpeg-3.0.1-win64-static/bin/ ../bin/windows/amd64/ffmpeg-3.0.1-win64-static/bin/ffmpeg.exe
package ffmpeg
