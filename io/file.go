package io

import (
	"bufio"
	"encoding/binary"
	"os"
)

type FileType byte

const (
	MagicNumber = "ed"

	PostingListV1 FileType = 0
	StringIndexV1 FileType = 1
)

type FileWriter struct {
	file   *os.File
	writer *bufio.Writer
	size   int
}

func newFileWriter(name string) (*FileWriter, error) {
	f, err := os.Create(name)
	if err != nil {
		return nil, err
	}
	fw := &FileWriter{
		file:   f,
		writer: bufio.NewWriter(f),
		size:   0,
	}
	return fw, nil
}

func (fw *FileWriter) Close() error {
	err := fw.writer.Flush()
	if err != nil {
		return err
	}
	err = fw.file.Sync()
	if err != nil {
		return err
	}
	err = fw.file.Close()
	if err != nil {
		return err
	}
	return nil
}

func (fw *FileWriter) writeHeader(fileType FileType) error {
	fw.writer.WriteString(MagicNumber)
	fw.writer.WriteByte(byte(fileType))
}

func (fw *FileWriter) writeBytes(bytes []bytes) error {
	binary.PutUvarint(buf, x)
	fw.writer.Write(bytes)
}
