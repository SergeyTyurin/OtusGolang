package main

import (
	"errors"
	"io"
	"os"

	"github.com/cheggaaa/pb/v3"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
	ErrIsDirectory           = errors.New("path is directory, not a file")
	ErrIsSamePaths           = errors.New("from file and to file are same")
	ChunkBytes               = int64(1024)
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	fileInfo, err := os.Stat(fromPath)
	if err != nil {
		return ErrUnsupportedFile
	}

	if fileInfo.IsDir() {
		return ErrIsDirectory
	}

	if fromPath == toPath {
		return ErrIsSamePaths
	}

	if fileInfo.Size() < offset {
		return ErrOffsetExceedsFileSize
	}

	fromFile, err := os.Open(fromPath)
	if err != nil {
		return ErrUnsupportedFile
	}

	toFile, err := os.Create(toPath)
	if err != nil {
		return ErrUnsupportedFile
	}

	defer fromFile.Close()
	defer toFile.Close()

	chunk := ChunkBytes
	total := fileInfo.Size() - offset
	if limit > 0 && total > limit {
		total = limit
		if limit < ChunkBytes {
			chunk = limit
		}
	}

	copiedBytes := int64(0)
	fromFile.Seek(offset, io.SeekStart)
	bar := pb.Full.Start64(total)
	barReader := bar.NewProxyReader(fromFile)
	for {
		copied, err := io.CopyN(toFile, barReader, chunk)
		copiedBytes += copied
		if errors.Is(err, io.EOF) || copiedBytes == total {
			break
		}
		if total-copiedBytes < ChunkBytes {
			chunk = total - copiedBytes
		}
	}
	bar.Finish()
	return nil
}
