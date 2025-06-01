package digest

import (
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"fmt"
	"github.com/xperimental/gosha/internal/config"
	"hash"
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

func Worker(wg *sync.WaitGroup, resCh chan<- Result, fileCh <-chan string, algorithm config.Algorithm) {
	wg.Add(1)
	go func() {
		defer wg.Done()

		for fileName := range fileCh {
			start := time.Now()
			log.Printf("Processing file %s", fileName)
			hash, bytes, err := hashFile(fileName, algorithm)
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

func createDigest(algorithm config.Algorithm) hash.Hash {
	switch algorithm {
	case config.AlgorithmSha1:
		return sha1.New()
	case config.AlgorithmSha256:
		return sha256.New()
	case config.AlgorithmSha512:
		return sha512.New()
	}

	return nil
}

func hashFile(fileName string, algorithm config.Algorithm) (string, int64, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return "", 0, fmt.Errorf("can not open file: %w", err)
	}
	defer file.Close()

	hash := createDigest(algorithm)
	if hash == nil {
		return "", 0, fmt.Errorf("unknown algorithm: %s", algorithm)
	}

	bytes, err := io.Copy(hash, file)
	if err != nil {
		return "", 0, fmt.Errorf("error while hashing file: %w", err)
	}
	log.Printf("file: %s size: %d bytes", fileName, bytes)

	return fmt.Sprintf("%x", hash.Sum(nil)), bytes, nil
}
