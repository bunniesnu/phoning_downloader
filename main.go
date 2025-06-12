package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/schollz/progressbar/v3"
)

func main() {
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
		arr := m[k].([]interface{})
		bar := progressbar.Default(int64(len(arr)))
		for _, item := range arr {
			obj := item.(map[string]interface{})
			id := obj["id"]
			isAudio := obj["a"]
			dir := "calls"
			if k == "p" {
				dir = "podcasts"
			}
			ext := "mp4"
			if isAudio.(bool) {
				ext = "m4a"
			}
			ffmpeg(
				fmt.Sprintf("https://cdn.newjeans.app/stream/%s/%d.m3u8", dir, int(id.(float64))),
				fmt.Sprintf("%s/%s%d.%s", outDir, k, int(id.(float64)), ext),
			)
			bar.Add(1)
		}
	}
}