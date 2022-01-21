package main

import (
	"github.com/fsnotify/fsnotify"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/spf13/viper"
	"log"
)

func init() {
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath("./conf/")
	viper.OnConfigChange(func(e fsnotify.Event) {
		log.Printf("config file changed: %s", e.Name)
	})
	viper.WatchConfig()

	if err := viper.ReadInConfig(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	g := NewGame()
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle(viper.GetString("name"))
	ebiten.SetFullscreen(false)
	err := ebiten.RunGame(g)
	if err != nil {
		log.Fatal(err)
	}
}
