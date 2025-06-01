package main

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/xperimental/gosha/internal/config"
	"github.com/xperimental/gosha/internal/digest"
)

const (
	megaBytes = 1024 * 1024
)

func main() {
	log.SetFlags(0)

	cfg, err := config.Parse(os.Args[0], os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}

	wg := &sync.WaitGroup{}
	fileCh := make(chan string, len(cfg.FileNames))
	resCh := make(chan digest.Result, len(cfg.FileNames))

	for _, fileName := range cfg.FileNames {
		info, err := os.Stat(fileName)
		if err != nil {
			log.Fatalf("can not stat file: %s", err)
		}

		if info.IsDir() {
			log.Printf("Skipping directory: %s", fileName)
			continue
		}

		fileCh <- fileName
	}
	close(fileCh)

	workers := min(cfg.Workers, len(cfg.FileNames))
	for i := 0; i < workers; i++ {
		digest.Worker(wg, resCh, fileCh)
	}

	wg.Wait()
	close(resCh)

	totalBytes := int64(0)
	totalDuration := time.Duration(0)
	for res := range resCh {
		if res.Error != nil {
			log.Fatal(res.Error)
		}

		totalBytes += res.Size
		totalDuration += res.Duration

		fmt.Printf("%s *%s\n", res.Hash, res.Filename)
	}

	fmt.Printf("Total Bytes: %d Duration: %s Speed: %.2fMB/s\n", totalBytes, totalDuration.Round(time.Second), (float64(totalBytes)/megaBytes)/float64(totalDuration.Seconds()))
}
