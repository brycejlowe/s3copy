package locations

import (
	"io"
	"os"
)

type File struct {
	path string
	file *os.File
}

func NewFile(path string) *File {
	return &File{
		path: path,
	}
}

func (f *File) Name() string {
	return "file"
}

func (f *File) Path() string {
	return f.path
}

func (f *File) Read() (io.Reader, error) {
	// open the file to read from
	if fp, err := os.Open(f.Path()); err == nil {
		f.file = fp
		return fp, nil
	} else {
		return nil, err
	}
}

func (f *File) Done() {
	if f.file != nil {
		f.file.Close()
	}
}

func (f *File) Write(source Source) (int64, error) {
	defer source.Done()

	inputToWrite, err := source.Read()
	if err != nil {
		return -1, err
	}

	// open the file to write to
	fp, err := os.OpenFile(f.Path(), os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return -1, err
	}

	defer fp.Close()

	// write out the buffer to a local file
	if written, err := fp.ReadFrom(inputToWrite); err == nil {
		return written, nil
	} else {
		return -1, err
	}
}
