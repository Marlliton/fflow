package ffmpeg

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type WriteStage interface {
	// Ss adiciona a flag -ss após os inputs (-i), realizando um seek preciso no output.
	//
	// Ss adds the -ss flag after the inputs (-i), performing a precise seek on the output.
	Ss(time.Duration) WriteStage

	// To adiciona a flag -to após os inputs (-i), definindo o tempo final absoluto do output.
	//
	// To adds the -to flag after the inputs (-i), defining the absolute end time of the output.
	To(time.Duration) WriteStage

	// T adiciona a flag -t após os inputs (-i), limitando a duração do output.
	//
	// T adds the -t flag after the inputs (-i), limiting the output duration.
	T(time.Duration) WriteStage

	// VideoCodec define o codec de vídeo do output (-c:v).
	//
	// VideoCodec sets the output video codec (-c:v).
	VideoCodec(codec string) WriteStage

	// AudioCodec define o codec de áudio do output (-c:a).
	//
	// AudioCodec sets the output audio codec (-c:a).
	AudioCodec(codec string) WriteStage

	// SubtitleCodec define o codec de legenda do output (-c:s).
	//
	// SubtitleCodec sets the output subtitle codec (-c:s).
	SubtitleCodec(codec string) WriteStage

	// CopyVideo copia o stream de vídeo sem recodificar (-c:v copy).
	//
	// CopyVideo copies the video stream without re-encoding (-c:v copy).
	CopyVideo() WriteStage

	// CopyAudio copia o stream de áudio sem recodificar (-c:a copy).
	//
	// CopyAudio copies the audio stream without re-encoding (-c:a copy).
	CopyAudio() WriteStage

	// CodecFor define o codec de um stream específico do output (-c:<stream>:<index>).
	//
	// CodecFor sets the codec for a specific output stream (-c:<stream>:<index>).
	CodecFor(stream StreamType, index int, codec string) WriteStage

	// CRF define o fator de qualidade constante para encoders de vídeo.
	//
	// CRF sets the constant quality factor for video encoders.
	CRF(value int) WriteStage

	// Output define o arquivo de saída.
	//
	// Output sets the output file path.
	Output(path string) WriteStage

	// Build monta o comando FFmpeg final respeitando a ordem semântica.
	//
	// Build assembles the final FFmpeg command respecting semantic order.
	Build() string
}

type writeCtx struct{ b *ffmpegBuilder }

func (c *writeCtx) T(d time.Duration) WriteStage {
	c.b.write = append(c.b.write, "-t", fmtDuration(d))
	return c
}

func (c *writeCtx) Ss(d time.Duration) WriteStage {
	c.b.write = append(c.b.write, "-ss", fmtDuration(d))
	return c
}

func (c *writeCtx) To(d time.Duration) WriteStage {
	c.b.write = append(c.b.write, "-to", fmtDuration(d))
	return c
}

func (c *writeCtx) VideoCodec(codec string) WriteStage {
	c.b.write = append(c.b.write, "-c:v", codec)
	return c
}

func (c *writeCtx) AudioCodec(codec string) WriteStage {
	c.b.write = append(c.b.write, "-c:a", codec)
	return c
}

func (c *writeCtx) SubtitleCodec(codec string) WriteStage {
	c.b.write = append(c.b.write, "-c:s", codec)
	return c
}

func (c *writeCtx) CodecFor(stream StreamType, index int, codec string) WriteStage {
	c.b.write = append(c.b.write, fmt.Sprintf("-c:%s:%d", stream, index), codec)
	return c
}

func (c *writeCtx) CopyVideo() WriteStage {
	c.b.write = append(c.b.write, "-c:v", "copy")
	return c
}

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

func (c *writeCtx) Build() string {
	var args []string

	args = append(args, c.b.global...)
	args = append(args, c.b.read...)
	if len(c.b.filters) > 0 {
		pipeline := Pipeline{Nodes: c.b.filters}
		if pipeline.NeedsComplex() {
			args = append(args, "-filter_complex", pipeline.String())
		} else {
			args = append(args, "-vf", pipeline.String())
		}
	}
	args = append(args, c.b.write...)
	args = append(args, c.b.output)
	return strings.Join(args, " ")
}
