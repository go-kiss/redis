package protocol

import (
	"bytes"
	"io/ioutil"
	"testing"
)

func TestWriteArgs(t *testing.T) {
	buf := new(bytes.Buffer)
	wr := NewWriter(buf)

	err := wr.WriteArgs([]interface{}{
		"string",
		12,
		34.56,
		[]byte{'b', 'y', 't', 'e', 's'},
		true,
		nil,
	})

	expected := []byte("*6\r\n" +
		"$6\r\nstring\r\n" +
		"$2\r\n12\r\n" +
		"$5\r\n34.56\r\n" +
		"$5\r\nbytes\r\n" +
		"$1\r\n1\r\n" +
		"$0\r\n" +
		"\r\n")

	if err != nil || bytes.Equal(buf.Bytes(), expected) {
		t.Fatal("WriteArgs faild")
	}
}

func BenchmarkWriteBuffer_Append(b *testing.B) {
	buf := NewWriter(ioutil.Discard)

	for i := 0; i < b.N; i++ {
		err := buf.WriteArgs("hello", "world", "foo", "bar")
		if err != nil {
			panic(err)
		}

		err = buf.Flush()
		if err != nil {
			panic(err)
		}
	}
}
