// Package ffmpeg fornece um builder fluente para compor comandos FFmpeg.
package ffmpeg

type StreamType string

const (
	Video    StreamType = "v"
	Audio    StreamType = "a"
	Subtitle StreamType = "s"
)

type ffmpegBuilder struct {
	global           []string
	read             []string
	write            []string
	filters          []Filter
	simpleFilterFlag string
	output           string
}

// New inicia um novo construtor de comando FFmpeg, retornando uma GlobalStage.
// Este é o ponto de entrada para construir qualquer comando FFmpeg, permitindo a configuração de opções globais
// antes de especificar as entradas.
//
// New starts a new FFmpeg command builder, returning a GlobalStage.
// This is the entry point for building any FFmpeg command, allowing the configuration of global options
// before specifying inputs.
func New() *globalCtx {
	return &globalCtx{
		b: &ffmpegBuilder{global: []string{"ffmpeg"}},
	}
}
