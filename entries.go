package exif

import (
    "os"
    "fmt"
    "encoding/binary"
)

const EntryStringFormat = "<%s:%#v>"

const (
    SOI  = 0xd8 // Start of Image

    APP0 = 0xe0 // JFIF application segment
    APP1 = 0xe1 // Other APP segments
    APP2 = 0xe2
    APP3 = 0xe3
    APP4 = 0xe4
    APP5 = 0xe5
    APP6 = 0xe6
    APP7 = 0xe7
    APP8 = 0xe8
    APP9 = 0xe9
    APPa = 0xea
    APPb = 0xeb
    APPc = 0xec
    APPd = 0xed
    APPe = 0xee
    APPf = 0xef

    DQT  = 0xdb // Quantization Table

    SOF0 = 0xc0 // Start of Frame
    SOF1 = 0xc1
    SOF2 = 0xc2
    SOF3 = 0xc3
    SOF5 = 0xc5
    SOF6 = 0xc6
    SOF7 = 0xc7
    SOF9 = 0xc9
    SOFa = 0xca
    SOFb = 0xcb
    SOFd = 0xcd
    SOFe = 0xce
    SOFf = 0xcf

    SOS  = 0xda // Start of Scan

    DHT  = 0xc4 // Huffman Table

    JPG  = 0xc8
    JPG0 = 0xf0
    JPGd = 0xfd

    DAC  = 0xcc // Define Arithmetic Table
    DNL  = 0xdc
    DRI  = 0xdd // Define Restart Interval
    DHP  = 0xde
    EXP  = 0xdf

    RST0 = 0xd0
    RST1 = 0xd1
    RST2 = 0xd2
    RST3 = 0xd3
    RST4 = 0xd4
    RST5 = 0xd5
    RST6 = 0xd6
    RST7 = 0xd7

    TEM  = 0x01

    COM  = 0xfe // Comment

    EOI  = 0xd9 // End of Image
)

var EntryName = map[Byte]string{
    SOI: "SOI",
    APP0: "APP0",
    APP1: "APP1",
    APP2: "APP2",
    APP3: "APP3",
    APP4: "APP4",
    APP5: "APP5",
    APP6: "APP6",
    APP7: "APP7",
    APP8: "APP8",
    APP9: "APP9",
    APPa: "APPa",
    APPb: "APPb",
    APPc: "APPc",
    APPd: "APPd",
    APPe: "APPe",
    APPf: "APPf",
    DQT: "DQT",
    SOF0: "SOF0",
    SOF1: "SOF1",
    SOF2: "SOF2",
    SOF3: "SOF3",
    SOF5: "SOF5",
    SOF6: "SOF6",
    SOF7: "SOF7",
    SOF9: "SOF9",
    SOFa: "SOFa",
    SOFb: "SOFb",
    SOFd: "SOFd",
    SOFe: "SOFe",
    SOFf: "SOFf",
    DHT: "DHT",
    SOS: "SOS",

    JPG: "JPG",
    JPG0: "JPG0",
    JPGd: "JPGd",

    DAC: "DAC",
    DNL: "DNL",
    DRI: "DRI",
    DHP: "DHP",
    EXP: "EXP",

    RST0: "RST0",
    RST1: "RST1",
    RST2: "RST2",
    RST3: "RST3",
    RST4: "RST4",
    RST5: "RST5",
    RST6: "RST6",
    RST7: "RST7",

    TEM: "TEM",

    COM: "COM",

    EOI: "EOI",
}

type Entry interface {
    IsValid() bool
    // .Read(fd) fills the Entry and returns its "tail" (if Entry.Length exists)
    Read(fd *os.File) error
    HasData() bool
    GetData() []byte
    GetId() Byte
    String() string
    Write(Writer) error
    Pos() int64
    Len() int64
}

type lookupHeader struct {
    Xff0 Byte
    XID Byte
    Len Word
}
func (lkp *lookupHeader) IsValid() bool {
    return lkp.Xff0 == 255 &&
          (lkp.XID == SOI  ||
           lkp.XID == APP0 || lkp.XID == APP1 || lkp.XID == APP2 ||
           lkp.XID == APP3 || lkp.XID == APP4 || lkp.XID == APP5 ||
           lkp.XID == APP6 || lkp.XID == APP7 || lkp.XID == APP8 ||
           lkp.XID == APP9 || lkp.XID == APPa || lkp.XID == APPb ||
           lkp.XID == APPc || lkp.XID == APPd || lkp.XID == APPe ||
           lkp.XID == APPf ||
           lkp.XID == DQT  ||
           lkp.XID == SOF0 || lkp.XID == SOF1 || lkp.XID == SOF2 ||
           lkp.XID == SOF3 || lkp.XID == SOF5 || lkp.XID == SOF6 ||
           lkp.XID == SOF7 || lkp.XID == SOF9 || lkp.XID == SOFa ||
           lkp.XID == SOFb || lkp.XID == SOFd || lkp.XID == SOFe ||
           lkp.XID == SOFf ||
           lkp.XID == DHT  ||
           lkp.XID == SOS  ||
           lkp.XID == EOI)
}
func (lkp *lookupHeader) Read(fd *os.File) error {
    if e := ReadStructHere(fd, lkp); e != nil {
        return e
    }
    if !lkp.IsValid() {
        return fmt.Errorf("Invalid entry %+v", lkp)
    }
    return nil
}
func (lkp *lookupHeader) ReadData(fd *os.File) ([]byte, error) {
    var e error
    if lkp.Len == 0 {
        return nil, nil
    }
    pos, e := Tell(fd)
    if e != nil {
        return nil, e
    }
    if _, e = fd.Seek(pos + 4, 0); e != nil {
        return nil, e
    }
    data := make([]byte, lkp.Len - 2)
    if _, e = fd.Read(data); e != nil {
        return nil, e
    }
    return data, nil
}

type anEntry struct {
    Xff0 Byte   // +0
    XID Byte    // +1
}
func (ent *anEntry) IsValid() bool {
    return ent.Xff0 == 255 &&
          (ent.XID == SOI || ent.XID == APP0 || ent.XID == APP1 ||
           ent.XID == APP2 || ent.XID == APP3 || ent.XID == APP4 ||
           ent.XID == APP5 || ent.XID == APP6 || ent.XID == APP7 ||
           ent.XID == APP8 || ent.XID == APP9 || ent.XID == APPa ||
           ent.XID == APPb || ent.XID == APPc || ent.XID == APPd ||
           ent.XID == APPe || ent.XID == APPf || ent.XID == DQT ||
           ent.XID == SOF0 || ent.XID == SOF1 || ent.XID == SOF2 ||
           ent.XID == SOF3 || ent.XID == SOF5 || ent.XID == SOF6 ||
           ent.XID == SOF7 || ent.XID == SOF9 || ent.XID == SOFa ||
           ent.XID == SOFb || ent.XID == SOFd || ent.XID == SOFe ||
           ent.XID == SOFf ||
           ent.XID == DHT || ent.XID == SOS || ent.XID == EOI)
}
// .Read(fd) fills the Entry and returns its "tail" (if Entry.Length exists)
func (ent *anEntry) ReadEntry(fd *os.File) (Entry, error) {
    if e := ReadStructHere(fd, ent); e != nil {
        return nil, e
    }
    if !ent.IsValid() {
        return nil, fmt.Errorf("Invalid entry %+v", ent)
    }
    switch ent.XID {
    case SOI : { nent := new(SoiEntry); return nent, nent.Read(fd); }
    case APP0: { nent := new(App0Entry); return nent, nent.Read(fd); }
    case APP1, APP2, APP3, APP4, APP5, APP6, APP7, APP8, APP9, APPa, APPb, APPc,
         APPd, APPe, APPf:
         { nent := new(AppnEntry); return nent, nent.Read(fd); }
    case DQT : { nent := new(DqtEntry); return nent, nent.Read(fd); }
    case SOF0, SOF1, SOF2, SOF3, SOF5, SOF6, SOF7, SOFa, SOFb, SOFd, SOFe, SOFf:
        { nent := new(SofEntry); return nent, nent.Read(fd); }
    case DHT : { nent := new(DhtEntry); return nent, nent.Read(fd); }
    case SOS : { nent := new(SosEntry); return nent, nent.Read(fd); }
    case EOI : { nent := new(EoiEntry); return nent, nent.Read(fd); }
    }

    return nil, fmt.Errorf("Shit happened in entry %+v", ent)
}

type SoiEntry struct {
    Xff0 Byte   // +0
    ID Byte    // +1

    pos int64
}
func (soi *SoiEntry) HasData() bool { return false; }
func (soi *SoiEntry) Pos() int64 { return soi.pos; }
func (soi *SoiEntry) Len() int64 { return 2; }
func (soi *SoiEntry) GetId() Byte { return soi.ID; }
func (soi *SoiEntry) Write(fd Writer) error {
    var e error
    if e = WriteByte(fd, soi.Xff0); e != nil { return e; }
    if e = WriteByte(fd, soi.ID); e != nil { return e; }
    return nil
}
func (soi *SoiEntry) String() string { return fmt.Sprintf("<%s:%d>", EntryName[soi.ID], soi.Pos()); }
func (soi *SoiEntry) GetData() []byte { return nil; }
func (soi *SoiEntry) IsValid() bool {
    return soi.Xff0 == 255 && soi.ID == SOI
}
// .Read(fd) fills the Entry and returns its "tail" (if Entry.Length exists)
func (soi *SoiEntry) Read(fd *os.File) error {
    var e error
    if soi.pos, e = fd.Seek(0, 0); e != nil {
        return e
    }
    var tmp = make([]byte, 2)
    if _, e = fd.Read(tmp); e != nil {
        return e
    }
    soi.Xff0 = Byte(tmp[0])
    soi.ID = Byte(tmp[1])
    if !soi.IsValid() {
        return fmt.Errorf("Invalid SOI entry %+v", soi)
    }
    return nil // no .Length, no "tail"
}

type App0Entry struct {
    Xff0 Byte   // +0
    ID Byte   // +1
    Length Word // +2 // including the .Length field itself
    Identifier [5]byte  // +4 // 4A 46 49 46 00 == "JFIF\0"
    Version [2]Byte     // +9
    Units Byte          // +11
    Xdensity, Ydensity Word // +12, +14
    Xthumbnail, Ythumbnail Byte // +15, +16
    Data []byte

    pos int64
}
func (app0 *App0Entry) HasData() bool { return app0.Data != nil; }
func (app0 *App0Entry) Pos() int64 { return app0.pos; }
func (app0 *App0Entry) Len() int64 { return int64(app0.Length) + 2; }
func (app0 *App0Entry) GetId() Byte { return app0.ID; }
func (app0 *App0Entry) Write(fd Writer) error {
    var e error
    if e = WriteByte(fd, app0.Xff0); e != nil { return e; }
    if e = WriteByte(fd, app0.ID); e != nil { return e; }
    if e = WriteWordBE(fd, app0.Length); e != nil { return e; }
    if _, e = fd.Write(app0.Identifier[:]); e != nil { return e; }
    if e = WriteBytes(fd, app0.Version[:]); e != nil { return e; }
    if e = WriteByte(fd, app0.Units); e != nil { return e; }
    if e = WriteWordBE(fd, app0.Xdensity); e != nil { return e; }
    if e = WriteWordBE(fd, app0.Ydensity); e != nil { return e; }
    if e = WriteByte(fd, app0.Xthumbnail); e != nil { return e; }
    if e = WriteByte(fd, app0.Ythumbnail); e != nil { return e; }
    if _, e = fd.Write(app0.Data); e != nil { return e; }

    return nil
}
func (app0 *App0Entry) String() string {
    return fmt.Sprintf("<%s:%d[%d] %q Version=%d.%d Units=%d Xdensity=%d Ydensity=%d Xthumbnail=%d Ythumbnail=%d Data=%v>",
                       EntryName[app0.ID], app0.Pos(), app0.Length,
                       string(app0.Identifier[:]),
                       app0.Version[0], app0.Version[1],
                       app0.Units,
                       app0.Xdensity, app0.Ydensity,
                       app0.Xthumbnail, app0.Ythumbnail,
                       app0.Data)
}
func (app0 *App0Entry) GetData() []byte { return app0.Data; }
func (app0 *App0Entry) IsValid() bool {
    return app0.Xff0 == 255 && app0.ID == APP0 &&
           app0.Identifier[0] == 'J' && app0.Identifier[1] == 'F' &&
           app0.Identifier[2] == 'I' && app0.Identifier[3] == 'F' &&
           app0.Identifier[4] == 0 &&
           (app0.Units == 0 || app0.Units == 1 || app0.Units == 2)
}
// .Read(fd) fills the Entry and returns its "tail" (if Entry.Length exists)
func (app0 *App0Entry) Read(fd *os.File) error {
    app0.pos, _ = Tell(fd)
    var tmp struct {
        Xff0 Byte   // +0
        APP0 Byte   // +1
        Length Word // +2 // including the .Length field itself
        Identifier [5]byte  // +4 // 4A 46 49 46 00 == "JFIF\0"
        Version [2]Byte     // +9
        Units Byte          // +11
        Xdensity, Ydensity Word // +12, +14
        Xthumbnail, Ythumbnail Byte // +15, +16
    }
    e := binary.Read(fd, binary.BigEndian, &tmp)
    if e != nil { return e; }
    app0.Xff0 = tmp.Xff0
    app0.ID = tmp.APP0
    app0.Length = tmp.Length
    app0.Identifier = tmp.Identifier
    app0.Version = tmp.Version
    app0.Units = tmp.Units
    app0.Xdensity = tmp.Xdensity
    app0.Ydensity = tmp.Ydensity
    app0.Xthumbnail = tmp.Xthumbnail
    app0.Ythumbnail = tmp.Ythumbnail
    app0.Data = nil

    if !app0.IsValid() {
        return fmt.Errorf("Invalid header %+v", app0)
    }

    tsize := 3 * int(app0.Xthumbnail) * int(app0.Ythumbnail)
    if tsize > 0 {
        app0.Data = make([]byte, tsize)
        if _, e := fd.Read(app0.Data); e != nil {
            return e
        }
    }

    return nil
}

type AppnEntry struct {
    Xff0 Byte
    ID Byte
    Length Word
    Data []byte

    pos int64
}
func (app *AppnEntry) HasData() bool { return app.Data != nil; }
func (app *AppnEntry) Pos() int64 { return app.pos; }
func (app *AppnEntry) Len() int64 { return int64(app.Length) + 2; }
func (app *AppnEntry) GetId() Byte { return app.ID; }
func (app *AppnEntry) Write(fd Writer) error {
    var e error
    if e = WriteByte(fd, app.Xff0); e != nil { return e; }
    if e = WriteByte(fd, app.ID); e != nil { return e; }
    if e = WriteWordBE(fd, app.Length); e != nil { return e; }
    if _, e = fd.Write(app.Data); e != nil { return e; }

    return nil
}
func (app *AppnEntry) String() string { return fmt.Sprintf(EntryStringFormat, EntryName[app.ID], app); }
func (app *AppnEntry) GetData() []byte { return app.Data; }
func (app *AppnEntry) IsValid() bool {
    return app.Xff0 == 255 &&
          (app.ID == APP1 || app.ID == APP2 || app.ID == APP3 ||
           app.ID == APP4 || app.ID == APP5 || app.ID == APP6 ||
           app.ID == APP7 || app.ID == APP8 || app.ID == APP9 ||
           app.ID == APPa || app.ID == APPb || app.ID == APPc ||
           app.ID == APPd || app.ID == APPe || app.ID == APPf)
}
// .Read(fd) fills the Entry and returns its "tail" (if Entry.Length exists)
func (app *AppnEntry) Read(fd *os.File) error {
    app.pos, _ = Tell(fd)
    var e error
    var tmp lookupHeader
    if e = tmp.Read(fd); e != nil {
        return e
    }
    app.Xff0 = tmp.Xff0
    app.ID = tmp.XID
    app.Length = tmp.Len

    if !app.IsValid() {
        return fmt.Errorf("Invalid header %+v", app)
    }

    app.Data, e = tmp.ReadData(fd)
    if e != nil {
        return e
    }

    return nil
}

type DqtEntry struct {
    Xff0 Byte
    ID Byte
    Length Word
    Data []byte

    pos int64
}
func (dqt *DqtEntry) HasData() bool { return dqt.Data != nil; }
func (dqt *DqtEntry) Pos() int64 { return dqt.pos; }
func (dqt *DqtEntry) Len() int64 { return int64(dqt.Length) + 2; }
func (dqt *DqtEntry) GetId() Byte { return dqt.ID; }
func (dqt *DqtEntry) Write(fd Writer) error {
    var e error
    if e = WriteByte(fd, dqt.Xff0); e != nil { return e; }
    if e = WriteByte(fd, dqt.ID); e != nil { return e; }
    if e = WriteWordBE(fd, dqt.Length); e != nil { return e; }
    if _, e = fd.Write(dqt.Data); e != nil { return e; }

    return nil
}
func (dqt *DqtEntry) String() string {
    return fmt.Sprintf("<%s:%d[%d] data=%+v>", EntryName[dqt.ID], dqt.Pos(), dqt.Length, dqt.Data)
}
func (dqt *DqtEntry) GetData() []byte { return dqt.Data; }
func (dqt *DqtEntry) IsValid() bool {
    return dqt.Xff0 == 255 && dqt.ID == DQT
}
// .Read(fd) fills the Entry and returns its "tail" (if Entry.Length exists)
func (dqt *DqtEntry) Read(fd *os.File) error {
    dqt.pos, _ = Tell(fd)
    var e error
    var tmp lookupHeader
    if e = tmp.Read(fd); e != nil {
        return e
    }
    dqt.Xff0 = tmp.Xff0
    dqt.ID = tmp.XID
    dqt.Length = tmp.Len

    if !dqt.IsValid() {
        return fmt.Errorf("Invalid header %+v", dqt)
    }

    dqt.Data, e = tmp.ReadData(fd)
    if e != nil {
        return e
    }

    return nil
}

type SofEntry struct {
    Xff0 Byte
    ID Byte
    Length Word
    Precision Byte
    Height Word
    Width Word
    Components Byte
    Data []byte

    pos int64
}
func (sof *SofEntry) HasData() bool { return sof.Data != nil; }
func (sof *SofEntry) Pos() int64 { return sof.pos; }
func (sof *SofEntry) Len() int64 { return int64(sof.Length) + 2; }
func (sof *SofEntry) GetId() Byte { return sof.ID; }
func (sof *SofEntry) Write(fd Writer) error {
    var e error
    if e = WriteByte(fd, sof.Xff0); e != nil { return e; }
    if e = WriteByte(fd, sof.ID); e != nil { return e; }
    if e = WriteWordBE(fd, sof.Length); e != nil { return e; }
    if e = WriteByte(fd, sof.Precision); e != nil { return e; }
    if e = WriteWordLE(fd, sof.Height); e != nil { return e; }
    if e = WriteWordLE(fd, sof.Width); e != nil { return e; }
    if e = WriteByte(fd, sof.Components); e != nil { return e; }
    if _, e = fd.Write(sof.Data); e != nil { return e; }

    return nil
}
func (sof *SofEntry) String() string {
    return fmt.Sprintf("<%s:%d[%d] b/px=%d (H%d x W%d) %d:data=%+v>",
                       EntryName[sof.ID], sof.Pos(), sof.Length,
                       sof.Precision,
                       sof.Height, sof.Width,
                       sof.Components, sof.Data)
}
func (sof *SofEntry) GetData() []byte { return sof.Data; }
func (sof *SofEntry) IsValid() bool {
    return sof.Xff0 == 255 && (
           sof.ID == SOF0 || sof.ID == SOF1 || sof.ID == SOF2 ||
           sof.ID == SOF3 || sof.ID == SOF5 || sof.ID == SOF6 ||
           sof.ID == SOF7 || sof.ID == SOF9 || sof.ID == SOFa ||
           sof.ID == SOFb || sof.ID == SOFd || sof.ID == SOFe ||
           sof.ID == SOFf)
}
// .Read(fd) fills the Entry and returns its "tail" (if Entry.Length exists)
func (sof *SofEntry) Read(fd *os.File) error {
    sof.pos, _ = Tell(fd)
    var e error
    var tmp lookupHeader
    if e = tmp.Read(fd); e != nil {
        return e
    }
    sof.Xff0 = tmp.Xff0
    sof.ID = tmp.XID
    sof.Length = tmp.Len

    if !sof.IsValid() {
        return fmt.Errorf("Invalid header %+v", sof)
    }

    data, e := tmp.ReadData(fd)
    if e != nil {
        return e
    }
    sof.Precision, data = Byte(data[0]), data[1:]
    sof.Height, data = GetWordLE(data), data[2:]
    sof.Width, data = GetWordLE(data), data[2:]
    sof.Components, data = Byte(data[0]), data[1:]
    sof.Data = data

    return nil
}

type DhtEntry struct {
    Xff0 Byte
    ID Byte
    Length Word

    HtInfo Byte
    NumOfSymbols []byte

    Data []byte

    pos int64
}
func (dht *DhtEntry) HasData() bool { return dht.Data != nil; }
func (dht *DhtEntry) Pos() int64 { return dht.pos; }
func (dht *DhtEntry) Len() int64 { return int64(dht.Length) + 2; }
func (dht *DhtEntry) GetId() Byte { return dht.ID; }
func (dht *DhtEntry) Write(fd Writer) error {
    var e error
    if e = WriteByte(fd, dht.Xff0); e != nil { return e; }
    if e = WriteByte(fd, dht.ID); e != nil { return e; }
    if e = WriteWordBE(fd, dht.Length); e != nil { return e; }
    if e = WriteByte(fd, dht.HtInfo); e != nil { return e; }
    if _, e = fd.Write(dht.NumOfSymbols); e != nil { return e; }
    if _, e = fd.Write(dht.Data); e != nil { return e; }

    return nil
}
func (dht *DhtEntry) String() string {
    return fmt.Sprintf("<%s:%d[%d] %02x %v data=%+v>",
                       EntryName[dht.ID], dht.Pos(), dht.Length,
                       dht.HtInfo, dht.NumOfSymbols, dht.Data)
}
func (dht *DhtEntry) GetData() []byte { return dht.Data; }
func (dht *DhtEntry) IsValid() bool {
    return dht.Xff0 == 255 && dht.ID == DHT
}
// .Read(fd) fills the Entry and returns its "tail" (if Entry.Length exists)
func (dht *DhtEntry) Read(fd *os.File) error {
    dht.pos, _ = Tell(fd)
    var e error
    var tmp lookupHeader
    if e = tmp.Read(fd); e != nil {
        return e
    }
    dht.Xff0 = tmp.Xff0
    dht.ID = tmp.XID
    dht.Length = tmp.Len

    if !dht.IsValid() {
        return fmt.Errorf("Invalid header %+v", dht)
    }

    data, e := tmp.ReadData(fd)
    if e != nil {
        return e
    }
    dht.HtInfo, data = Byte(data[0]), data[1:]
    dht.NumOfSymbols, data = data[0:16], data[16:]
    dht.Data = data

    if len(dht.NumOfSymbols) != 16 { return fmt.Errorf("Bad NumOfSymbols in %+v", dht); }
    sum := 0
    for _, b := range dht.NumOfSymbols { sum += int(b); }
    if sum > 256 { return fmt.Errorf("Bad NumOfSymbols in %+v (%v > 256)", dht, sum); }
    if len(dht.Data) != sum { return fmt.Errorf("Bad NumOfSymbols in %+v: %v != %v", dht, len(dht.Data), sum); }

    return nil
}

type SosComponent struct {
    Id Byte
    Ht Byte
}
type SosEntry struct {
    Xff0 Byte
    ID Byte
    Length Word

    ComponentCount Byte
    Components []SosComponent

    Data []byte

    Image []byte

    pos int64
}
func (sos *SosEntry) HasData() bool { return sos.Data != nil; }
func (sos *SosEntry) Pos() int64 { return sos.pos; }
func (sos *SosEntry) Len() int64 { return int64(sos.Length) + 2; }
func (sos *SosEntry) GetId() Byte { return sos.ID; }
func (sos *SosEntry) Write(fd Writer) error {
    var e error
    if e = WriteByte(fd, sos.Xff0); e != nil { return e; }
    if e = WriteByte(fd, sos.ID); e != nil { return e; }
    if e = WriteWordBE(fd, sos.Length); e != nil { return e; }
    if e = WriteByte(fd, sos.ComponentCount); e != nil { return e; }

    for _, c := range sos.Components {
        if e = WriteByte(fd, c.Id); e != nil { return e; }
        if e = WriteByte(fd, c.Ht); e != nil { return e; }
    }

    if _, e = fd.Write(sos.Data); e != nil { return e; }
    if _, e = fd.Write(sos.Image); e != nil { return e; }

    return nil
}
func (sos *SosEntry) String() string {
    return fmt.Sprintf("<%s:%d[%d] components[%d]%v data=%v><IMAGE[%d]>",
                       EntryName[sos.ID], sos.Pos(), sos.Length,
                       sos.ComponentCount, sos.Components,
                       sos.Data,
                       len(sos.Image))
}
func (sos *SosEntry) GetData() []byte { return sos.Data; }
func (sos *SosEntry) IsValid() bool {
    return sos.Xff0 == 255 && sos.ID == SOS
}
// .Read(fd) fills the Entry and returns its "tail" (if Entry.Length exists)
func (sos *SosEntry) Read(fd *os.File) error {
    sos.pos, _ = Tell(fd)
    var e error
    var tmp lookupHeader
    if e = tmp.Read(fd); e != nil {
        return e
    }
    sos.Xff0 = tmp.Xff0
    sos.ID = tmp.XID
    sos.Length = tmp.Len

    if !sos.IsValid() {
        return fmt.Errorf("Invalid header %+v", sos)
    }

    data, e := tmp.ReadData(fd)
    if e != nil {
        return e
    }
    sos.ComponentCount, data = Byte(data[0]), data[1:]
    for i := 0; i < int(sos.ComponentCount); i++ {
        c := new (SosComponent)
        c.Id, c.Ht, data = Byte(data[0]), Byte(data[1]), data[2:]
        sos.Components = append(sos.Components, *c)
    }
    sos.Data = data

    var image []byte
    var oxff byte
    var b = make([]byte, 1)
    var img_size uint64 = 0

    for {
        _, e = fd.Read(b)
        if e != nil {
            return e
        }
        img_size += 1
        if img_size == 1 {
            oxff = b[0]
            continue
        }
        if oxff == 255 && b[0] == EOI {
            img_size -= 2
            _, e = fd.Seek(-2, 1)
            if e != nil {
                return e
            }
            break
        }
        image = append(image, oxff)
        oxff = b[0]
    } // here only via `break`
    sos.Image = image

    return nil
}

type EoiEntry struct {
    Xff0 Byte   // +0
    ID Byte    // +1

    pos int64
}
func (eoi *EoiEntry) HasData() bool { return false; }
func (eoi *EoiEntry) Pos() int64 { return eoi.pos; }
func (eoi *EoiEntry) Len() int64 { return 2; }
func (eoi *EoiEntry) GetId() Byte { return eoi.ID; }
func (eoi *EoiEntry) Write(fd Writer) error {
    var e error
    if e = WriteByte(fd, eoi.Xff0); e != nil { return e; }
    if e = WriteByte(fd, eoi.ID); e != nil { return e; }

    return nil
}
func (eoi *EoiEntry) String() string { return fmt.Sprintf("<%s:%d>", EntryName[eoi.ID], eoi.Pos()); }
func (eoi *EoiEntry) GetData() []byte { return nil; }
func (eoi *EoiEntry) IsValid() bool {
    return eoi.Xff0 == 255 && eoi.ID == EOI
}
// .Read(fd) fills the Entry and returns its "tail" (if Entry.Length exists)
func (eoi *EoiEntry) Read(fd *os.File) error {
    var e error
    if eoi.pos, e = Tell(fd); e != nil {
        return e
    }
    var tmp = make([]byte, 2)
    if _, e = fd.Read(tmp); e != nil {
        return e
    }
    eoi.Xff0 = Byte(tmp[0])
    eoi.ID = Byte(tmp[1])
    if !eoi.IsValid() {
        return fmt.Errorf("Invalid EOI entry %+v", eoi)
    }
    return nil // no .Length, no "tail"
}

/* vim: set ft=go ai et ts=4 sts=4 sw=4: EOF */
