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
	// Parse flags
	conc := flag.String("c", "5", "Number of concurrent downloads")
	output_in := flag.String("o", "", "Output directory")
	verbose := flag.Bool("verbose", false, "Enable verbose logging")
	help := flag.Bool("h", false, "Show help")
	flag.Parse()
	if *help {
		flag.Usage()
		os.Exit(0)
	}
	var concurrency int
	if _, err := fmt.Sscanf(*conc, "%d", &concurrency); err != nil || concurrency < 1 {
		fmt.Fprintln(os.Stderr, "Invalid concurrency value. Please provide a positive integer with -c flag.")
		flag.Usage()
		os.Exit(1)
	}
	if *verbose {
		fmt.Printf("Concurrency value: %d\n", concurrency)
	}
	if concurrency > 20 {
		fmt.Println("Warning: If you set the concurrency value too high, your system may crash or become unresponsive.")
	}

	// Check if FFmpeg is installed
	if *verbose {
		fmt.Println("Checking FFmpeg installation...")
	}
	ffmpegCheckCmd := exec.Command("ffmpeg", "-version")
	if err := ffmpegCheckCmd.Run(); err != nil {
		fmt.Println("Your system does not have FFmpeg installed or not in PATH. Refer: https://ffmpeg.org")
		os.Exit(1)
	}
	if *verbose {
		fmt.Println("FFmpeg verified.")
	}

	// Check if the data file exists
	if *verbose {
		fmt.Println("Checking data file existance...")
	}
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
		if *verbose {
			fmt.Println("Copied docs/data.json to data.json")
		}
	} else if *verbose {
		fmt.Println("Data file exists.")
	}

	if *verbose {
		fmt.Println("Validating data file...")
	}
	m, errMsg := validateJson(dataFile)
	if m == nil {
		fmt.Fprintln(os.Stderr, errMsg)
		os.Exit(1)
	}
	if *verbose {
		fmt.Println("Data file is valid.")
	}

	// Prompt user to select output directory
	var outDir string
	if *output_in == "" {
		fmt.Print("Enter output directory (default: output): ")
		fmt.Scanln(&outDir)
		if outDir == "" {
			outDir = "output"
		}
	} else {
		outDir = *output_in
	}
	if err := os.MkdirAll(outDir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create output directory: %v\n", err)
		os.Exit(1)
	}
	if *verbose {
		fmt.Printf("Output directory: %s\n", outDir)
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
				inputUrl := fmt.Sprintf("https://cdn.newjeans.app/stream/%s/%d.m3u8", dir, int(id.(float64)))
				outputFile := fmt.Sprintf("%s/%s%d.%s", outDir, k, int(id.(float64)), ext)
				res, err := ffmpeg(
					inputUrl,
					outputFile,
				)
				if *verbose {
					bar.Clear()
					fmt.Println("Downloaded:", inputUrl, "to", outputFile)
				}
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