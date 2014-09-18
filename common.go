package common

import (
	"bufio"
	"bytes"
	"crypto/rand"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	TitaniumEndpoint = "https://titanium-dot-atomos-release.appspot.com/"
	OxygenEndpoint   = "https://oxygen-dot-atomos-release.appspot.com/"
)

type PrintReader struct {
	reader io.Reader
}

func (printReader *PrintReader) Read(p []byte) (n int, err error) {
	n, err = printReader.reader.Read(p)
	fmt.Printf("%s", p[:n])
	return n, err
}

func NewPrintReader(reader io.Reader) *PrintReader {
	return &PrintReader{
		reader: reader,
	}
}

type CountingReader struct {
	reader   io.Reader
	count    int
	stepping int
}

func (r *CountingReader) Read(p []byte) (n int, err error) {
	r.count += 1
	if r.count%r.stepping == 0 {
		fmt.Printf("Count:%d, Length: %d\n", r.count, len(p))
	}
	return r.reader.Read(p)
}

func NewCountingReader(reader io.Reader, stepping int) *CountingReader {
	return &CountingReader{
		reader:   reader,
		stepping: stepping,
	}
}

type RandomReader struct {
}

func (randomReader *RandomReader) Read(p []byte) (n int, err error) {
	return rand.Read(p)
}

func NewRandomReader() *RandomReader {
	return &RandomReader{}
}

type LimitRandomReader struct {
	reader io.Reader
}

func (limitRandomReader *LimitRandomReader) Read(p []byte) (n int, err error) {
	return limitRandomReader.reader.Read(p)
}

func NewLimitRandomReader(length int64) *LimitRandomReader {
	return &LimitRandomReader{
		reader: io.LimitReader(NewRandomReader(), length),
	}
}

func NewEmptyReader() io.Reader {
	return &EmptyReader{}
}

type EmptyReader struct{}

func (reader *EmptyReader) Read(p []byte) (n int, err error) {
	return 0, io.EOF
}

type BufferedReadCloser struct {
	bufReader *bufio.Reader
	reader    io.ReadCloser
}

func NewBufferedReadCloser(reader io.ReadCloser) *BufferedReadCloser {
	return &BufferedReadCloser{
		reader:    reader,
		bufReader: bufio.NewReader(reader),
	}
}

func (reader *BufferedReadCloser) Close() error {
	return reader.reader.Close()
}

func (reader *BufferedReadCloser) Read(p []byte) (n int, err error) {
	return reader.bufReader.Read(p)
}

func HandlerPatternRegistered(pattern string, serveMux *http.ServeMux) bool {
	r, err := http.NewRequest("GET", pattern, bytes.NewReader([]byte{}))
	if err != nil {
		log.Fatalf("%s", err)
	}

	_, p := serveMux.Handler(r)
	if p == "" {
		return false
	}

	return true
}

func WriteReadCloserToFile(name string, data io.ReadCloser) error {
	defer data.Close()

	buf, err := ioutil.ReadAll(data)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(name, buf, 0750)
	if err != nil {
		return err
	}

	return nil
}
