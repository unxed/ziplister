# ziplister

`ziplister` is a simple demo command-line utility that lists the contents of a ZIP archive, leveraging the `zipcharset` micro-library to resolve legacy character encodings.

## Requirements

* Go 1.22 or higher

## Getting Started

1. Clone or navigate to the directory.
2. Build the application:

```bash
go build -o ziplister main.go
```

3. Run the lister against a legacy ZIP archive:

```bash
./ziplister path/to/legacy_archive.zip
```

## How It Works

The program uses `github.com/klauspost/compress/zip` instead of the standard library `archive/zip`. It instantiates the reader using `NewReaderWithOptions`, passing a `zipcharset.NewNameDecoder()` callback. As the archive is parsed, legacy metadata is intercepted and transcoded into UTF-8 before filename validation occurs.
```