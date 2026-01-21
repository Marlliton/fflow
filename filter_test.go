package fflow

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFilterElements(t *testing.T) {
	scale := "scale"

	t.Run("AtomicFilter String and NeedsComplex", func(t *testing.T) {
		tests := []struct {
			name         string
			filter       AtomicFilter
			expectedStr  string
			needsComplex bool
		}{
			{
				name:         "Scale filter with params",
				filter:       AtomicFilter{Name: scale, Params: []string{"1280", "-1"}},
				expectedStr:  "scale=1280:-1",
				needsComplex: false,
			},
			{
				name:         "Hflip filter without params",
				filter:       AtomicFilter{Name: "hflip", Params: []string{}},
				expectedStr:  "hflip",
				needsComplex: false,
			},
			{
				name:         "Empty filter name",
				filter:       AtomicFilter{Name: "", Params: []string{"param"}},
				expectedStr:  "=param",
				needsComplex: false,
			},
		}

		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				assert.Equal(t, tc.expectedStr, tc.filter.String())
				assert.Equal(t, tc.needsComplex, tc.filter.NeedsComplex())
			})
		}
	})

	t.Run("Chaing String and NeedsComplex", func(t *testing.T) {
		tests := []struct {
			name        string
			inputs      []string
			filter      []AtomicFilter
			output      []string
			expectedStr string
		}{
			{
				name:        "Single input, scale filter, single output label",
				inputs:      []string{"0:v"},
				filter:      []AtomicFilter{{Name: scale, Params: []string{"1280", "-1"}}},
				output:      []string{"out"},
				expectedStr: "[0:v]scale=1280:-1[out]",
			},
			{
				name:        "Multiple inputs, overlay filter, single output label",
				inputs:      []string{"main", "logo"},
				filter:      []AtomicFilter{{Name: "overlay", Params: []string{"W-w-10:10"}}},
				output:      []string{"final_video"},
				expectedStr: "[main][logo]overlay=W-w-10:10[final_video]",
			},
			{
				name:        "No output label",
				inputs:      []string{"0:v"},
				filter:      []AtomicFilter{{Name: scale, Params: []string{"640", "-1"}}},
				output:      []string{},
				expectedStr: "[0:v]scale=640:-1",
			},
			{
				name:        "Single input, split filter, multiple output labels",
				inputs:      []string{"0:v"},
				filter:      []AtomicFilter{{Name: "split", Params: []string{"2"}}},
				output:      []string{"main", "blur"},
				expectedStr: "[0:v]split=2[main][blur]",
			},
		}

		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				c := Chain{Inputs: tc.inputs, Filter: tc.filter, Output: tc.output}
				assert.Equal(t, tc.expectedStr, c.String())
				assert.True(t, c.NeedsComplex())
			})
		}
	})

	t.Run("Pipeline String and NeedsComplex", func(t *testing.T) {
		tests := []struct {
			name         string
			nodes        []filter
			expectedStr  string
			needsComplex bool
		}{
			{
				name: "Pipeline of simple atomic filters (treated as separate chains)",
				nodes: []filter{
					AtomicFilter{Name: scale, Params: []string{"1280", "-1"}},
					AtomicFilter{Name: "hflip", Params: []string{}},
				},
				expectedStr:  "scale=1280:-1,hflip",
				needsComplex: false,
			},
			{
				name: "Pipeline including a complex chain",
				nodes: []filter{
					AtomicFilter{Name: "format", Params: []string{"yuv420p"}},
					Chain{Inputs: []string{"0:v"}, Filter: []AtomicFilter{{Name: "fade", Params: []string{"in", "0", "30"}}}, Output: []string{"faded_video"}},
					AtomicFilter{Name: "setsar", Params: []string{"1"}},
				},
				expectedStr:  "format=yuv420p;[0:v]fade=in:0:30[faded_video];setsar=1",
				needsComplex: true, // INFO: One node needs complex, so pipeline needs complex
			},
			{
				name: "Pipeline of only complex chains",
				nodes: []filter{
					Chain{Inputs: []string{"0:v"}, Filter: []AtomicFilter{{Name: scale, Params: []string{"640", "-1"}}}, Output: []string{"scaled"}},
					Chain{Inputs: []string{"scaled", "1:v"}, Filter: []AtomicFilter{{Name: "overlay", Params: []string{"W-w-10:10"}}}, Output: []string{"final"}},
				},
				expectedStr:  "[0:v]scale=640:-1[scaled];[scaled][1:v]overlay=W-w-10:10[final]",
				needsComplex: true, // INFO: All nodes need complex
			},
			{
				name:         "Empty pipeline",
				nodes:        []filter{},
				expectedStr:  "",
				needsComplex: false,
			},
		}

		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				p := Pipeline{Nodes: tc.nodes}
				assert.Equal(t, tc.expectedStr, p.String())
				assert.Equal(t, tc.needsComplex, p.NeedsComplex())
			})
		}
	})
}

func TestFilterStages(t *testing.T) {
	t.Run("filterCtx (Entry Point)", func(t *testing.T) {
		b := &ffmpegBuilder{}
		fCtx := filterCtx{b}

		t.Run("Simple returns non-nil SimpleFilter", func(t *testing.T) {
			simpleStage := fCtx.Simple(FilterVideo)
			assert.NotNil(t, simpleStage)
			_, ok := simpleStage.(*simpleFilterCtx)
			assert.True(t, ok, "Simple() should return a *simpleFilterCtx")
		})

		t.Run("Complex returns non-nil ComplexFilter", func(t *testing.T) { // Corrected
			complexStage := fCtx.Complex()
			assert.NotNil(t, complexStage)
			_, ok := complexStage.(*complexFilterCtx)
			assert.True(t, ok, "Complex() should return a *complexFilterCtx")
		})
	})

	t.Run("simpleFilterCtx (Simple Filter Builder)", func(t *testing.T) {
		b := &ffmpegBuilder{}
		sCtx := simpleFilterCtx{b}

		t.Run("Add appends AtomicFilter to builder.filters", func(t *testing.T) { // Changed description
			filter1 := AtomicFilter{Name: "scale", Params: []string{"1280", "-1"}}
			filter2 := AtomicFilter{Name: "hflip"}

			sCtx.Add(filter1)
			sCtx.Add(filter2)

			expectedFilters := []filter{filter1, filter2} // Now []Filter
			assert.Equal(t, expectedFilters, b.filters)
		})

		t.Run("Done returns non-nil WriteStage", func(t *testing.T) {
			writeStage := sCtx.Done()
			assert.NotNil(t, writeStage)
			_, ok := writeStage.(*writeCtx)
			assert.True(t, ok, "Done() should return a *writeCtx")
		})
	})

	t.Run("complexFilterCtx (Complex Filter Builder)", func(t *testing.T) {
		// Instantiate fCtx to obtain cCtx correctly
		b := &ffmpegBuilder{}
		fCtx := filterCtx{b}
		cCtx := fCtx.Complex().(*complexFilterCtx) // Get the instance of complexFilterCtx

		t.Run("Chaing appends Chaing object to builder.filters", func(t *testing.T) { // Changed description
			chain1 := Chain{
				Inputs: []string{"0:v"},
				Filter: []AtomicFilter{{Name: "scale", Params: []string{"640", "-1"}}},
				Output: []string{"out"},
			}
			chain2 := Chain{
				Inputs: []string{"1:a"},
				Filter: []AtomicFilter{{Name: "aformat", Params: []string{"fltp"}}},
				Output: []string{},
			}
			cCtx.Chain(chain1.Inputs, chain1.Filter, chain1.Output)
			cCtx.Chain(chain2.Inputs, chain2.Filter, chain2.Output)

			expectedFilters := []filter{chain1, chain2} // Now []Filter
			assert.Equal(t, expectedFilters, b.filters)
		})

		t.Run("Done returns non-nil WriteStage", func(t *testing.T) {
			// Using the same cCtx from above setup
			writeStage := cCtx.Done()
			assert.NotNil(t, writeStage)
			_, ok := writeStage.(*writeCtx)
			assert.True(t, ok, "Done() should return a *writeCtx")
		})
	})
}
