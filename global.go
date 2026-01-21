package fflow

import "time"

type beforeReadStage interface {
	// Raw adiciona um argumento bruto ao comando FFmpeg, antes do -i
	//
	// Raw adds a raw argument to the FFmpeg command, before -i flag
	Raw(value string) beforeReadStage

	// Input adiciona um arquivo de entrada (-i) e transiciona para o ReadStage.
	//
	// Input adds an input file (-i) and transitions to ReadStage.
	Input(path string) readStagee

	// Ss adiciona a flag -ss antes do -i, realizando um seek rápido na entrada.
	//
	// Ss adds the -ss flag before -i, performing a fast seek on the input.
	Ss(d time.Duration) beforeReadStage

	// To adiciona a flag -to antes do -i, definindo o tempo final absoluto da leitura.
	//
	// To adds the -to flag before -i, defining the absolute end time of the input read.
	To(d time.Duration) beforeReadStage

	// T adiciona a flag -t antes do -i, limitando quanto da entrada será lida.
	//
	// T adds the -t flag before -i, limiting how much of the input is read.
	T(d time.Duration) beforeReadStage
}

type beforeReadCtx struct{ b *ffmpegBuilder }

func (c *beforeReadCtx) Input(path string) readStagee {
	read := &readCtx{c.b}
	read.Input(path)
	return read
}

func (c *beforeReadCtx) Raw(value string) beforeReadStage {
	c.b.beforeRead = append(c.b.beforeRead, value)
	return c
}

func (c *beforeReadCtx) T(d time.Duration) beforeReadStage {
	c.b.beforeRead = append(c.b.beforeRead, "-t", fmtDuration(d))
	return c
}

func (c *beforeReadCtx) Ss(d time.Duration) beforeReadStage {
	c.b.beforeRead = append(c.b.beforeRead, "-ss", fmtDuration(d))
	return c
}

func (c *beforeReadCtx) To(d time.Duration) beforeReadStage {
	c.b.beforeRead = append(c.b.beforeRead, "-to", fmtDuration(d))
	return c
}
