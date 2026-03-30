package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"

	"dcaiovinicius/httpchecker/internal/config"
	"dcaiovinicius/httpchecker/internal/logger"
	"dcaiovinicius/httpchecker/internal/worker"
)

func main() {
	cfg, err := config.FromFlag()
	if err != nil {
		logger.Error("%v", err)
		os.Exit(1)
	}

	if cfg.Debug {
		logger.EnableDebug()
		logger.Info("debug mode enabled")
	}

	urls, err := loadURLs()
	if err != nil {
		logger.Error("%v", err)
		os.Exit(1)
	}

	w := worker.NewWorker(urls, cfg)
	w.Run()

	logger.Notice("finished checking %d URLs\n", len(urls))
}

func loadURLs() ([]string, error) {
	urls := make([]string, 0, len(flag.Args()))
	urls = append(urls, flag.Args()...)

	pipedURLs, err := readURLsFromStdin()
	if err != nil {
		return nil, err
	}
	urls = append(urls, pipedURLs...)

	if len(urls) == 0 {
		return nil, fmt.Errorf("no URLs provided (use args or pipe via stdin)")
	}
	return urls, nil
}

func readURLsFromStdin() ([]string, error) {
	stat, err := os.Stdin.Stat()
	if err != nil {
		return nil, fmt.Errorf("read stdin: %w", err)
	}
	if (stat.Mode() & os.ModeCharDevice) != 0 {
		return nil, nil
	}

	urls := []string{}
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		parts := strings.FieldsFunc(line, func(r rune) bool {
			return r == ',' || r == ' ' || r == '\t'
		})
		for _, p := range parts {
			u := strings.TrimSpace(p)
			if u != "" {
				urls = append(urls, u)
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scan stdin: %w", err)
	}
	return urls, nil
}
