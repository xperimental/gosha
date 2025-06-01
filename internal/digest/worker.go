package digest

import (
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"os"
	"sync"
	"time"
)

type Result struct {
	Filename string
	Hash     string
	Size     int64
	Duration time.Duration
	Error    error
}

func Worker(wg *sync.WaitGroup, resCh chan<- Result, fileCh <-chan string) {
	wg.Add(1)
	go func() {
		defer wg.Done()

		for fileName := range fileCh {
			start := time.Now()
			log.Printf("Processing file %s", fileName)
			hash, bytes, err := hashFile(fileName)
			duration := time.Since(start)

			resCh <- Result{
				Filename: fileName,
				Hash:     hash,
				Size:     bytes,
				Duration: duration,
				Error:    err,
			}
		}
	}()
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
