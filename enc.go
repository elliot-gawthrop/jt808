package jt808

import (
	"bufio"
	"bytes"
	"encoding"
	"encoding/binary"
	"fmt"
	"io"
)

func JT808Split(data []byte, atEOF bool) (int, []byte, error) {
	var start, end int = -1, -1

	if start = bytes.IndexByte(data, 0x7e); start != -1 {
		if end = bytes.IndexByte(data[start+1:], 0x7e); end != -1 {
			return start + end + 2, data[start : start+end+2], nil
		}
	}

	if atEOF {
		return len(data), nil, io.EOF
	}

	return start, nil, nil
}

func NewScanner(r io.Reader) *bufio.Scanner {
	scanner := bufio.NewScanner(r)
	scanner.Split(JT808Split)
	return scanner
}

type BufDecoder struct {
	bs []byte
	l  int
	s  int
}

func NewBufDecoder(bs []byte) *BufDecoder {
	return &BufDecoder{bs: bs, l: len(bs), s: 0}
}

func (d *BufDecoder) Read(bs []byte) (int, error) {
	l := len(bs)
	if err := d.check(l); err != nil {
		return 0, err
	}
	i := copy(bs, d.bs[d.s:d.s+l])
	d.s += l
	return i, nil
}

func (d *BufDecoder) ReadByte() (byte, error) {
	if err := d.check(1); err != nil {
		return 0, err
	}

	b := d.bs[d.s]
	d.s += 1

	return b, nil
}

func (d *BufDecoder) Uint16() (uint16, error) {
	if err := d.check(2); err != nil {
		return 0, err
	}
	i := binary.BigEndian.Uint16(d.bs[d.s : d.s+2])
	d.s += 2
	return i, nil
}

func (d *BufDecoder) Uint32() (uint32, error) {
	if err := d.check(4); err != nil {
		return 0, err
	}

	i := binary.BigEndian.Uint32(d.bs[d.s : d.s+4])
	d.s += 4

	return i, nil
}

func (d *BufDecoder) Uint64() (uint64, error) {
	if err := d.check(8); err != nil {
		return 0, err
	}

	i := binary.BigEndian.Uint64(d.bs[d.s : d.s+8])
	d.s += 8

	return i, nil
}

func (d *BufDecoder) BCD(n int) (string, error) {
	if err := d.check(n); err != nil {
		return "", err
	}

	s := bcd2str(d.bs[d.s : d.s+n])
	d.s += n

	return s, nil
}

func (d *BufDecoder) check(size int) error {
	if d.s+size > d.l {
		return io.ErrShortBuffer
	}

	return nil
}

/*
type BufEncoder struct {
	*bytes.Buffer
}

func NewBufEncoder() *BufEncoder {
	return &BufEncoder{Buffer: new(bytes.Buffer)}
}

func (b *BufEncoder) PutUint16(i uint16) {
	bs := make([]byte, 2, 2)
	binary.BigEndian.PutUint16(bs, i)
	b.Write(bs)
}

func (b *BufEncoder) PutUint32(i uint32) {
	bs := make([]byte, 4, 4)
	binary.BigEndian.PutUint32(bs, i)
	b.Write(bs)
}

func (b *BufEncoder) PutUint64(i uint64) {
	bs := make([]byte, 8, 8)
	binary.BigEndian.PutUint64(bs, i)
	b.Write(bs)
}
*/

func checkSum(msg []byte, sum byte) ([]byte, error) {
	/*
		var check byte = 0x00
		for _, b := range msg {
			check ^= b
		}

		if check != sum {
			return nil, fmt.Errorf("check sum error")
		}
	*/
	return msg, nil
}

func Restore(bs []byte) ([]byte, error) {

	l := len(bs)

	if l == 0 || bs[0] != 0x7e || bs[l-1] != 0x7e {
		return nil, fmt.Errorf("parser error")
	}

	i := 1
	l -= 1

	var buf bytes.Buffer

	for i < l {
		if bs[i] == 0x7d {
			i++
			if bs[i] == 0x02 {
				buf.WriteByte(0x7e)
			} else if bs[i] == 0x01 {
				buf.WriteByte(0x7d)
			}
			i++
		} else {
			buf.WriteByte(bs[i])
			i++
		}
	}

	msg := buf.Bytes()
	l = len(msg)
	if l < 13 {
		return nil, io.ErrShortBuffer
	}

	return checkSum(msg[:l-1], msg[l-1])
}

func JT808SendMsg(writer io.Writer, msghead, msgbody encoding.BinaryMarshaler) error {
	w := bufio.NewWriterSize(writer, 1024)

	var head, body []byte
	var err error

	if head, err = msghead.MarshalBinary(); err != nil {
		return err
	}

	if body, err = msgbody.MarshalBinary(); err != nil {
		return err
	}

	var check byte = 0x00
	w.WriteByte(0x7e)

	for _, b := range head {
		check ^= b
		if b == 0x7e {
			w.WriteByte(0x7d)
			w.WriteByte(0x02)
		} else if b == 0x7d {
			w.WriteByte(0x7d)
			w.WriteByte(0x01)
		} else {
			w.WriteByte(b)
		}
	}

	for _, b := range body {
		check ^= b
		if b == 0x7e {
			w.WriteByte(0x7d)
			w.WriteByte(0x02)
		} else if b == 0x7d {
			w.WriteByte(0x7d)
			w.WriteByte(0x01)
		} else {
			w.WriteByte(b)
		}
	}

	w.WriteByte(check)
	w.WriteByte(0x7e)

	return w.Flush()
}
