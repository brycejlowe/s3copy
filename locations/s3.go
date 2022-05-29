package locations

import (
	"context"
	"fmt"
	"io"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go/aws"
)

type S3 struct {
	path           string
	readerTempFile *os.File
}

func NewS3(path string) *S3 {
	return &S3{
		// strip off the s3 protocol, we already know we're dealing with s3
		path: strings.TrimLeft(path, "s3://"),
	}
}

func (s *S3) Name() string {
	return "S3"
}

func (s *S3) GetBucket() string {
	// return the bucket (the string up to the first /)
	return strings.Split(s.path, "/")[0]
}

func (s *S3) GetPath() string {
	// return the path (the remainder after the first /)
	parts := strings.SplitN(s.Path(), "/", 2)
	return parts[len(parts)-1]
}

func (s *S3) Path() string {
	return s.path
}

func (s *S3) s3Client() (*s3.Client, error) {
	if cfg, err := config.LoadDefaultConfig(context.TODO()); err == nil {
		return s3.NewFromConfig(cfg), nil
	} else {
		return nil, err
	}
}

func (s *S3) Read() (io.Reader, error) {
	s3Client, err := s.s3Client()
	if err != nil {
		return nil, err
	}

	s3Downloader := manager.NewDownloader(s3Client, func(d *manager.Downloader) {
		d.PartSize = 2 * 1024 * 1024 * 1024 // 2GB
	})

	// we're going to shove a temporary file in the current working directory
	currentDirectory, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	tempFilePrefix := ".tmp." + path.Base(os.Args[0]) + "-"
	dirEntries, err := os.ReadDir(currentDirectory)
	if err != nil {
		return nil, err
	}

	foundCount := 0
	for _, dirEntry := range dirEntries {
		if strings.HasPrefix(dirEntry.Name(), tempFilePrefix) {
			foundCount++
		}
	}

	tempFilePath := path.Join(currentDirectory, fmt.Sprint(tempFilePrefix+strconv.Itoa(foundCount)))

	/*
		TODO: FIX THIS
		The AWS SDK V2 for Go implements concurrent downloads from S3 for improved performance.  The object that
		receives the downloaded byte stream needs to be a WriteAtBuffer.  One of the things the WriteAtBuffer does
		is implement the WriteAt method, which can take a chunk of the downloaded byte stream and write it to a
		particular location.  Because of the way I have implemented the hand-off between a source and destination (via
		a io.Reader, which might not be the best method, but I am still learning) the byte stream at least needs to be
		sequential to support all of the output types (namely stdout).

		This means that any download from s3 is buffered to a file prior to being sent on to its destination, which
		means it can chew up double dis space.  For now this is OK as the primary use-case is uploading to S3.
	*/
	wfp, err := os.OpenFile(tempFilePath, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return nil, err
	}

	defer wfp.Close()

	_, err = s3Downloader.Download(context.TODO(), wfp, &s3.GetObjectInput{
		Bucket: aws.String(s.GetBucket()),
		Key:    aws.String(s.GetPath()),
	})
	if err != nil {
		return nil, err
	}

	if fp, err := os.Open(tempFilePath); err == nil {
		s.readerTempFile = fp
		return fp, nil
	} else {
		return nil, err
	}
}

func (s *S3) Done() {
	if s.readerTempFile != nil {
		s.readerTempFile.Close()
		os.Remove(s.readerTempFile.Name())
	}
}

func (s *S3) Write(source Source) (int64, error) {
	defer source.Done()

	inputToWrite, err := source.Read()
	if err != nil {
		return -1, err
	}

	s3Client, err := s.s3Client()
	if err != nil {
		return -1, err
	}

	// get an s3 client
	s3Uploader := manager.NewUploader(s3Client, func(u *manager.Uploader) {
		u.PartSize = 1 * 1024 * 1024 * 1024 // 1GB
	})

	_, err = s3Uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(s.GetBucket()),
		Key:    aws.String(s.GetPath()),
		Body:   inputToWrite,
	})
	if err != nil {
		return -1, err
	}

	// TODO: return the number of bytes read
	return 0, nil
}
