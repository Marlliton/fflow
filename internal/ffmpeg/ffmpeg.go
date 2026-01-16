// Package ffmpeg fornece um builder fluente para compor comandos FFmpeg.
package ffmpeg

import (
	"fmt"
	"strings"
	"time"
)

type StreamType string

const (
	Video    StreamType = "v"
	Audio    StreamType = "a"
	Subtitle StreamType = "s"
)

type ffmpegBuilder struct {
	args []string
}

func New() *ffmpegBuilder {
	return &ffmpegBuilder{
		args: []string{"ffmpeg"},
	}
}

// Input adiciona um arquivo de entrada (-i).
func (f *ffmpegBuilder) Input(path string) *ffmpegBuilder {
	return f.addArgs("-i", path)
}

// Output adiciona um arquivo de saída.
func (f *ffmpegBuilder) Output(path string) *ffmpegBuilder {
	return f.addArgs(path)
}

// Override adiciona a flag de sobrescrita automática (-y).
func (f *ffmpegBuilder) Override() *ffmpegBuilder {
	return f.addArgs("-y")
}

// VideoCodec define o codec de vídeo (-c:v).
func (f *ffmpegBuilder) VideoCodec(codec string) *ffmpegBuilder {
	return f.addArgs("-c:v", codec)
}

// AudioCodec define o codec de áudio (-c:a).
func (f *ffmpegBuilder) AudioCodec(codec string) *ffmpegBuilder {
	return f.addArgs("-c:a", codec)
}

// SubtitleCodec define o codec de legenda (-c:s).
func (f *ffmpegBuilder) SubtitleCodec(codec string) *ffmpegBuilder {
	return f.addArgs("-c:s", codec)
}

// CodecFor define o codec de um stream específico (-c:<stream>:<index>).
func (f *ffmpegBuilder) CodecFor(stream StreamType, index int, codec string) *ffmpegBuilder {
	return f.addArgs(fmt.Sprintf("-c:%s:%d", stream, index), codec)
}

// CopyVideo apenas copia os pacotes do input para o output (sem encoder)
func (f *ffmpegBuilder) CopyVideo() *ffmpegBuilder {
	return f.addArgs("-c:v", "copy")
}

// CopyAudio apenas copia os pacotes do input para o output (sem encoder)
func (f *ffmpegBuilder) CopyAudio() *ffmpegBuilder {
	return f.addArgs("-c:a", "copy")
}

// LimitDuration adiciona a flag -t.
// O efeito depende da posição no comando:
// - antes do Input: limita a leitura do input
// - antes do output: limita a escrita do output
// Ex: LimitDuration(30 * time.Second) = 30s
func (f *ffmpegBuilder) LimitDuration(d time.Duration) *ffmpegBuilder {
	return f.addArgs("-t", fmtDuration(d))
}

// addArgs adiciona argumentos ao comando FFmpeg.
func (f *ffmpegBuilder) addArgs(args ...string) *ffmpegBuilder {
	f.args = append(f.args, args...)
	return f
}

// Build retorna o comando FFmpeg final como string.
func (f *ffmpegBuilder) Build() string {
	var sb strings.Builder

	for _, arg := range f.args {
		fmt.Fprintf(&sb, "%s ", arg)
	}

	return strings.TrimSpace(sb.String())
}
