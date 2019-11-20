package jfif

import (
    "os"
    "path"
    "io/ioutil"
)
import "testing"

const testImagePath = "assets"
const testOutputSuffix = ".jfif"

func compare(t *testing.T, x *Jfif, p1, p2 string) bool {
    d1, e := ioutil.ReadFile(p1)
    if e != nil {
        t.Fatalf("ioutil.ReadFile(%q): %v", p1, e)
    }
    d2, e := ioutil.ReadFile(p2)
    if e != nil {
        t.Fatalf("ioutil.ReadFile(%q): %v", p2, e)
    }
    if len(d1) != len(d2) {
        t.Errorf("Size mismatch: %d vs %d", len(d1), len(d2))
        return false
    }
    ok := true
    for i, b1 := range d1 {
        b2 := d2[i]
        if b1 != b2 {
            s := x.SectionAt(int64(i))
            t.Logf("%5d %08x %02x %02x @ %v", i, i, b1, b2, s)
            ok = false
        }
    }
    return ok
}

func TestInject(t *testing.T) {
    var x = JfifData{
        "UserComment": String("sample user comment"),
        "gps.latitude": Rational{Num:550, Den:10},
    }
    // imgPath := path.Join(os.TempDir(), "go-basler-pylon-test")
    imgPath := testImagePath
    var e error
    var f *os.File
    var lst []string
    if f, e = os.Open(imgPath); e != nil {
        t.Fatalf("Cannot open %#v: %v", imgPath, e)
    }
    defer f.Close()
    if lst, e = f.Readdirnames(0); e != nil {
        t.Fatalf("Cannot read %#v: %v", imgPath, e)
    }
    if len(lst) == 0 {
        t.Fatalf("No files in %#v", imgPath)
    }
    for _, name := range lst {
        if path.Ext(name) != ".jpg" {
            t.Logf("Skipping %#v", name)
            continue
        }
        path := path.Join(imgPath, name)
        var X Jfif
        if e = X.Load(path); e != nil {
            t.Fatalf("Cannot load %#v: %v", path, e)
        }
        if e = X.Inject(x); e != nil {
            t.Fatalf("Cannot Inject(%#v, %#v): %v", path, x, e)
        }
        if e = X.SaveTo(path + testOutputSuffix); e != nil {
            t.Fatalf("Cannot save %#v: %v", path + testOutputSuffix, e)
        } else {
            t.Logf("File %#v saved", path)
        }
        if compare(t, &X, path, path + testOutputSuffix) {
            t.Logf("same files")
        } else {
            t.Fatalf("Files differ")
        }
    }
}
/* vim: set ft=go ai et ts=4 sts=4 sw=4: EOF */
