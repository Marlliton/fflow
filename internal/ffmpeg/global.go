package ffmpeg

type GlobalStage interface {
	Override() GlobalStage
	// LogLevel(level string) GlobalStage

	Input(path string) ReadStage
}

type globalCtx struct{ b *ffmpegBuilder }

// Input adiciona um arquivo de entrada (-i).
func (c *globalCtx) Input(path string) ReadStage {
	read := &readCtx{c.b}
	read.Input(path)
	return read
}

// Override adiciona a flag de sobrescrita automática (-y).
func (c *globalCtx) Override() GlobalStage {
	c.b.global = append(c.b.global, "-y")
	return c
}
