package ffmpeg

import (
	"testing"

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

func TestFFmpegBuilder(t *testing.T) {
	t.Run("Fluent interface mantém o mesmo estado", func(t *testing.T) {
		f := New()
		f2 := f.Input("test.mp4")

		assert.NotNil(t, f2)
	})
}
