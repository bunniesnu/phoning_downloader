package main

import (
	"fmt"
	"os"
	"os/exec"
)

func main() {
	// Check if FFmpeg is installed
	ffmpegCheckCmd := exec.Command("ffmpeg", "-version")
	if err := ffmpegCheckCmd.Run(); err != nil {
		fmt.Println("Your system does not have FFmpeg installed. Refer: https://ffmpeg.org")
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
}