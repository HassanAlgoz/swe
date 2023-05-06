package s3

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/url"
	"path"
	"sync"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
)

// mockedS3Client is a Thread-Safe AWS S3 Mock for Golang,
// it only works with: GetObject, PutObject, DeleteObject, HeadBucket
// source: https://gist.github.com/kamermans/6d84d3a35f5809e67c71f657cbb16e02
// which goes accroding to
// the official documentation: https://aws.amazon.com/blogs/developer/mocking-out-then-aws-sdk-for-go-for-unit-testing/
type mockedS3Client struct {
	// By embedding the interface into the mock struct, the struct can be used as
	// a drop-in replacement for the real s3iface.S3API interface, even though it
	// doesn't provide any actual implementation for the methods.
	s3iface.S3API

	mu    sync.Mutex
	files map[string][]byte
	tags  map[string]map[string]string
}

func newMockedS3Client() *mockedS3Client {
	return &mockedS3Client{
		files: map[string][]byte{},
		tags:  map[string]map[string]string{},
	}
}

func (m *mockedS3Client) HeadBucket(in *s3.HeadBucketInput) (*s3.HeadBucketOutput, error) {
	return &s3.HeadBucketOutput{}, nil
}

func (m *mockedS3Client) PutObject(in *s3.PutObjectInput) (out *s3.PutObjectOutput, err error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	key := path.Join(*in.Bucket, *in.Key)
	m.files[key], err = ioutil.ReadAll(in.Body)

	m.tags[key] = map[string]string{}
	if in.Tagging != nil {
		u, err := url.Parse("/?" + *in.Tagging)
		if err != nil {
			panic(fmt.Errorf("Unable to parse AWS S3 Tagging string %q: %w", *in.Tagging, err))
		}

		q := u.Query()
		for k := range q {
			m.tags[key][k] = q.Get(k)
		}
	}

	return &s3.PutObjectOutput{}, nil
}

func (m *mockedS3Client) GetObject(in *s3.GetObjectInput) (*s3.GetObjectOutput, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	key := path.Join(*in.Bucket, *in.Key)
	if _, ok := m.files[key]; !ok {
		return nil, errors.New("Key does not exist")
	}

	return &s3.GetObjectOutput{
		Body: ioutil.NopCloser(bytes.NewReader(m.files[key])),
	}, nil
}

func (m *mockedS3Client) DeleteObject(in *s3.DeleteObjectInput) (*s3.DeleteObjectOutput, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.files[*in.Key]; ok {
		delete(m.files, *in.Key)
	}
	if _, ok := m.tags[*in.Key]; ok {
		delete(m.tags, *in.Key)
	}
	return &s3.DeleteObjectOutput{}, nil
}
