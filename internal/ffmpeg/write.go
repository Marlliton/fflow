package ffmpeg

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type WriteStage interface {
	Ss(time.Duration) WriteStage
	To(time.Duration) WriteStage
	T(time.Duration) WriteStage
	VideoCodec(codec string) WriteStage
	AudioCodec(codec string) WriteStage
	SubtitleCodec(codec string) WriteStage
	CopyVideo() WriteStage
	CopyAudio() WriteStage
	CRF(value int) WriteStage

	Output(path string) WriteStage

	Build() string
}

type writeCtx struct{ b *ffmpegBuilder }

func (c *writeCtx) T(d time.Duration) WriteStage {
	c.b.write = append(c.b.write, "-t", fmtDuration(d))
	return c
}

// Ss adiciona a flag -ss (seek).
// Antes do -i: seek rápido (menos preciso).
// Depois do -i: seek preciso (mais lento).
func (c *writeCtx) Ss(d time.Duration) WriteStage {
	c.b.write = append(c.b.write, "-ss", fmtDuration(d))
	return c
}

// To adiciona a flag -to (tempo final absoluto).
// Processa até atingir o timestamp informado.
// Geralmente usado em conjunto com -ss.
func (c *writeCtx) To(d time.Duration) WriteStage {
	c.b.write = append(c.b.write, "-to", fmtDuration(d))
	return c
}

// VideoCodec define o codec de vídeo (-c:v).
func (c *writeCtx) VideoCodec(codec string) WriteStage {
	c.b.write = append(c.b.write, "-c:v", codec)
	return c
}

// AudioCodec define o codec de áudio (-c:a).
func (c *writeCtx) AudioCodec(codec string) WriteStage {
	c.b.write = append(c.b.write, "-c:a", codec)
	return c
}

// SubtitleCodec define o codec de legenda (-c:s).
func (c *writeCtx) SubtitleCodec(codec string) WriteStage {
	c.b.write = append(c.b.write, "-c:s", codec)
	return c
}

// CodecFor define o codec de um stream específico (-c:<stream>:<index>).
func (c *writeCtx) CodecFor(stream StreamType, index int, codec string) WriteStage {
	c.b.write = append(c.b.write, fmt.Sprintf("-c:%s:%d", stream, index), codec)
	return c
}

// CopyVideo apenas copia os pacotes do input para o output (sem encoder)
func (c *writeCtx) CopyVideo() WriteStage {
	c.b.write = append(c.b.write, "-c:v", "copy")
	return c
}

// CopyAudio apenas copia os pacotes do input para o output (sem encoder)
func (c *writeCtx) CopyAudio() WriteStage {
	c.b.write = append(c.b.write, "-c:a", "copy")
	return c
}

func (c *writeCtx) CRF(value int) WriteStage {
	c.b.write = append(c.b.write, "-crf", strconv.Itoa(value))
	return c
}

func (c *writeCtx) Output(path string) WriteStage {
	c.b.output = path
	return c
}

// TODO: escrever lógica de build
func (c *writeCtx) Build() string {
	var args []string

	args = append(args, c.b.global...)
	args = append(args, c.b.read...)
	args = append(args, c.b.write...)
	args = append(args, c.b.output)
	return strings.Join(args, " ")
}
