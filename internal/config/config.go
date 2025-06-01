package config

import (
	"fmt"
	"github.com/spf13/pflag"
	"os"
	"path/filepath"
	"runtime"
)

type Algorithm string

func (a *Algorithm) String() string {
	return string(*a)
}

func (a *Algorithm) Set(val string) error {
	switch val {
	case "1":
		*a = AlgorithmSha1
	case "256":
		*a = AlgorithmSha256
	case "512":
		*a = AlgorithmSha512
	default:
		return fmt.Errorf("hash algorithm unsupported: %s", val)
	}

	return nil
}

func (a *Algorithm) Type() string {
	return "algorithm"
}

const (
	AlgorithmSha1   Algorithm = "1"
	AlgorithmSha256 Algorithm = "256"
	AlgorithmSha384 Algorithm = "384"
	AlgorithmSha512 Algorithm = "512"
)

type Config struct {
	Workers   int
	FileNames []string
	Algorithm Algorithm
}

func Parse(cmd string, args []string) (*Config, error) {
	cfg := defaultConfig(cmd)

	flags := pflag.NewFlagSet(cmd, pflag.ContinueOnError)
	flags.IntVarP(&cfg.Workers, "workers", "w", cfg.Workers, "number of worker threads")
	flags.VarP(&cfg.Algorithm, "algorithm", "a", "hashing algorithm")

	if err := flags.Parse(os.Args[1:]); err != nil {
		return nil, fmt.Errorf("error parsing flags: %w", err)
	}

	fileNames := flags.Args()
	if len(args) < 1 {
		fileNames = []string{"-"}
	}
	cfg.FileNames = fileNames

	return cfg, nil
}

func defaultAlgorithm(cmd string) Algorithm {
	switch filepath.Base(cmd) {
	case "shasum", "sha1sum":
		return AlgorithmSha1
	case "sha256":
		return AlgorithmSha256
	case "sha384":
		return AlgorithmSha384
	case "sha512":
		return AlgorithmSha512
	default:
		return AlgorithmSha256
	}
}

func defaultConfig(cmd string) *Config {
	return &Config{
		Workers:   runtime.NumCPU(),
		Algorithm: defaultAlgorithm(cmd),
	}
}
