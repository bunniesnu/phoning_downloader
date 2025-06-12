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

	choice, err := promptChoice("You do not have a data file.\nPlease choose an option:",
		"Download all files",
		"Quit and generate a data file",
	)

    if err != nil {
        fmt.Fprintln(os.Stderr, "Error reading input:", err)
        os.Exit(1)
    }

	switch choice {
	case 1:
		fmt.Println("Downloading all files...")
	case 2:
		fmt.Println("Quitting...\nRefer to the documentation: https://github.com/bunniesnu/phoningdb_downloader#readme")
		os.Exit(0)
	}
}