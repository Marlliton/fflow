package ffmpeg

import (
	"fmt"
	"time"
)

// fmtDuration formata uma duração de tempo (time.Duration) para o formato de string HH:MM:SS.ms.
//
// Formats a time.Duration into an HH:MM:SS.ms string.
func fmtDuration(d time.Duration) string {
	h := int(d.Hours())
	m := int(d.Minutes()) % 60
	s := int(d.Seconds()) % 60
	ms := int(d.Milliseconds()) % 1000

	return fmt.Sprintf("%02d:%02d:%02d.%03d", h, m, s, ms)
}
