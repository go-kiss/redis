package protocol

import (
	"bytes"
	"testing"
)

func BenchmarkReaderParseReplyStatus(b *testing.B) {
	buf := new(bytes.Buffer)
	for i := 0; i < b.N; i++ {
		buf.WriteString("+OK\r\n")
	}
	p := NewReader(buf)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		if m, err := p.ReadStatusReply(); err != nil || m != "OK" {
			b.Fatal(err)
		}
	}
}

func BenchmarkReaderParseReplyInt(b *testing.B) {
	buf := new(bytes.Buffer)
	for i := 0; i < b.N; i++ {
		buf.WriteString(":1\r\n")
	}
	p := NewReader(buf)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		if i, err := p.ReadIntReply(); err != nil || i != 1 {
			b.Fatal(err)
		}
	}
}

func BenchmarkReaderParseReplyError(b *testing.B) {
	buf := new(bytes.Buffer)
	for i := 0; i < b.N; i++ {
		buf.WriteString("-Error message\r\n")
	}
	p := NewReader(buf)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		if _, err := p.ReadInt(); err.Error() != "Error message" {
			b.Fatal(err)
		}
	}
}

func BenchmarkReaderParseReplyString(b *testing.B) {
	buf := new(bytes.Buffer)
	for i := 0; i < b.N; i++ {
		buf.WriteString("$5\r\nhello\r\n")
	}
	p := NewReader(buf)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		if buf, err := p.ReadBytesReply(); err != nil || !bytes.Equal(buf, []byte("hello")) {
			b.Fatal(err)
		}
	}
}

func BenchmarkReaderParseReplyArray(b *testing.B) {
	buf := new(bytes.Buffer)
	for i := 0; i < b.N; i++ {
		buf.WriteString("*3\r\n$5\r\nhello\r\n$3\r\nfoo\r\n$3\r\nbar\r\n")
	}
	p := NewReader(buf)
	b.ResetTimer()

	bufs := [][]byte{[]byte("hello"), []byte("foo"), []byte("bar")}

	for i := 0; i < b.N; i++ {
		l, err := p.ReadArrayLenReply()
		if err != nil || l != 3 {
			b.Fatal(err)
		}

		for j := 0; j < l; j++ {
			if buf, err := p.ReadBytesReply(); err != nil || !bytes.Equal(buf, bufs[j]) {
				b.Fatal(err)
			}
		}

	}
}
