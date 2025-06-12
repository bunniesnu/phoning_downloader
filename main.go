package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"sync"

	"github.com/schollz/progressbar/v3"
)

func main() {
	// Parse concurrency flag
	conc := flag.String("c", "5", "Number of concurrent downloads")
	flag.Parse()
	var concurrency int
	if _, err := fmt.Sscanf(*conc, "%d", &concurrency); err != nil || concurrency < 1 {
		fmt.Fprintln(os.Stderr, "Invalid concurrency value. Please provide a positive integer with -c flag.")
		flag.Usage()
		os.Exit(1)
	}
	// Check if FFmpeg is installed
	ffmpegCheckCmd := exec.Command("ffmpeg", "-version")
	if err := ffmpegCheckCmd.Run(); err != nil {
		fmt.Println("Your system does not have FFmpeg installed or not in PATH. Refer: https://ffmpeg.org")
		os.Exit(1)
	}

	// Check if the data file exists
	dataFile := "data.json"
	if _, err := os.Stat(dataFile); err != nil {
		choice, err := promptChoice("You do not have a data file.\nPlease choose an option:",
			"Download all files",
			"Quit and generate a data file",
		)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error reading input:", err)
			os.Exit(1)
		}
		if choice == 2 {
			fmt.Println("Quitting...\nRefer to the documentation: https://github.com/bunniesnu/phoningdb_downloader#readme")
			os.Exit(0)
		}
		src := "docs/data.json"
		dst := "data.json"
		input, err := os.ReadFile(src)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to read %s: %v\n", src, err)
			os.Exit(1)
		}
		if err := os.WriteFile(dst, input, 0644); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to write %s: %v\n", dst, err)
			os.Exit(1)
		}
		fmt.Println("Copied docs/data.json to data.json.")
	}

	m, errMsg := validateJson(dataFile)
	if m == nil {
		fmt.Fprintln(os.Stderr, errMsg)
		os.Exit(1)
	}

	// Prompt user to select output directory
	fmt.Print("Enter output directory (default: output): ")
	var outDir string
	fmt.Scanln(&outDir)
	if outDir == "" {
		outDir = "output"
	}
	if err := os.MkdirAll(outDir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create output directory: %v\n", err)
		os.Exit(1)
	}

	// Download calls and podcasts
	for _, k := range []string{"c", "p"} {
		dir := "calls"
		if k == "p" {
			dir = "podcasts"
		}
		arr := m[k].([]interface{})
		fmt.Printf("Downloading %s (%d items)...\n", dir, len(arr))
		bar := progressbar.Default(int64(len(arr)))
		var wg sync.WaitGroup
		sem := make(chan struct{}, concurrency) // limit concurrent downloads
		wg.Add(len(arr))
		// Use context to handle cancellation on error
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		errCh := make(chan error, 1)

		for _, item := range arr {
			go func(item interface{}) {
				sem <- struct{}{} // acquire a slot
				defer func() { <-sem }() // release the slot
				defer wg.Done()
				select {
				case <-ctx.Done():
					return
				default:
				}
				obj := item.(map[string]interface{})
				id := obj["id"]
				isAudio := obj["a"]
				ext := "mp4"
				if isAudio.(bool) {
					ext = "m4a"
				}
				res, err := ffmpeg(
					fmt.Sprintf("https://cdn.newjeans.app/stream/%s/%d.m3u8", dir, int(id.(float64))),
					fmt.Sprintf("%s/%s%d.%s", outDir, k, int(id.(float64)), ext),
				)
				if !res || err != nil {
					select {
					case errCh <- fmt.Errorf("failed to download %s%d: %v", k, int(id.(float64)), err):
					default:
					}
					cancel()
					return
				}
				bar.Add(1)
			}(item)
		}

		done := make(chan struct{})
		go func() {
			wg.Wait()
			close(done)
		}()

		select {
		case err := <-errCh:
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		case <-done:
			fmt.Printf("Finished downloading %s.\n", dir)
		}
	}
}