# Phoning Downloader
[![go report card](https://goreportcard.com/badge/github.com/bunniesnu/phoning_downloader)](https://goreportcard.com/report/github.com/bunniesnu/phoning_downloader)

A tool for downloading calls and podcasts of Phoning. Built with Go.

## Requirements
Your system needs FFmpeg installed and added to your PATH to run this program. See [FFmpeg Documentation](https://ffmpeg.org)

To check FFmpeg installation, run ```ffmpeg --version``` in your terminal.

No minimum version specified, but the higher the better.
## Install
Download the binary/executable from [Source]
## Usage
```
phoningdb_downloader [OPTIONS]
```
### Options
```
-c string
      Number of concurrent downloads (default "5")
-h    Show help
-o string
      Output directory
-verbose
      Enable verbose logging
```
## License
MIT