package fflow

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
				name:     "CodecFor video stream",
				builder:  New().Input(in).Output(out).CodecFor(Video, 0, "libx264"),
				expected: "ffmpeg -i video.mp4 -c:v:0 libx264 out.mp4",
			},
			{
				name:     "CodecFor audio stream",
				builder:  New().Input(in).Output(out).CodecFor(Audio, 1, "libmp3lame"),
				expected: "ffmpeg -i video.mp4 -c:a:1 libmp3lame out.mp4",
			},
			{
				name:     "CodecFor subtitle stream",
				builder:  New().Input(in).Output(out).CodecFor(Subtitle, 0, "mov_text"),
				expected: "ffmpeg -i video.mp4 -c:s:0 mov_text out.mp4",
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
					T(22 * time.Second).
					Output("out.mkv"),
				expected: "ffmpeg -i movie.mkv -t 00:00:22.000 out.mkv",
			},
			{
				name: "seek, procura tempo de vídeo com -ss (depois do -i)",
				builder: New().
					Input("movie.mkv").
					Ss(22 * time.Second).
					Output("out.mkv"),
				expected: "ffmpeg -i movie.mkv -ss 00:00:22.000 out.mkv",
			},
			{
				name: "Define tempo final absoluto do vídeo (depois do -i)",
				builder: New().
					Input("movie.mkv").
					To(52 * time.Second).
					Output("out.mkv"),
				expected: "ffmpeg -i movie.mkv -to 00:00:52.000 out.mkv",
			},
		})
	})

	t.Run("Outros", func(t *testing.T) {
		run(t, []testCase{
			{
				name:     "Preset",
				builder:  New().Input(in).Output(out).VideoCodec("libx264").Preset("slow"),
				expected: "ffmpeg -i video.mp4 -c:v libx264 -preset slow out.mp4",
			},
			{
				name:     "Argumento bruto (depois do -i)",
				builder:  New().Input(in).Output(out).Raw("-f"),
				expected: "ffmpeg -i video.mp4 -f out.mp4",
			},
		})
	})

	t.Run("filter", func(t *testing.T) {
		t.Run("Complex filter with overlay and map", func(t *testing.T) {
			builder := New().
				Input(in).
				Input("logo.png").
				Filter().
				Complex().
				Chain([]string{"0:v"}, AtomicFilter{Name: "scale", Params: []string{"1280", "-1"}}, []string{"main"}).
				Chain([]string{"1:v"}, AtomicFilter{Name: "scale", Params: []string{"400", "-1"}}, []string{"logo"}).
				Chain([]string{"main", "logo"}, AtomicFilter{Name: "overlay", Params: []string{"W-w-10", "10"}}, []string{"out"}).
				Done().
				Map("out").
				Output(out).
				VideoCodec("libx264").
				CRF(22)

			expected := "ffmpeg -i video.mp4 -i logo.png " +
				"-filter_complex [0:v]scale=1280:-1[main];[1:v]scale=400:-1[logo];[main][logo]overlay=W-w-10:10[out] " +
				"-map [out] -c:v libx264 -crf 22 out.mp4"
			assert.Equal(t, expected, builder.Build())
		})

		t.Run("Simple filter chain with -vf", func(t *testing.T) {
			builder := New().
				Input(in).
				Filter().
				Simple(FilterVideo).
				Add(AtomicFilter{Name: "scale", Params: []string{"640", "-1"}}).
				Add(AtomicFilter{Name: "hflip"}).
				Done().
				Output(out)

			expected := "ffmpeg -i video.mp4 -vf scale=640:-1,hflip out.mp4"
			assert.Equal(t, expected, builder.Build())
		})

		t.Run("Simple filter chain with -af", func(t *testing.T) {
			builder := New().
				Input(in).
				Filter().
				Simple(FilterAudio).
				Add(AtomicFilter{Name: "volume", Params: []string{"0.5"}}).
				Add(AtomicFilter{Name: "atempo", Params: []string{"2.0"}}).
				Done().
				Output(out)

			expected := "ffmpeg -i video.mp4 -af volume=0.5,atempo=2.0 out.mp4"
			assert.Equal(t, expected, builder.Build())
		})

		t.Run("Real-world complex filter", func(t *testing.T) {
			builder := New().
				Override().
				Ss(10*time.Second).
				Input("input.mp4").
				Input("logo.png").
				Filter().
				Complex().
				Chain(
					[]string{"0:v"},
					AtomicFilter{Name: "crop", Params: []string{"1280", "720", "0", "0"}},
					[]string{"cropped"},
				).
				Chain(
					[]string{"1:v"},
					AtomicFilter{Name: "scale", Params: []string{"200", "-1"}},
					[]string{"logo_scaled"},
				).
				Chain(
					[]string{"cropped", "logo_scaled"},
					AtomicFilter{Name: "overlay", Params: []string{"W-w-10", "10"}},
					[]string{"video_out"},
				).
				Chain(
					[]string{"0:a"},
					AtomicFilter{Name: "atempo", Params: []string{"1.1"}},
					[]string{"audio_out"},
				).
				Done().
				Map("video_out").
				Map("audio_out").
				VideoCodec("libx264").
				AudioCodec("aac").
				Preset("fast").
				CRF(23).
				Output("output.mp4")

			expected := "ffmpeg -y -ss 00:00:10.000 -i input.mp4 -i logo.png " +
				"-filter_complex [0:v]crop=1280:720:0:0[cropped];[1:v]scale=200:-1[logo_scaled];[cropped][logo_scaled]overlay=W-w-10:10[video_out];[0:a]atempo=1.1[audio_out] " +
				"-map [video_out] -map [audio_out] -c:v libx264 -c:a aac -preset fast -crf 23 output.mp4"
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
