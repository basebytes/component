package oss

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

func NewOss(bucket *oss.Bucket, opts ...Option) (oss *Oss) {
	oss = &Oss{bucket: bucket, endPoint: defaultEndPoint, scheme: defaultScheme, expired: defaultExpired.Milliseconds() / 1000}
	for _, opt := range opts {
		opt(oss)
	}
	return
}

type Oss struct {
	endPoint string
	scheme   string
	domain   string
	expired  int64
	bucket   *oss.Bucket
}

func (s *Oss) Upload(filename string, reader io.Reader) error {
	return s.bucket.PutObject(filename, reader)
}

func (s *Oss) Download(filename string) ([]byte, error) {
	if content, err := s.bucket.GetObject(filename); err == nil {
		defer content.Close()
		return io.ReadAll(content)
	} else {
		return nil, err
	}
}

func (s *Oss) Delete(filename ...string) (err error) {
	_, err = s.bucket.DeleteObjects(filename)
	return
}

func (s *Oss) Exists(filename string) (bool, error) {
	return s.bucket.IsObjectExist(filename)
}

func (s *Oss) GetSignUrl(filename string) (url string, err error) {
	return s.bucket.SignURL(filename, http.MethodGet, s.expired)
}

func (s *Oss) GetInnerUrl(filename string) string {
	return fmt.Sprintf("%s://%s.%s/%s", s.scheme, s.bucket.BucketName, s.endPoint, filename)
}

func (s *Oss) GetDomainSignUrl(filename string) (url string, err error) {
	if s.domain == "" {
		err = errors.New("oss domain not set")
	} else if url, err = s.GetSignUrl(filename); err == nil {
		url = strings.Replace(url, fmt.Sprintf("%s.%s", s.bucket.BucketName, s.endPoint), s.domain, 1)
	}
	return
}

type Option func(s *Oss)

func Scheme(scheme string) Option {
	return func(s *Oss) {
		if scheme == SchemeHttp {
			s.scheme = SchemeHttp
		}
	}
}

func Endpoint(endPoint string) Option {
	return func(s *Oss) {
		if endPoint != "" {
			s.endPoint = endPoint
		}
	}
}

func ExpiredInSec(expired int64) Option {
	return func(s *Oss) {
		if expired > 0 {
			s.expired = expired
		}
	}
}

func Domain(domain string) Option {
	return func(s *Oss) {
		if domain != "" {
			s.domain = domain
		}
	}
}
