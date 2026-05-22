package main

import (
	"fmt"
	"log"
	"os"

	"github.com/klauspost/compress/zip"
	"github.com/unxed/zipcharset"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("Usage: %s <archive.zip>", os.Args[0])
	}

	filePath := os.Args[1]
	f, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("failed to open archive: %v", err)
	}
	defer f.Close()

	stat, err := f.Stat()
	if err != nil {
		log.Fatalf("failed to stat file: %v", err)
	}

	// Create zip.ReaderOptions and inject our micro-library decoder
	opts := zip.ReaderOptions{
		NameDecoder: zipcharset.NewNameDecoder(),
	}

	zr, err := zip.NewReaderWithOptions(f, stat.Size(), opts)
	if err != nil {
		log.Fatalf("failed to read zip archive: %v", err)
	}

	fmt.Printf("Listing files in %s:\n", filePath)
	fmt.Println("----------------------------------------")
	for _, file := range zr.File {
		fmt.Printf("- %s (Size: %d bytes)\n", file.Name, file.UncompressedSize64)
	}
	fmt.Println("----------------------------------------")
}