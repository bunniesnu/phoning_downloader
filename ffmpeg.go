package main

import (
	"fmt"
	"os"
	"os/exec"
)

func ffmpeg(inputUrl string, outputFile string) {
	cmd := exec.Command(
		"ffmpeg",
		"-y",
		"-i", inputUrl,
		"-c", "copy",
		outputFile,
	)
	cmd.Stdout = nil
	cmd.Stderr = nil
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "ffmpeg error: %v\n", err)
		os.Exit(1)
	}
}