package config

import (
	"fmt"
	"github.com/spf13/pflag"
	"log"
	"os"
	"runtime"
)

type Config struct {
	Workers   int
	FileNames []string
}

func Parse(cmd string, args []string) (*Config, error) {
	cfg := defaultConfig()

	flags := pflag.NewFlagSet(cmd, pflag.ContinueOnError)
	flags.IntP("workers", "w", cfg.Workers, "number of worker threads")

	if err := flags.Parse(os.Args[1:]); err != nil {
		return nil, fmt.Errorf("error parsing flags: %w", err)
	}

	fileNames := flags.Args()
	if len(args) < 1 {
		log.Fatalf("Usage: %s file [file ...]", os.Args[0])
	}
	cfg.FileNames = fileNames

	return cfg, nil
}

func defaultConfig() *Config {
	return &Config{
		Workers: runtime.NumCPU(),
	}
}
