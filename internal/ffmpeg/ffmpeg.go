// Package ffmpeg fornece um builder fluente para compor comandos FFmpeg.
package ffmpeg

type StreamType string

const (
	Video    StreamType = "v"
	Audio    StreamType = "a"
	Subtitle StreamType = "s"
)

type ffmpegBuilder struct {
	global  []string
	read    []string
	write   []string
	filters []string
	output  string
}

func New() *globalCtx {
	return &globalCtx{
		b: &ffmpegBuilder{global: []string{"ffmpeg"}},
	}
}
