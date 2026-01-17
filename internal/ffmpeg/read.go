package ffmpeg

import "time"

type ReadStage interface {
	Ss(d time.Duration) ReadStage
	T(d time.Duration) ReadStage
	Input(path string) ReadStage

	Filter() FilterStage

	Output(path string) WriteStage
}

type readCtx struct{ b *ffmpegBuilder }

// T adiciona a flag -t (duração).
// Define por quanto tempo o FFmpeg processa a partir do início
// ou do ponto definido por -ss.
// Ex: T(30 * time.Second) → 30s
func (c *readCtx) T(d time.Duration) ReadStage {
	c.b.read = append([]string{"-t", fmtDuration(d)}, c.b.read...)
	return c
}

// Ss adiciona a flag -ss (seek).
// Antes do -i: seek rápido (menos preciso).
// Depois do -i: seek preciso (mais lento).
func (c *readCtx) Ss(d time.Duration) ReadStage {
	c.b.read = append([]string{"-ss", fmtDuration(d)}, c.b.read...)
	return c
}

// To adiciona a flag -to (tempo final absoluto).
// Processa até atingir o timestamp informado.
// Geralmente usado em conjunto com -ss.
func (c *readCtx) To(d time.Duration) ReadStage {
	c.b.read = append([]string{"-to", fmtDuration(d)}, c.b.read...)
	return c
}

// Input adiciona um arquivo de entrada (-i).
func (c *readCtx) Input(path string) ReadStage {
	c.b.read = append(c.b.read, "-i", path)
	return c
}

func (c *readCtx) Filter() FilterStage {
	return &filterCtx{c.b}
}

func (c *readCtx) Output(path string) WriteStage {
	write := &writeCtx{c.b}
	write.Output(path)
	return write
}
