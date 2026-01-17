package ffmpeg

type FilterStage interface {
	// Video(expr string) FilterStage
	// Audio(expr string) FilterStage
	// Complex(expr string) FilterStage
	//
	// Output(path string) WriteStage
}

type filterCtx struct{ b *ffmpegBuilder }
