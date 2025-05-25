package main

import (
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"os"
	"sync"
	"time"
)

const (
	megaBytes = 1024 * 1024
)

type result struct {
	Filename string
	Hash     string
	Size     int64
	Duration time.Duration
	Error    error
}

func main() {
	log.SetFlags(0)
	if len(os.Args) < 2 {
		log.Fatalf("Usage: %s file [file ...]", os.Args[0])
	}

	wg := &sync.WaitGroup{}
	resCh := make(chan result, len(os.Args[1:]))

	for _, fileName := range os.Args[1:] {
		info, err := os.Stat(fileName)
		if err != nil {
			log.Fatalf("can not stat file: %s", err)
		}

		if info.IsDir() {
			log.Printf("Skipping directory: %s", fileName)
			continue
		}

		wg.Add(1)
		go func() {
			defer wg.Done()

			start := time.Now()
			log.Printf("Processing file %s", fileName)
			hash, bytes, err := hashFile(fileName)
			duration := time.Since(start)

			resCh <- result{
				Filename: fileName,
				Hash:     hash,
				Size:     bytes,
				Duration: duration,
				Error:    err,
			}
		}()
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

func hashFile(fileName string) (string, int64, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return "", 0, fmt.Errorf("can not open file: %w", err)
	}
	defer file.Close()

	hash := sha256.New()
	bytes, err := io.Copy(hash, file)
	if err != nil {
		return "", 0, fmt.Errorf("error while hashing file: %w", err)
	}
	log.Printf("file: %s size: %d bytes", fileName, bytes)

	return fmt.Sprintf("%x", hash.Sum(nil)), bytes, nil
}
