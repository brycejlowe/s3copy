package locations

import (
	"io"
	"os"
)

type Input struct {
	name   string
	reader io.Reader
}

func NewStdin() *Input {
	return &Input{
		name:   "stdin",
		reader: os.Stdin,
	}
}

func (i *Input) Name() string {
	return i.name
}

func (i *Input) Done() {
	// no-op
}

func (i *Input) Read() (io.Reader, error) {
	return i.reader, nil
}
