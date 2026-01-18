package ffmpeg

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestReadStage(t *testing.T) {
	t.Run("Definições de tempo de leitura", func(t *testing.T) {
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
				name: "seek, procura tempo de vídeo com -ss (antes do -i)",
				builder: New().
					Input("movie.mkv").
					Ss(22 * time.Second).
					Output("out.mkv"),
				expected: "ffmpeg -ss 00:00:22.000 -i movie.mkv out.mkv",
			},
			{
				name: "Define tempo final absoluto do vídeo (antes do -i)",
				builder: New().
					Input("movie.mkv").
					To(52 * time.Second).
					Output("out.mkv"),
				expected: "ffmpeg -to 00:00:52.000 -i movie.mkv out.mkv",
			},
			{
				name: "Argumento bruto (antes do -i)",
				builder: New().
					Raw("-re").
					Input("movie.mkv").
					Output("out.mkv"),
				expected: "ffmpeg -re -i movie.mkv out.mkv",
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

		// cmd := New().
		// 	Input("movie.mkv").
		// 	Ss(22 * time.Second).
		// 	To(60 * time.Second).
		// 	Input("movie2.mkv").
		// 	Ss(120*time.Second).
		// 	To(630*time.Second).
		// 	Filter().
		// 	Simple().Add(AtomicFilter{Name: "scale", Params: []string{"v:"}})
	})
}
