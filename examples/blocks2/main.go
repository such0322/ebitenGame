package main

import (
	"flag"
	"log"
	"os"
	"runtime/pprof"

	"ebitenGame/examples/blocks2/blocks"
	"github.com/hajimehoshi/ebiten/v2"
)

var cpuProfile = flag.String("cpuprofile", "", "write cpu profile to file")

func main() {
	flag.Parse()
	if *cpuProfile != "" {
		f, err := os.Create(*cpuProfile)
		if err != nil {
			log.Fatal(err)
		}
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal(err)
		}
		defer pprof.StopCPUProfile()
	}

	ebiten.SetWindowSize(blocks.ScreenWidth*2, blocks.ScreenHeight*2)
	ebiten.SetWindowTitle("Blocks (Ebiten Demo)")
	if err := ebiten.RunGame(&blocks.Game{}); err != nil {
		log.Fatal(err)
	}
}
