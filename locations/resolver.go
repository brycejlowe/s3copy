package locations

import (
	"errors"
	"fmt"
	"io"
	"strings"
)

type Destination interface {
	Name() string
	Path() string
	Write(source Source) (int64, error)
}

type Source interface {
	Name() string
	Read() (io.Reader, error)
	Done()
}

func ResolveDestination(destination string) (Destination, error) {
	// hypen says we want to fetch from stdin
	if destination == "-" {
		return NewStdout(), nil
	}

	// look for the s3 prefix
	if isS3(destination) {
		return NewS3(destination), nil
	}

	// other things should be file locations
	if destination != "" {
		return NewFile(destination), nil
	}

	return nil, errors.New(fmt.Sprint("Unable to Resolve Destination:", destination))
}

func ResolveSource(source string) (Source, error) {
	if source == "-" {
		return NewStdin(), nil
	}

	// look for the s3 prefix
	if isS3(source) {
		return NewS3(source), nil
	}

	if source != "" {
		return NewFile(source), nil
	}

	return nil, errors.New(fmt.Sprint("Unable to Resolve Source:", source))
}

func isS3(input string) bool {
	return strings.HasPrefix(input, "s3://")
}
