package main

import (
	"embed"
	"log"

	"go-walis/internal/config"
	"go-walis/internal/containers"
	"go-walis/internal/core/db"
	"go-walis/internal/extras"
	"go-walis/internal/images"
	"go-walis/internal/networks"
	"go-walis/internal/volumes"

	"github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	database, err := db.InitDB()
	if err != nil {
		log.Fatalf("Erro ao inicializar banco SQLite: %v", err)
	}
	defer database.Close()

	configService := config.NewConfigService(database)
	containerService := containers.NewContainerService(database)
	extraService := extras.NewExtraService(database)
	imageService := images.NewImageService(database)
	networkService := networks.NewNetworkService(database)
	volumeService := volumes.NewVolumeService(database)

	app := application.New(application.Options{
		Name:        "DockSea",
		Description: "DockSea Container Manager",
		Services: []application.Service{
			application.NewService(configService),
			application.NewService(containerService),
			application.NewService(extraService),
			application.NewService(imageService),
			application.NewService(networkService),
			application.NewService(volumeService),
		},
		Assets: application.AssetOptions{
			Handler: application.AssetFileServerFS(assets),
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
	})

	app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:  "DockSea",
		Width:  1200,
		Height: 750,
		Mac: application.MacWindow{
			InvisibleTitleBarHeight: 50,
			Backdrop:                application.MacBackdropTranslucent,
			TitleBar:                application.MacTitleBarHiddenInset,
		},
		BackgroundColour: application.NewRGB(15, 23, 42),
		URL:              "/",
	})

	err = app.Run()
	if err != nil {
		log.Fatal(err)
	}
}

