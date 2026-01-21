package fflow

import (
	"fmt"
	"strings"
)

// SimpleFilterType representa filtros simples aplicados a um único stream. FilterVideo para -vf ou FilterAudio para -af.
//
// SimpleFilterType represents simple filters applied to a single stream. FilterVideo to -vf or FilterAudio to -af.
type SimpleFilterType string

const (
	// FilterVideo representa filtros simples de vídeo (-vf).
	//
	// FilterVideo represents simple video filters (-vf).
	FilterVideo SimpleFilterType = "-vf"

	// FilterAudio representa filtros simples de áudio (-af).
	//
	// FilterAudio represents simple audio filters (-af).
	FilterAudio SimpleFilterType = "-af"
)

type filter interface {
	// String retorna a representação textual do filtro no formato aceito pelo ffmpeg.
	//
	// String returns the textual representation of the filter in a format accepted by ffmpeg.
	String() string

	// NeedsComplex indica se o filtro requer o uso de -filter_complex.
	//
	// NeedsComplex indicates whether the filter requires -filter_complex.
	NeedsComplex() bool
}
type filterStage interface {
	// Simple inicia a construção de filtros simples (-vf ou -af),
	// aplicados diretamente a um único stream.
	//
	// Simple starts building simple filters (-vf or -af),
	// applied directly to a single stream.
	Simple(t SimpleFilterType) simpleFilter

	// Complex inicia a construção de filtros complexos (-filter_complex),
	// permitindo múltiplos inputs e outputs.
	//
	// Complex starts building complex filters (-filter_complex),
	// allowing multiple inputs and outputs.
	Complex() complexFilter
}

type simpleFilter interface {
	// Add adiciona um filtro atômico à cadeia de filtros simples.
	//
	// Add appends an atomic filter to the simple filter chain.
	Add(filter AtomicFilter) simpleFilter

	// Done finaliza a construção dos filtros simples
	// e avança para o próximo estágio do pipeline.
	//
	// Done finalizes the simple filter construction
	// and advances to the next pipeline stage.
	Done() writeStage
}

type complexFilter interface {
	// Chain adiciona um filtro complexo à cadeia,
	// conectando explicitamente inputs e outputs.
	//
	// Chain adds a complex filter to the chain,
	// explicitly connecting inputs and outputs.
	Chain(in []string, filter []AtomicFilter, out []string) complexFilter

	// Done finaliza a construção dos filtros complexos
	// e avança para o próximo estágio do pipeline.
	//
	// Done finalizes the complex filter construction
	// and advances to the next pipeline stage.
	Done() writeStage
}

type (
	filterCtx        struct{ b *ffmpegBuilder }
	simpleFilterCtx  struct{ b *ffmpegBuilder }
	complexFilterCtx struct{ b *ffmpegBuilder }
)

func (c *filterCtx) Simple(t SimpleFilterType) simpleFilter {
	c.b.simpleFilterFlag = string(t)
	return &simpleFilterCtx{c.b}
}

func (c *filterCtx) Complex() complexFilter {
	return &complexFilterCtx{c.b}
}

func (sf *simpleFilterCtx) Add(filter AtomicFilter) simpleFilter {
	sf.b.filters = append(sf.b.filters, filter)
	return sf
}

func (sf *simpleFilterCtx) Done() writeStage {
	return &writeCtx{sf.b}
}

func (cf *complexFilterCtx) Chain(in []string, filter []AtomicFilter, out []string) complexFilter {
	chain := Chain{Inputs: in, Filter: filter, Output: out}
	cf.b.filters = append(cf.b.filters, chain)
	return cf
}

func (cf *complexFilterCtx) Done() writeStage {
	return &writeCtx{cf.b}
}

type AtomicFilter struct {
	Name   string
	Params []string
}

func (f AtomicFilter) String() string {
	if len(f.Params) == 0 {
		return f.Name
	}
	return fmt.Sprintf("%s=%s", f.Name, strings.Join(f.Params, ":"))
}

func (f AtomicFilter) NeedsComplex() bool {
	return false
}

type Chain struct {
	Inputs []string
	Filter []AtomicFilter
	Output []string
}

func (c Chain) String() string {
	var sb strings.Builder

	for _, in := range c.Inputs {
		sb.WriteString("[")
		sb.WriteString(in)
		sb.WriteString("]")
	}

	for i, f := range c.Filter {
		if i > 0 {
			sb.WriteString(",")
		}
		sb.WriteString(f.String())
	}

	for _, out := range c.Output {
		sb.WriteString("[")
		sb.WriteString(out)
		sb.WriteString("]")
	}

	return sb.String()
}

func (c Chain) NeedsComplex() bool {
	return true
}

type Pipeline struct {
	Nodes []filter
}

func (p Pipeline) String() string {
	var parts []string
	for _, n := range p.Nodes {
		parts = append(parts, n.String())
	}
	if p.NeedsComplex() {
		return strings.Join(parts, ";")
	}

	return strings.Join(parts, ",")
}

func (p Pipeline) NeedsComplex() bool {
	for _, n := range p.Nodes {
		if n.NeedsComplex() {
			return true
		}
	}
	return false
}
