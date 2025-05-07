package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"solution/pkg/ipcounter"

	_ "github.com/joho/godotenv/autoload"

	"github.com/kelseyhightower/envconfig"
)

type config struct {
	Input     string `envconfig:"INPUT"`
	Mode      string `envconfig:"MODE"`
	Host      string `envconfig:"HOST"`
	Port      int    `envconfig:"PORT"`
	Namespace string `envconfig:"NAMESPACE"`
	Set       string `envconfig:"SET"`
	Output    string `envconfig:"OUTPUT"`
}

func main() {
	ctx := context.Background()

	var c config
	if err := envconfig.Process("", &c); err != nil {
		log.Printf("Failed to load config: %v", err)
		return
	}

	storage, err := ipcounter.NewAeroSpike(c.Host, c.Port)
	if err != nil {
		log.Printf("Failed to initialize storage: %v", err)
		return
	}
	defer storage.Close()

	service := ipcounter.NewService(storage)
	handler := &ipcounter.FileHandler{FS: os.DirFS(".")}

	reader, err := handler.OpenRead(c.Input)
	if err != nil {
		log.Printf("Failed to read input file: %v", err)
		return
	}
	defer reader.Close()

	elapsed, err := service.Import(ctx, reader, c.Namespace, c.Set, c.Mode)
	if err != nil {
		log.Printf("Failed to read input: %v", err)
		return
	}

	content, err := service.Export(ctx, c.Namespace, c.Set, c.Mode)
	if err != nil {
		log.Printf("Failed to scan: %v", err)
		return
	}

	content += fmt.Sprintf("\nTotal time %d seconds", elapsed)
	writer, err := handler.CreateWrite(c.Output, content)
	if err != nil {
		log.Printf("Failed to write output file: %v", err)
		return
	}
	defer writer.Close()

	log.Println("Done")
}
