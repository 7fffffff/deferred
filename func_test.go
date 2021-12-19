package deferred_test

import (
	"errors"
	"io"
	"testing"

	"github.com/7fffffff/deferred"
)

var (
	errWriteFail = errors.New("bad write")
	errCloseFail = errors.New("bad close")
)

type writeTest struct {
	writeErr error
	closeErr error
}

var writeTests = []writeTest{
	{
		writeErr: nil,
		closeErr: nil,
	},
	{
		writeErr: errWriteFail,
		closeErr: nil,
	},
	{
		writeErr: nil,
		closeErr: errCloseFail,
	},
	{
		writeErr: errWriteFail,
		closeErr: errCloseFail,
	},
}

func TestFuncErr(t *testing.T) {
	blobs := [][]byte{
		[]byte("example"),
	}
	for i, test := range writeTests {
		w := &discarder{
			writeErr: test.writeErr,
			closeErr: test.closeErr,
		}
		err := writeAndClose(w, blobs)
		if test.writeErr != nil {
			if err != test.writeErr {
				t.Errorf("write %d: expected %v, got %v", i, test.writeErr, err)
			}
			continue
		}
		if test.closeErr != nil {
			if err != test.closeErr {
				t.Errorf("close %d: expected %v, got %v", i, test.closeErr, err)
			}
			continue
		}
		if err != nil {
			t.Errorf("%d: %v", i, err)
		}
	}
}

type discarder struct {
	writeErr error
	closeErr error
}

func (d *discarder) Close() error {
	return d.closeErr
}

func (d *discarder) Write(b []byte) (int, error) {
	return len(b), d.writeErr
}

func writeAndClose(output io.WriteCloser, records [][]byte) (err error) {
	defer deferred.FuncErr(output.Close)(&err)
	for _, record := range records {
		_, err = output.Write(record)
		if err != nil {
			return err
		}
	}
	return nil
}
