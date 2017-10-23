package godgdat

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
	"os"
)

type DGDat struct {
	fp *os.File
}

func Open(file string) (*DGDat, error) {
	d := &DGDat{}

	if err := d.open(file); err != nil {
		return nil, err
	}

	if err := d.readHeader(); err != nil {
		return nil, err
	}

	return d, nil
}

func (dg *DGDat) Close() {
	dg.fp.Close()
}

func (dg *DGDat) open(path string) error {
	fp, err := os.Open(path)
	if err != nil {
		return err
	}

	dg.fp = fp

	return nil
}

func (dg *DGDat) readHeader() error {

	// Read magic number
	dg.readUint32()
	dg.readBytes(1)

	// Unknown read ops
	dg.readUint32()
	dg.readUint32()

	dg.readPackedValue()
	dg.readPackedValue()
	dg.readPackedValue()
	dg.readPackedValue()

	tbllen := int(dg.readBytes(1)[0])

	currentPosition, _ := dg.fp.Seek(0, os.SEEK_CUR)
	fmt.Println(currentPosition)
	fmt.Println(tbllen)

	for i := currentPosition; i < (currentPosition + int64(tbllen)); {
		len := int(dg.readBytes(1)[0])
		chunk := string(dg.readBytes(len))

		size := dg.readPackedValue()
		fmt.Printf("%s\tlen:0x%X\n", chunk, size)

		switch chunk {
		case "name", "cpt", "fbn", "lang", "stat":
			fmt.Println(chunk)
		default:
			temp := dg.readWideString(size)
			fmt.Printf("%s\t'%s'\n", chunk, temp)
		}

		// At end of loop
		i, _ = dg.fp.Seek(0, os.SEEK_CUR)
	}

	return nil
}

func (dg *DGDat) readBytes(n int) []byte {
	msg := make([]byte, n)
	dg.fp.Read(msg)
	return msg
}

func (dg *DGDat) readString() string {
	n := int64(dg.readBytes(1)[0])
	msg := make([]byte, n)
	dg.fp.Read(msg)
	return string(msg)
}

func (dg *DGDat) readUint16() uint16 {
	msg := make([]byte, 2)
	dg.fp.Read(msg)
	v := binary.BigEndian.Uint16(msg)
	return v
}

func (dg *DGDat) readUint32() uint32 {
	msg := make([]byte, 4)
	dg.fp.Read(msg)
	v := binary.BigEndian.Uint32(msg)
	return v
}

func (dg *DGDat) readUint64() uint64 {
	msg := make([]byte, 8)
	dg.fp.Read(msg)
	v := binary.BigEndian.Uint64(msg)
	return v
}

func (dg *DGDat) readPackedValue() int {
	size := int(dg.readBytes(1)[0])

	if size >= 0xF0 {
		s := dg.readBytes(4)
		v := (int(s[0]) << 24) | (int(s[1]) << 16) | (int(s[2]) << 8) | int(s[3])
		return v
	}

	if size >= 0xE0 {
		s := dg.readBytes(3)
		v := (size ^ 0xE0<<24) | (int(s[0]) << 16) | (int(s[1]) << 8) | int(s[2])
		return v
	}

	if size >= 0xC0 {
		s := dg.readBytes(2)
		v := (size ^ 0xC0<<16) | (int(s[0]) << 8) | int(s[1])
		return v
	}

	if size >= 0x80 {
		s := dg.readBytes(1)
		v := (size ^ 0x80<<8) | int(s[0])
		return v
	}

	return size
}

func (dg *DGDat) readWideString(size int) string {
	buf := dg.readBytes(size)

	// x2 := dg.readPackedValue()
	// x2 := dg.readPackedValue()
	x2 := int(buf[0])

	var sbuf bytes.Buffer
	zbuf := buf[1:x2]

	if len(buf) > x2+1 {
		buf = buf[x2+1:]

		mcount := int(buf[0])
		arr := buf[1:mcount]
		buf = buf[mcount+1:]

		iter := 0
		for len(buf) > 0 {
			v := int(buf[0])
			buf = buf[1:]
			count := int(math.Floor(float64(v) / float64(mcount)))
			offset := v % mcount

			if count == 0 {
				l := len(zbuf)
				for j := iter; j < l; j++ {
					sbuf.WriteByte(zbuf[j])
					sbuf.WriteByte(arr[offset])
				}
			} else {
				for j := iter; j < count; j++ {
					sbuf.WriteByte(zbuf[j])
					sbuf.WriteByte(arr[offset])
					iter += 2
				}
			}
		}
	} else {
		return string(zbuf)
	}

	return sbuf.String()
}
