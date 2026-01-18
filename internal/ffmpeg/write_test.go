package ffmpeg

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWriteStage(t *testing.T) {
	in := "video.mp4"
	out := "out.mp4"

	t.Run("Codecs", func(t *testing.T) {
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
			{
				name:     "CodecFor",
				builder:  New().Input(in).Output(out).CodecFor("v", 0, "libx264"),
				expected: "ffmpeg -i video.mp4 -c:v:0 libx264 out.mp4",
			},
		})
	})

	t.Run("Helpers de copy", func(t *testing.T) {
		run(t, []testCase{
			{
				name:     "CopyVideo",
				builder:  New().Input(in).Output(out).CopyVideo(),
				expected: "ffmpeg -i video.mp4 -c:v copy out.mp4",
			},
			{
				name:     "CopyAudio",
				builder:  New().Input(in).Output(out).CopyAudio(),
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
				name: "Duração de escrita (depois do -i)",
				builder: New().
					Input("movie.mkv").
					Output("out.mkv").
					T(22 * time.Second),
				expected: "ffmpeg -i movie.mkv -t 00:00:22.000 out.mkv",
			},
			{
				name: "seek, procura tempo de vídeo com -ss (depois do -i)",
				builder: New().
					Input("movie.mkv").
					Output("out.mkv").
					Ss(22 * time.Second),
				expected: "ffmpeg -i movie.mkv -ss 00:00:22.000 out.mkv",
			},
			{
				name: "Define tempo final absoluto do vídeo (depois do -i)",
				builder: New().
					Input("movie.mkv").
					Output("out.mkv").
					To(52 * time.Second),
				expected: "ffmpeg -i movie.mkv -to 00:00:52.000 out.mkv",
			},
		})
	})

	t.Run("filter", func(t *testing.T) {
		t.Run("should returns a string filter", func(t *testing.T) {
			builder := New().
				Input(in).
				Filter().Complex().
				Chaing([]string{"0:v"}, AtomicFilter{Name: "scale", Params: []string{"1280", "-1"}}, "main").
				Chaing([]string{"1:v"}, AtomicFilter{Name: "scale", Params: []string{"400", "-1"}}, "logo").
				Chaing([]string{"main", "logo"}, AtomicFilter{Name: "overlay", Params: []string{"W-w-10", "10"}}, "out").Done().
				Output(out).
				VideoCodec("libx264").
				CRF(22)

			// expected: " ",
			expected := "ffmpeg -i video.mp4 " + "-filter_complex [0:v]scale=1280:-1[main];[1:v]scale=400:-1[logo];[main][logo]overlay=W-w-10:10[out] " + "-c:v libx264 -crf 22 out.mp4"
			assert.Equal(t, expected, builder.Build())
		})
	})

	t.Run("Build não altera estado", func(t *testing.T) {
		b := New().Input("x.mp4").Output("out.mp4")

		cmd1 := b.Build()
		cmd2 := b.Build()

		require.Equal(t, cmd1, cmd2)
	})
}
