package locations

import (
	"bufio"
	"io"
	"os"
)

type Output struct {
	name   string
	path   string
	output io.Writer
}

func NewStdout() *Output {
	return &Output{
		name:   "stdout",
		path:   "-",
		output: os.Stdout,
	}
}

func (s *Output) Name() string {
	return s.name
}

func (s *Output) Path() string {
	return s.path
}

func (s *Output) Write(source Source) (int64, error) {
	defer source.Done()

	inputToWrite, err := source.Read()
	if err != nil {
		return -1, err
	}

	// TODO: double buffering
	output := bufio.NewWriter(s.output)
	defer output.Flush()

	return output.ReadFrom(inputToWrite)
}
