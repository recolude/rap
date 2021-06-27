package binary

import "io"

//go:generate mockgen -destination=../../mocks/reader.go -package=mocks io Reader

// https://dave.cheney.net/2019/01/27/eliminate-error-handling-by-eliminating-errors
type ErrReader struct {
	io.Reader
	err error
	n   int
}

func NewErrReader(inStream io.Reader) *ErrReader {
	return &ErrReader{Reader: inStream}
}

func (e *ErrReader) Error() error {
	return e.err
}

func (e *ErrReader) Read(p []byte) (int, error) {
	if e.err != nil {
		return 0, e.err
	}

	n, err := io.ReadFull(e.Reader, p)
	e.err = err
	e.n += n
	return n, e.err
}

func (e *ErrReader) ReadByte() (byte, error) {
	b := []byte{0}
	e.Read(b)
	return b[0], e.err
}

func (e *ErrReader) TotalRead() int {
	return e.n
}
