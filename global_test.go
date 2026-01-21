package fflow

import (
	"testing"
)

func TestGlobalStage(t *testing.T) {
	in := "video.mp4"
	out := "out.mp4"

	t.Run("Argumentos básicos", func(t *testing.T) {
		run(t, []testCase{
			{
				name:     "Input e Output",
				builder:  New().Input(in).Output(out),
				expected: "ffmpeg -loglevel error -y -i video.mp4 out.mp4",
			},
		})
	})
}
