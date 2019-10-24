package exif

import (
    "fmt"
    "io"
    "encoding/binary"
)

type Seeker interface{
    Seek(int64, int) (int64, error)
}

// return current position in a file
func Tell(fd Seeker) (int64, error) {
    pos, e := fd.Seek(0, 1)
    if e != nil {
        return -1, e
    }
    return pos, nil
}

type SeekingReader interface {
    io.Reader
    Seeker
}

type Writer io.Writer

func GetWordBE(b []byte) Word {
    return Word(binary.BigEndian.Uint16(b))
}

func GetWordLE(b []byte) Word {
    return Word(binary.LittleEndian.Uint16(b))
}

func WriteWordBE(fd Writer, v Word) error {
    b := make([]byte, 2)
    binary.BigEndian.PutUint16(b, uint16(v))
    _, e := fd.Write(b)
    return e
}

func WriteWordLE(fd Writer, v Word) error {
    b := make([]byte, 2)
    binary.LittleEndian.PutUint16(b, uint16(v))
    _, e := fd.Write(b)
    return e
}

func WriteBytes(fd Writer, v []Byte) error {
    b := make([]byte, len(v))
    for i, x := range v { b[i] = byte(x); }
    _, e := fd.Write(b)
    return e
}

func WriteByte(fd Writer, v Byte) error {
    b := make([]byte, 1)
    b[0] = byte(v)
    _, e := fd.Write(b)
    return e
}

func ReadStructHere(fd SeekingReader, data interface{}) error {
    pos, e := Tell(fd)
    if e != nil {
        return e
    }
    if e = binary.Read(fd, binary.BigEndian, data); e != nil {
        return e
    }
    if _, e = fd.Seek(pos, 0); e != nil {
        return e
    }
    return nil
}

// Write()s dump (in `hexdump -C` format) of `data` to `out` with `title`
func Dump(out io.Writer, title string, data []byte) {
    output := fmt.Sprintf("******** %s (%d)\n", title, len(data))
    out.Write([]byte(output))
    for offset := 0; offset < len(data); offset += 16 {
        output = fmt.Sprintf("%08x", offset)
        ep := offset + 16
        if ep > len(data) { ep = len(data); }
        part := data[offset:ep]
        for _, c := range part {
            output += fmt.Sprintf(" %02x", c)
        }
        for len(output) < 8 + 3 * 16 {
            output += "   "
        }
        output += " |"
        for _, c := range part {
            /*
            s := fmt.Sprintf("%q", c)
            if s[1] == '\\' {
                if s[2] == '\\' {
                    output += "\\" // the backslash itself
                } else {
                    output += "." // non-printable stuff
                }
            } else {
                output += string(s[1]) // printable
            }
            */
            if c < ' ' || c > '~' { // is it a dirty hack? well, maybe.
                output += "." // non-printable stuff
            } else {
                output += string(c) // printable
            }
        }
        for len(output) < 8 + 3 * 16 + 2 + 16 {
            output += " "
        }
        output += "|\n"
        out.Write([]byte(output))
    }
    output = fmt.Sprintf("%08x = %d\n", len(data), len(data))
    out.Write([]byte(output))
}
/* vim: set ft=go ai et ts=4 sts=4 sw=4: EOF */
