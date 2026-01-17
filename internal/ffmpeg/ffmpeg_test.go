package ffmpeg

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testCase struct {
	name     string
	builder  WriteStage
	expected string
}

func run(t *testing.T, tests []testCase) {
	t.Helper()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.expected, tt.builder.Build())
		})
	}
}

func Test_FFmpegBuilder(t *testing.T) {
	in := "video.mp4"
	out := "out.mp4"

	t.Run("Argumentos básicos", func(t *testing.T) {
		run(t, []testCase{
			{
				name:     "Input e Output",
				builder:  New().Input(in).Output(out),
				expected: "ffmpeg -i video.mp4 out.mp4",
			},
			{
				name:     "Override",
				builder:  New().Override().Input(in).Output(out),
				expected: "ffmpeg -y -i video.mp4 out.mp4",
			},
		})
	})

	t.Run("Codecs simples", func(t *testing.T) {
		run(t, []testCase{
			{
				name:     "VideoCodec",
				builder:  New().Input(in).Output(out).VideoCodec("libx264"),
				expected: "ffmpeg -i video.mp4 -c:v libx264 out.mp4",
			},
			{
				name:     "AudioCodec",
				builder:  New().Input(in).Output(out).AudioCodec("aac"),
				expected: "ffmpeg -i video.mp4 -c:a aac out.mp4",
			},
			{
				name:     "SubtitleCodec",
				builder:  New().Input(in).Output(out).SubtitleCodec("srt"),
				expected: "ffmpeg -i video.mp4 -c:s srt out.mp4",
			},
		})
	})

	t.Run("Helpers de copy", func(t *testing.T) {
		run(t, []testCase{
			{
				name:     "CopyVideo",
				builder:  New().Input(in).Output(out).VideoCodec("copy"),
				expected: "ffmpeg -i video.mp4 -c:v copy out.mp4",
			},
			{
				name:     "CopyAudio",
				builder:  New().Input(in).Output(out).AudioCodec("copy"),
				expected: "ffmpeg -i video.mp4 -c:a copy out.mp4",
			},
		})
	})

	t.Run("CRF", func(t *testing.T) {
		run(t, []testCase{
			{
				name:     "Define qualidade de vídeo",
				builder:  New().Input(in).Output(out).VideoCodec("libx264").CRF(22),
				expected: "ffmpeg -i video.mp4 -c:v libx264 -crf 22 out.mp4",
			},
		})
	})

	t.Run("Definições de tempo", func(t *testing.T) {
		run(t, []testCase{
			{
				name: "Duração de leitura (antes do -i)",
				builder: New().
					Input("movie.mkv").
					T(30 * time.Second).
					Output("out.mkv"),
				expected: "ffmpeg -t 00:00:30.000 -i movie.mkv out.mkv",
			},
			{
				name: "Duração de escrita (depois do -i)",
				builder: New().
					Input("movie.mkv").
					Output("out.mkv").
					T(22 * time.Second),
				expected: "ffmpeg -i movie.mkv -t 00:00:22.000 out.mkv",
			},
		})
	})

	t.Run("Múltiplos inputs", func(t *testing.T) {
		cmd := New().
			Override().
			Input("movie.mkv").
			Input("audio.mp3").
			Output("out.mkv").
			VideoCodec("libx264").
			AudioCodec("aac").
			SubtitleCodec("srt").
			Build()

		require.Equal(
			t,
			"ffmpeg -y -i movie.mkv -i audio.mp3 -c:v libx264 -c:a aac -c:s srt out.mkv",
			cmd,
		)
	})

	t.Run("Fluent interface mantém o mesmo estado", func(t *testing.T) {
		f := New()
		f2 := f.Input("test.mp4")

		assert.NotNil(t, f2)
	})

	t.Run("Build não altera estado", func(t *testing.T) {
		b := New().Input("x.mp4").Output("out.mp4")

		cmd1 := b.Build()
		cmd2 := b.Build()

		require.Equal(t, cmd1, cmd2)
	})
}
