package main

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/klauspost/compress/zip"
	"github.com/unxed/zipcharset"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("Usage: %s <archive.zip>", os.Args[0])
	}

	if err := listArchive(os.Args[1], os.Stdout); err != nil {
		log.Fatalf("Error: %v", err)
	}
}

// listArchive parses the ZIP archive, decodes legacy names, and writes output to the provided writer.
func listArchive(filePath string, out io.Writer) error {
	f, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open archive: %w", err)
	}
	defer f.Close()

	stat, err := f.Stat()
	if err != nil {
		return fmt.Errorf("failed to stat file: %w", err)
	}

	// Inject the custom NameDecoder callback
	opts := zip.ReaderOptions{
		NameDecoder: zipcharset.NewNameDecoder(),
	}

	zr, err := zip.NewReaderWithOptions(f, stat.Size(), opts)
	if err != nil {
		return fmt.Errorf("failed to read zip archive: %w", err)
	}

	fmt.Fprintf(out, "Listing files in %s:\n", filePath)
	fmt.Fprintln(out, "----------------------------------------")
	for _, file := range zr.File {
		fmt.Fprintf(out, "- %s (Size: %d bytes)\n", file.Name, file.UncompressedSize64)
	}
	fmt.Fprintln(out, "----------------------------------------")
	return nil
}