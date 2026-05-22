package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/klauspost/compress/zip"
)

func TestListArchive_LegacyCP866(t *testing.T) {
	tmp := t.TempDir()
	zipPath := filepath.Join(tmp, "test_legacy.zip")

	// 1. Create a dummy legacy ZIP file
	f, err := os.Create(zipPath)
	if err != nil {
		t.Fatalf("failed to create temp zip file: %v", err)
	}

	zw := zip.NewWriter(f)

	// "Привет.txt" encoded in CP866 (legacy Cyrillic DOS OEM encoding)
	cp866Name := string([]byte{0x8f, 0xe0, 0xa8, 0xa2, 0xa5, 0xe2, 0x2e, 0x74, 0x78, 0x74})

	fh := &zip.FileHeader{
		Name:           cp866Name,
		Method:         zip.Store,
		CreatorVersion: 0, // OS: MS-DOS/FAT to trigger OEMDecoder selection
	}

	w, err := zw.CreateHeader(fh)
	if err != nil {
		t.Fatalf("failed to create header: %v", err)
	}
	w.Write([]byte("legacy content"))

	if err := zw.Close(); err != nil {
		t.Fatalf("failed to close writer: %v", err)
	}
	f.Close()

	// 2. Run the lister logic on the created archive
	var buf bytes.Buffer
	if err := listArchive(zipPath, &buf); err != nil {
		t.Fatalf("listArchive failed: %v", err)
	}

	// 3. Verify that the output was properly decoded from CP866 to UTF-8
	output := buf.String()
	if !strings.Contains(output, "Привет.txt") {
		t.Errorf("expected output to contain decoded filename 'Привет.txt', got:\n%s", output)
	}
}

func TestListArchive_MissingFile(t *testing.T) {
	var buf bytes.Buffer
	err := listArchive("non_existent_file.zip", &buf)
	if err == nil {
		t.Error("expected error for non-existent file, got nil")
	}
}