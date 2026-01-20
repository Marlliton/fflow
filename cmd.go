package fflow

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"time"
)

type commandStage interface {
	// String retorna o comando final (debug/log).
	//
	// String returns the final command (debug/log).
	String() string

	// Cmd retorna *exec.Cmd pronto para executar.
	//
	// Cmd returns *exec.Cmd ready to execute.
	Cmd(ctx context.Context) *exec.Cmd

	// Run Executa o comando
	//
	// Run execute the command
	Run(ctx context.Context) error

	// RunWithProgress executa o comando e emite eventos de progresso.
	//
	// RunWithProgress executes the command and emits progress events.
	RunWithProgress(ctx context.Context) (<-chan Progress, error)
}

type Progress struct {
	Frame   int
	FPS     float64
	Bitrate string
	OutTime time.Duration
	Speed   string
}

type commandCtx struct{ b *ffmpegBuilder }

func (c *commandCtx) String() string {
	return c.tmpWritter().String()
}

func (c *commandCtx) Cmd(ctx context.Context) *exec.Cmd {
	return exec.CommandContext(ctx, "ffmpeg", c.tmpWritter().Args()...)
}

func (c *commandCtx) Run(ctx context.Context) error {
	cmd := c.Cmd(ctx)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (c *commandCtx) RunWithProgress(ctx context.Context) (<-chan Progress, error) {
	args := c.tmpWritter().Args()
	args = append(args, "-progress", "pipe:2", "-nostats", "-loglevel", "error")

	cmd := exec.CommandContext(ctx, "ffmpeg", args...)
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, err
	}

	err = cmd.Start()
	if err != nil {
		return nil, err
	}

	pch := make(chan Progress)

	go func() {
		defer close(pch)

		c.monitorProgress(stderr, pch)

		if err := cmd.Wait(); err != nil {
			fmt.Printf("FFmpeg terminou com erro: %v\n", err)
		}
	}()

	return pch, nil
}

func (c commandCtx) monitorProgress(stderr io.ReadCloser, pch chan Progress) {
	scanner := bufio.NewScanner(stderr)

	prog := Progress{}
	for scanner.Scan() {
		line := scanner.Text()

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := parts[0]
		value := parts[1]

		switch key {
		case "frame":
			fmt.Sscanf(value, "%d", &prog.Frame)

		case "fps":
			fmt.Sscanf(value, "%f", &prog.FPS)

		case "bitrate":
			prog.Bitrate = value

		case "out_time":
			var h, m int
			var sec float64
			fmt.Sscanf(value, "%d:%d:%f", &h, &m, &sec)
			d := time.Duration(h)*time.Hour +
				time.Duration(m)*time.Minute +
				time.Duration(sec*float64(time.Second))

			prog.OutTime = d

		case "speed":
			s := strings.TrimSuffix(value, "x")
			fmt.Sscanf(s, "%s", &prog.Speed)

		case "progress":
			pch <- prog

			if value == "end" {
				return
			}
		}
	}
}

func (c *commandCtx) tmpWritter() *writeCtx {
	return (&writeCtx{c.b})
}
