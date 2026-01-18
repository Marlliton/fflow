package ffmpeg

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
	String() string
	NeedsComplex() bool
}
type filterStage interface {
	Simple(t SimpleFilterType) simpleFilter
	Complex() complexFilter
}

type simpleFilter interface {
	Add(filter AtomicFilter) simpleFilter
	Done() WriteStage
}

type complexFilter interface {
	Chaing(in []string, filter AtomicFilter, out []string) complexFilter
	Done() WriteStage
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

func (sf *simpleFilterCtx) Done() WriteStage {
	return &writeCtx{sf.b}
}

func (cf *complexFilterCtx) Chaing(in []string, filter AtomicFilter, out []string) complexFilter {
	chain := Chaing{Inputs: in, Filter: filter, Output: out}
	cf.b.filters = append(cf.b.filters, chain)
	return cf
}

func (cf *complexFilterCtx) Done() WriteStage {
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

type Chaing struct {
	Inputs []string
	Filter AtomicFilter
	Output []string
}

func (c Chaing) String() string {
	var sb strings.Builder

	for _, in := range c.Inputs {
		sb.WriteString("[")
		sb.WriteString(in)
		sb.WriteString("]")
	}

	sb.WriteString(c.Filter.String())

	for _, out := range c.Output {
		sb.WriteString("[")
		sb.WriteString(out)
		sb.WriteString("]")
	}

	return sb.String()
}

func (c Chaing) NeedsComplex() bool {
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
