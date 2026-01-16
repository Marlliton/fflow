package ffmpeg

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type tests struct {
	name     string
	builder  *ffmpegBuilder
	expected string
}

func exec(t *testing.T, tests []tests) {
	t.Helper()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.expected, tt.builder.Build())
		})
	}
}

func Test_FFmpegBuilder(t *testing.T) {
	t.Run("Estado inicial", func(t *testing.T) {
		f := New()
		require.Equal(t, "ffmpeg", f.Build(), "Deve iniciar apenas com o binário")
	})

	t.Run("Argumentos básicos", func(t *testing.T) {
		tests := []tests{
			{
				name:     "Input",
				builder:  New().Input("video.mp4"),
				expected: "ffmpeg -i video.mp4",
			},
			{
				name:     "Override",
				builder:  New().Override(),
				expected: "ffmpeg -y",
			},
		}

		exec(t, tests)
	})

	t.Run("Codecs simples", func(t *testing.T) {
		tests := []tests{
			{
				name:     "VideoCodec",
				builder:  New().VideoCodec("libx264"),
				expected: "ffmpeg -c:v libx264",
			},
			{
				name:     "AudioCodec",
				builder:  New().AudioCodec("aac"),
				expected: "ffmpeg -c:a aac",
			},
			{
				name:     "SubtitleCodec",
				builder:  New().SubtitleCodec("srt"),
				expected: "ffmpeg -c:s srt",
			},
		}

		exec(t, tests)
	})

	t.Run("Copy helpers", func(t *testing.T) {
		tests := []tests{
			{
				name:     "CopyVideo",
				builder:  New().CopyVideo(),
				expected: "ffmpeg -c:v copy",
			},
			{
				name:     "CopyAudio",
				builder:  New().CopyAudio(),
				expected: "ffmpeg -c:a copy",
			},
		}

		exec(t, tests)
	})

	t.Run("Codecs por stream", func(t *testing.T) {
		tests := []tests{
			{
				name:     "Codec por stream de legenda",
				builder:  New().CodecFor(Subtitle, 1, "srt"),
				expected: "ffmpeg -c:s:1 srt",
			},
			{
				name:     "Codec por stream de vídeo",
				builder:  New().CodecFor(Video, 0, "libx264"),
				expected: "ffmpeg -c:v:0 libx264",
			},
		}

		exec(t, tests)
	})

	t.Run("Definições de tempo", func(t *testing.T) {
		tests := []tests{
			{
				name:     "Adiciona duração de entrada",
				builder:  New().LimitDuration(30 * time.Second).Input("movie.mkv"),
				expected: "ffmpeg -t 00:00:30.000 -i movie.mkv",
			},
			{
				name: "Adiciona duração de saída",
				builder: New().Input("movie.mkv").AudioCodec("aac").LimitDuration(22 * time.Second).
					Output("out.mkv"),
				expected: "ffmpeg -i movie.mkv -c:a aac -t 00:00:22.000 out.mkv",
			},
		}

		exec(t, tests)
	})

	t.Run("Encadeamento e ordem", func(t *testing.T) {
		cmd := New().
			Override().
			Input("movie.mkv").
			Input("audio.mp3").
			VideoCodec("libx264").
			AudioCodec("aac").
			SubtitleCodec("srt").
			CodecFor(Subtitle, 1, "srt").
			Build()

		require.Equal(
			t,
			"ffmpeg -y -i movie.mkv -i audio.mp3 -c:v libx264 -c:a aac -c:s srt -c:s:1 srt",
			cmd,
		)
	})

	t.Run("Fluent interface", func(t *testing.T) {
		f := New()
		f2 := f.Input("test.mp4")

		assert.Same(t, f, f2, "O builder deve retornar a mesma instância")
	})

	t.Run("Build não altera estado", func(t *testing.T) {
		b := New().Input("x.mp4")

		cmd1 := b.Build()
		cmd2 := b.Build()

		require.Equal(t, cmd1, cmd2)
	})
}
