// Package hexd formats bytes into hex. It's really beautiful and makes me cry.
///////////////////////////////////////////////////////////////////////////////
// Copyright 2019, Joshua J Baker
//
// Permission to use, copy, modify, and/or distribute this software for any
// purpose with or without fee is hereby granted, provided that the above
// copyright notice and this permission notice appear in all copies.
//
// THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES
// WITH REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF
// MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR ANY
// SPECIAL, DIRECT, INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES
// WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR PROFITS, WHETHER IN AN
// ACTION OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION, ARISING OUT OF OR
// IN CONNECTION WITH THE USE OR PERFORMANCE OF THIS SOFTWARE.
///////////////////////////////////////////////////////////////////////////////
// Summer laid her simple Hat
// On it's boundless Shelf --
// Unobserved -- a Ribin slipt,
// Snatch it for yourself.
//
// Summer laid her supple Glove
// In it's sylvan Drawer --
// Wheresoe'er, or was she --
// The demand of Awe?
//
// - Emily Dickinson
///////////////////////////////////////////////////////////////////////////////
package hexd

import (
	"bufio"
	"bytes"
	"io"
	"strconv"
	"strings"
	"unicode/utf8"
)

// IncludeHeader will include the following header line to the top of every
// dump.
//     Offset (h)  00 01 02 03 04 05 06 07 08 09 0A 0B 0C 0D 0E 0F
var IncludeHeader = false

const hch = "0123456789ABCDEF"
const head = "Offset (h)  00 01 02 03 04 05 06 07 08 09 0A 0B 0C 0D 0E 0F\n"

// WriteTo writes bytes in b as pretty hex output to writer wr.
func WriteTo(wr io.Writer, b []byte) (n int, err error) {
	w := bufio.NewWriter(wr)
	if IncludeHeader {
		nn, err := w.Write([]byte(head))
		n += nn
		if err != nil {
			return n, err
		}
	}
	var line []byte
	for i := 0; i < len(b); i += 16 {
		p := b[i:]
		line = append(line[:0], "0000000000"...)
		line = strconv.AppendInt(line, int64(i), 16)
		mark := len(line)
		line = append(line, ' ', ' ')
		for i := 0; i < 16; i++ {
			if i < len(p) {
				line = append(line, hch[p[i]>>4], hch[p[i]&15], ' ')
			} else {
				line = append(line, ' ', ' ', ' ')
			}
		}
		line = append(line, ' ')
		for i := 0; i < 16; i++ {
			if i < len(p) {
				if p[i] < ' ' {
					line = append(line, '.')
				} else if p[i] > 127 {
					line = append(line, 0, 0, 0, 0, 0, 0)
					n := utf8.EncodeRune(line[len(line)-6:], rune(p[i]))
					line = line[:len(line)-6+n]
				} else {
					line = append(line, p[i])
				}
			} else {
				line = append(line, ' ')
			}
		}

		line = append(line, '\n')
		nn, err := w.Write(line[mark-10:])
		n += nn
		if err != nil {
			return n, err
		}
	}
	return n, w.Flush()
}

// Dump returns bytes in b as pretty hex output.
func Dump(b []byte) []byte {
	var wr bytes.Buffer
	WriteTo(&wr, b)
	return wr.Bytes()
}

// DumpString returns bytes in b as pretty hex output.
func DumpString(b []byte) string {
	var wr strings.Builder
	WriteTo(&wr, b)
	return wr.String()
}
