package ffmpeg

import "time"

type ReadStage interface {
	// Ss adiciona a flag -ss antes do -i, realizando um seek rápido na entrada.
	//
	// Ss adds the -ss flag before -i, performing a fast seek on the input.
	Ss(d time.Duration) ReadStage

	// To adiciona a flag -to antes do -i, definindo o tempo final absoluto da leitura.
	//
	// To adds the -to flag before -i, defining the absolute end time of the input read.
	To(d time.Duration) ReadStage

	// T adiciona a flag -t antes do -i, limitando quanto da entrada será lida.
	//
	// T adds the -t flag before -i, limiting how much of the input is read.
	T(d time.Duration) ReadStage

	// Input adiciona um arquivo de entrada (-i).
	//
	// Input adds an input file (-i).
	Input(path string) ReadStage

	// Filter transiciona para a etapa de filtros da entrada atual.
	//
	// Filter transitions to the filter stage for the current input.
	Filter() FilterStage

	// Output define o arquivo de saída e transiciona para o WriteStage.
	//
	// Output sets the output file and transitions to WriteStage.
	Output(path string) WriteStage
}

type readCtx struct{ b *ffmpegBuilder }

func (c *readCtx) T(d time.Duration) ReadStage {
	c.b.read = append([]string{"-t", fmtDuration(d)}, c.b.read...)
	return c
}

func (c *readCtx) Ss(d time.Duration) ReadStage {
	c.b.read = append([]string{"-ss", fmtDuration(d)}, c.b.read...)
	return c
}

func (c *readCtx) To(d time.Duration) ReadStage {
	c.b.read = append([]string{"-to", fmtDuration(d)}, c.b.read...)
	return c
}

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
