# ffcmd - Construtor de Comandos FFmpeg Fluente para Go

[![Go Report Card](https://goreportcard.com/badge/github.com/Marlliton/ffcmd)](https://goreportcard.com/report/github.com/Marlliton/ffcmd)

`ffcmd` é uma biblioteca Go que fornece uma interface fluente e semântica para construir comandos `ffmpeg` de forma programática. Diga adeus à concatenação de strings e aos erros de ordenamento de flags.

## ✨ Recursos

- **API Fluente**: Construa comandos complexos encadeando métodos de forma legível.
- **Ordem Semântica**: A biblioteca garante a ordem correta das flags do FFmpeg (opções globais, de entrada e de saída).
- **Filtros Simples e Complexos**: Suporte nativo para `-vf`, `-af` e `-filter_complex` de forma organizada.
- **Type-Safe**: Evite erros comuns especificando se um filtro simples é para **vídeo** ou **áudio**.
- **Clareza**: Separação clara entre estágios de configuração (Global, Leitura, Filtro, Escrita).

## 📦 Instalação

```bash
go get github.com/Marlliton/ffcmd/ffmpeg
```

## 🚀 Uso e Exemplos

A utilização da biblioteca segue a lógica de construção de um comando `ffmpeg`: primeiro as opções globais, depois as entradas, os filtros e, por fim, a saída e suas opções.

### Exemplo 1: Conversão Básica

Converter um arquivo de vídeo para um formato diferente.

```go
package main

import (
	"fmt"
	"github.com/Marlliton/ffcmd/ffmpeg"
)

func main() {
	cmd := ffmpeg.New().
		Override(). // Adiciona a flag global -y para sobrescrever o arquivo de saída
		Input("input.mp4").
		Output("output.webm").
		VideoCodec("libvpx-vp9").
		AudioCodec("libopus").
		Build()

	fmt.Println(cmd)
	// Saída: ffmpeg -y -i input.mp4 -c:v libvpx-vp9 -c:a libopus output.webm
}
```

### Exemplo 2: Cortar um Vídeo (Trimming)

É possível usar `Ss` (seek) e `T` (duração) tanto na entrada (para um seek rápido) quanto na saída (para um corte preciso).

```go
package main

import (
	"fmt"
	"time"
	"github.com/Marlliton/ffcmd/ffmpeg"
)

func main() {
    // Seek rápido na entrada e duração limitada
	cmd := ffmpeg.New().
		Ss(1 * time.Minute). // Pula para 1 minuto do início do arquivo de entrada
		T(30 * time.Second).  // Lê apenas 30 segundos da entrada
		Input("input.mp4").
		Output("output.mp4").
		CopyVideo(). // Copia o stream de vídeo sem recodificar
		CopyAudio().  // Copia o stream de áudio sem recodificar
		Build()

	fmt.Println(cmd)
	// Saída: ffmpeg -ss 00:01:00.000 -t 00:00:30.000 -i input.mp4 -c:v copy -c:a copy output.mp4
}
```

### Exemplo 3: Filtros Simples (Vídeo e Áudio)

Aplique filtros a um único stream de vídeo (`-vf`) ou áudio (`-af`).

```go
package main

import (
	"fmt"
	"github.com/Marlliton/ffcmd/ffmpeg"
)

func main() {
    // Filtro de vídeo para redimensionar e inverter horizontalmente
	videoCmd := ffmpeg.New().
		Input("input.mp4").
		Filter().
		Simple(ffmpeg.FilterVideo). // Especifica que é um filtro de vídeo (-vf)
		Add(ffmpeg.AtomicFilter{Name: "scale", Params: []string{"1280", "-1"}}).
		Add(ffmpeg.AtomicFilter{Name: "hflip"}).
		Done().
		Output("video_filtered.mp4").
		Build()

	fmt.Println(videoCmd)
	// Saída: ffmpeg -i input.mp4 -vf scale=1280:-1,hflip video_filtered.mp4

    // Filtro de áudio para ajustar o volume
	audioCmd := ffmpeg.New().
		Input("input.mp4").
		Filter().
		Simple(ffmpeg.FilterAudio). // Especifica que é um filtro de áudio (-af)
		Add(ffmpeg.AtomicFilter{Name: "volume", Params: []string{"0.5"}}).
		Done().
		Output("audio_filtered.mp3").
		Build()

	fmt.Println(audioCmd)
    // Saída: ffmpeg -i input.mp4 -af volume=0.5 audio_filtered.mp3
}
```

### Exemplo 4: Filtro Complexo (Cenário Real)

Um exemplo mais avançado: cortar um vídeo, sobrepor uma marca d'água, acelerar o áudio e re-codificar com presets específicos.

```go
package main

import (
	"fmt"
	"time"
	"github.com/Marlliton/ffcmd/ffmpeg"
)

func main() {
	cmd := ffmpeg.New().
		Override().
		Input("input.mp4"). // Entrada de vídeo principal
		Input("logo.png").   // Entrada da imagem para a marca d'água
		Filter().
		Complex(). // Inicia um -filter_complex
		Chaing(
			[]string{"0:v"}, // Pega o vídeo da primeira entrada
			ffmpeg.AtomicFilter{Name: "scale", Params: []string{"1920", "-1"}},
			"scaled", // Nomeia a saída para uso posterior
		).
		Chaing(
			[]string{"scaled", "1:v"}, // Pega o vídeo redimensionado e a imagem da segunda entrada
			ffmpeg.AtomicFilter{Name: "overlay", Params: []string{"W-w-10", "10"}},
			"video_out", // Nomeia a saída de vídeo final
		).
		Chaing(
			[]string{"0:a"}, // Pega o áudio da primeira entrada
			ffmpeg.AtomicFilter{Name: "atempo", Params: []string{"1.5"}},
			"audio_out", // Nomeia a saída de áudio final
		).
		Done().
		Map("video_out"). // Mapeia a saída de vídeo do filtro complexo
		Map("audio_out"). // Mapeia a saída de áudio do filtro complexo
		VideoCodec("libx264").
		AudioCodec("aac").
		Preset("fast").
		CRF(23).
		Output("final_video.mp4").
		Build()

	fmt.Println(cmd)
	/*
	   Saída: ffmpeg -y -i input.mp4 -i logo.png -filter_complex [0:v]scale=1920:-1[scaled];[scaled][1:v]overlay=W-w-10:10[video_out];[0:a]atempo=1.5[audio_out] -map [video_out] -map [audio_out] -c:v libx264 -c:a aac -preset fast -crf 23 final_video.mp4
	*/
}
```

## 📖 Visão Geral da API

O builder é dividido em estágios para garantir uma construção lógica e semântica do comando.

1.  **`GlobalStage`**: Ponto de entrada (`New()`). Permite definir opções globais como `-y` (sobrescrever).
2.  **`ReadStage`**: Define as entradas (`Input()`) e suas opções, como `-ss` (seek) ou `-t` (duração).
3.  **`FilterStage`**: Permite a criação de filtros simples (`Simple()`) ou complexos (`Complex()`).
4.  **`WriteStage`**: Define a saída (`Output()`) e todas as suas opções, como codecs (`-c:v`), presets (`-preset`), CRF, etc. É o estágio final antes de construir o comando com `Build()`.
