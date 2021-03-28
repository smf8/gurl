package request

import (
	"bytes"
	"io"
	"os"
)

type Payload interface {
	//Data returns payload data in []byte. nil indicates error in payload
	Data() io.Reader
}

type FilePayload struct {
	FilePath string
	file     *os.File
}

type RawPayload struct {
	Content []byte
}

func (f *FilePayload) Data() io.Reader {
	file, err := os.Open(f.FilePath)
	if err != nil {
		return nil
	}

	f.file = file

	return file
}

func (f *FilePayload) CloseFile() error {
	return f.file.Close()
}

func (r *RawPayload) Data() io.Reader {
	return bytes.NewReader(r.Content)
}
