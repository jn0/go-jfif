package exif

import (
    "os"
    "io"
    "fmt"
)

type ExifData map[string]interface{}

type size_t uint64

type Byte uint8
type Word uint16
type Long uint32
type String string
type Rational struct { Num, Den int }

func loadSize(fd *os.File, size size_t) ([]byte, error) {
    var buffer = make([]byte, size)
    if size != size_t(len(buffer)) {
        return nil, fmt.Errorf(
            "exif.Load.loadSize(%s, %v): cannot allocate buffer",
            fd.Name(), size)
    }

    read, e := fd.ReadAt(buffer, 0)
    if e != nil && e != io.EOF {
        return nil, e
    }
    if read != len(buffer) {
        return nil, fmt.Errorf(
            "exif.Load.loadSize(%s, %v): cannot read entier file: got %v",
            fd.Name(), size, read)
    }
    return buffer, nil
}

type Exif struct {
    Path string
    Entries []Entry
    NoDataLeft bool
}

func (x *Exif) SectionAt(pos int64) Entry {
    for _, ent := range x.Entries {
        if ent.Pos() <= pos && pos <= ent.Pos() + ent.Len() {
            return ent
        }
    }
    return nil
}

func (x *Exif) Load(path string) error {
    fmt.Printf("exif.Load(%#v)\n", path)
    x.Path = path

    fd, e := os.Open(path)
    if e != nil {
        return fmt.Errorf("exif.Load.Open(%s): %v", path, e)
    }
    defer fd.Close()

    size, e := fd.Seek(0, 2)
    if e != nil {
        return fmt.Errorf("exif.Load.SeekEnd(%s): %v", path, e)
    }

    _, e = fd.Seek(0, 0)
    if e != nil {
        return fmt.Errorf("exif.Load.SeekStart(%s): %v", path, e)
    }

    fmt.Printf("exif.Load(%q): %v bytes\n", x.Path, size)

    var tmp anEntry

    for {
        entry, err := tmp.ReadEntry(fd)
        if err != nil {
            return fmt.Errorf("exif.Load.entry(%s): %v", path, err)
        }
        fmt.Printf("%s\n", entry)
        x.Entries = append(x.Entries, entry)
        if entry.GetId() == EOI { break; }
    }
    here, e := Tell(fd)
    if e != nil {
        return fmt.Errorf("exif.Load(%q): %v", path, e)
    }
    x.NoDataLeft = here == size
    if x.NoDataLeft {
        fmt.Printf("File data exhausted.\n")
    } else {
        fmt.Printf("File data dangle: expecting %v, got %v, delta %v\n",
                   size, here, size - here)
    }

    /*
    buffer, e := loadSize(fd, size_t(size))
    if e != nil {
        return fmt.Errorf("exif.Load(%s).loadSize(%v): %v", path, size, e)
    }
    */
    return nil
}

func (x *Exif) SaveTo(path string) error {
    fmt.Printf("exif.SaveTo(%#v)\n", path)
    fd, e := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0644)
    if e != nil { return e; }
    defer fd.Close()
    for _, entry := range x.Entries {
        fmt.Printf("\tsaving %s\n", entry)
        entry.Write(fd)
    }
    return nil
}

func (x *Exif) Save() error {
    fmt.Printf("exif.Save() -> %#v\n", x.Path)
    // return x.SaveTo(x.Path)
    return nil
}

func (x *Exif) Inject(xif ExifData) error {
    fmt.Printf("exif.Inject(%q, %#v)\n", x.Path, xif)
    return nil
}
/* vim: set ft=go ai et ts=4 sts=4 sw=4: EOF */
