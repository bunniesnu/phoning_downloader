package main

import (
	"os/exec"
)

func ffmpeg(inputUrl string, outputFile string) (bool, error) {
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
		return false, err
	}
	return true, nil
}