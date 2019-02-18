package jt808

import "strings"

var (
	bcdx [16]byte = [16]byte{
		'0', '1', '2', '3',
		'4', '5', '6', '7',
		'8', '9', 'A', 'B',
		'C', 'D', 'E', 'F',
	}
)

func bcd2str(bs []byte) string {
	var buf strings.Builder

	for _, b := range bs {
		buf.WriteByte(bcdx[(b&0xf0)>>4])
		buf.WriteByte(bcdx[b&0x0f])
	}

	return buf.String()
}
