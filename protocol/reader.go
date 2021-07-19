package protocol

import (
	"bufio"
	"fmt"
	"io"
	"strconv"

	"github.com/go-kiss/redis/util"
)

const (
	ErrorReply  = '-'
	StatusReply = '+'
	IntReply    = ':'
	StringReply = '$'
	ArrayReply  = '*'
)

//------------------------------------------------------------------------------

const Nil = RedisError("redis: nil")

type RedisError string

func (e RedisError) Error() string { return string(e) }

//------------------------------------------------------------------------------

type MultiBulkParse func(*Reader, int64) (interface{}, error)

type Reader struct {
	rd   *bufio.Reader
	_buf []byte
}

func NewReader(rd io.Reader) *Reader {
	return &Reader{
		rd:   bufio.NewReader(rd),
		_buf: make([]byte, 64),
	}
}

func (r *Reader) Reset(rd io.Reader) {
	r.rd.Reset(rd)
}

func (r *Reader) readLine() ([]byte, error) {
	line, isPrefix, err := r.rd.ReadLine()
	if err != nil {
		return nil, err
	}
	if isPrefix {
		return nil, bufio.ErrBufferFull
	}
	if len(line) == 0 {
		return nil, fmt.Errorf("redis: reply is empty")
	}
	if isNilReply(line) {
		return nil, Nil
	}
	return line, nil
}

func (r *Reader) ReadInterfaceReply() (interface{}, error) {
	line, err := r.readLine()
	if err != nil {
		return nil, err
	}
	return r.parseInterface(line)
}

func (r *Reader) parseInterface(line []byte) (interface{}, error) {
	switch line[0] {
	case ErrorReply:
		return nil, ParseErrorReply(line)
	case IntReply:
		return util.ParseInt(line[1:], 10, 64)
	case StatusReply:
		return string(line[1:]), nil
	case StringReply:
		b, err := r.readBytes(line)
		if err != nil {
			return nil, err
		}
		return string(b), nil
	case ArrayReply:
		vs, err := r.readArray(line)
		return vs, err
	default:
		return nil, fmt.Errorf("redis: can't parse int reply: %.100q", line)
	}
}

func (r *Reader) readArray(line []byte) ([]interface{}, error) {
	if isNilReply(line) {
		return nil, Nil
	}

	arrayLen, err := parseArrayLen(line)
	if err != nil {
		return nil, err
	}

	vs := make([]interface{}, 0, arrayLen)

	for i := 0; i < arrayLen; i++ {
		line, err := r.readLine()
		if err != nil {
			return nil, err
		}

		v, err := r.parseInterface(line)
		if err != nil {
			return nil, err
		}

		vs = append(vs, v)
	}

	return vs, nil
}

func (r *Reader) ReadIntReply() (int64, error) {
	line, err := r.readLine()
	if err != nil {
		return 0, err
	}
	switch line[0] {
	case ErrorReply:
		return 0, ParseErrorReply(line)
	case IntReply:
		return util.ParseInt(line[1:], 10, 64)
	default:
		return 0, fmt.Errorf("redis: can't parse int reply: %.100q", line)
	}
}

func (r *Reader) ReadStatusReply() (string, error) {
	line, err := r.readLine()
	if err != nil {
		return "", err
	}
	switch line[0] {
	case ErrorReply:
		return "", ParseErrorReply(line)
	case StatusReply:
		return string(line[1:]), nil
	case StringReply:
		buf, err := r.ReadBytesReply()
		if err != nil {
			return "", err
		}

		return string(buf), nil
	default:
		return "", fmt.Errorf("redis: can't parse reply=%.100q reading string", line)
	}
}

func (r *Reader) ReadBytesReply() ([]byte, error) {
	line, err := r.readLine()
	if err != nil {
		return nil, err
	}
	switch line[0] {
	case ErrorReply:
		return nil, ParseErrorReply(line)
	case StringReply:
		return r.readBytes(line)
	default:
		return nil, fmt.Errorf("redis: can't parse string reply: %.100q", line)
	}
}

func (r *Reader) readBytes(line []byte) ([]byte, error) {
	if isNilReply(line) {
		return nil, Nil
	}

	replyLen, err := strconv.Atoi(string(line[1:]))
	if err != nil {
		return nil, err
	}

	b := make([]byte, replyLen+2)
	_, err = io.ReadFull(r.rd, b)
	if err != nil {
		return nil, err
	}

	return b[:replyLen], nil
}

func (r *Reader) ReadArrayLenReply() (int, error) {
	line, err := r.readLine()
	if err != nil {
		return 0, err
	}
	switch line[0] {
	case ErrorReply:
		return 0, ParseErrorReply(line)
	case ArrayReply:
		return parseArrayLen(line)
	default:
		return 0, fmt.Errorf("redis: can't parse array reply: %.100q", line)
	}
}

func (r *Reader) ReadInt() (int64, error) {
	b, err := r.readTmpBytesReply()
	if err != nil {
		return 0, err
	}
	return util.ParseInt(b, 10, 64)
}

func (r *Reader) ReadUint() (uint64, error) {
	b, err := r.readTmpBytesReply()
	if err != nil {
		return 0, err
	}
	return util.ParseUint(b, 10, 64)
}

func (r *Reader) ReadFloat() (float64, error) {
	b, err := r.readTmpBytesReply()
	if err != nil {
		return 0, err
	}
	return util.ParseFloat(b, 64)
}

func (r *Reader) readTmpBytesReply() ([]byte, error) {
	line, err := r.readLine()
	if err != nil {
		return nil, err
	}
	switch line[0] {
	case ErrorReply:
		return nil, ParseErrorReply(line)
	case StringReply:
		return r._readTmpBytesReply(line)
	case StatusReply:
		return line[1:], nil
	default:
		return nil, fmt.Errorf("redis: can't parse string reply: %.100q", line)
	}
}

func (r *Reader) _readTmpBytesReply(line []byte) ([]byte, error) {
	if isNilReply(line) {
		return nil, Nil
	}

	replyLen, err := strconv.Atoi(string(line[1:]))
	if err != nil {
		return nil, err
	}

	buf := r.buf(replyLen + 2)
	_, err = io.ReadFull(r.rd, buf)
	if err != nil {
		return nil, err
	}

	return buf[:replyLen], nil
}

func (r *Reader) buf(n int) []byte {
	if d := n - cap(r._buf); d > 0 {
		r._buf = append(r._buf, make([]byte, d)...)
	}
	return r._buf[:n]
}

func isNilReply(b []byte) bool {
	return len(b) == 3 &&
		(b[0] == StringReply || b[0] == ArrayReply) &&
		b[1] == '-' && b[2] == '1'
}

func ParseErrorReply(line []byte) error {
	return RedisError(string(line[1:]))
}

func parseArrayLen(line []byte) (int, error) {
	if isNilReply(line) {
		return 0, Nil
	}
	return util.Atoi(line[1:])
}
