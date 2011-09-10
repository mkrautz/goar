// Copyright (c) 2011 Mikkel Krautz
// The use of this source code is goverened by a BSD-style
// license that can be found in the LICENSE-file.

package ar

import (
	"bytes"
	"io"
	"os"
	"testing"
)

var fbsd82Archive []archiveTest = []archiveTest{
	{
		&Header{
			Name:  "/",
			Mode:  0,
			Mtime: 1315607407,
			Uid:   0,
			Gid:   0,
			Size:  4,
		},
		[]byte{0x0, 0x0, 0x0, 0x0},
	},
	{
		&Header{
			Name:  "a",
			Mode:  0100644,
			Mtime: 1315607373,
			Uid:   1001,
			Gid:   1001,
			Size:  2,
		},
		[]byte{'a', '\n'},
	},
	{
		&Header{
			Name:  "b",
			Mode:  0100644,
			Mtime: 1315607374,
			Uid:   1001,
			Gid:   1001,
			Size:  2,
		},
		[]byte{'b', '\n'},
	},
	{
		&Header{
			Name:  "c",
			Mode:  0100644,
			Mtime: 1315607376,
			Uid:   1001,
			Gid:   1001,
			Size:  2,
		},
		[]byte{'c', '\n'},
	},
}

var lionArchive []archiveTest = []archiveTest{
	{
		&Header{
			Name:  "__.SYMDEF SORTED",
			Mode:  0100644,
			Mtime: 1315593186,
			Uid:   501,
			Gid:   20,
			Size:  8,
		},
		[]byte{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
	},
	{
		&Header{
			Name:  "a",
			Mode:  0100644,
			Mtime: 1315593158,
			Uid:   501,
			Gid:   20,
			Size:  8,
		},
		[]byte("a\n\n\n\n\n\n\n"),
	},
	{
		&Header{
			Name:  "b",
			Mode:  0100644,
			Mtime: 1315593165,
			Uid:   501,
			Gid:   20,
			Size:  8,
		},
		[]byte("b\n\n\n\n\n\n\n"),
	},
	{
		&Header{
			Name:  "c",
			Mode:  0100644,
			Mtime: 1315593166,
			Uid:   501,
			Gid:   20,
			Size:  8,
		},
		[]byte("c\n\n\n\n\n\n\n"),
	},
}

func testRead(t *testing.T, r io.Reader, testArchive []archiveTest) {
	ar := NewReader(r)
	for _, testEntry := range testArchive {
		hdr, err := ar.Next()
		if err != nil {
			t.Fatal(err)
		}
		if !headerCmp(hdr, testEntry.hdr) {
			t.Fatalf("header mismatch:\nread = %v\norig = %v", hdr, testEntry.hdr)
		}
		fbuf := make([]byte, hdr.Size)
		_, err = io.ReadFull(ar, fbuf)
		if err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(fbuf, testEntry.data) {
			t.Fatalf("data mismatch\nread = %v\norig = %v", fbuf, testEntry.data)
		}
	}
}

// Test the we can correctly read and parse a FreeBSD 8.2 generated ar file.
func TestReadFreeBSD82LibArchive(t *testing.T) {
	f, err := os.Open("testdata/test-bsd-freebsd82-libarchive.ar")
	if err != nil {
		t.Error(err)
		return
	}
	defer f.Close()
	testRead(t, f, fbsd82Archive)
}

// Test the we can correctly read and parse a Mac OS X Lion generated ar file.
// It is generated in the same way as the FreeBSD archive ahove, but ar on OS X
// seems to pad the archived files with a lot of newlines. 
// Attempting to "ar x" the archive also reproduces the newlines in the extracted
// files, so they are not a form of padding, but are intended to be there, somehow.
func TestReadMacOSXLionOld(t *testing.T) {
	f, err := os.Open("testdata/test-bsd-macosx.ar")
	if err != nil {
		t.Error(err)
		return
	}
	defer f.Close()
	testRead(t, f, lionArchive)
}

// Test that we can correctly skip data bytes in a file
func TestDataSkipping(t *testing.T) {
	f, err := os.Open("testdata/test-bsd-freebsd82-libarchive.ar")
	if err != nil {
		t.Error(err)
		return
	}
	r := NewReader(f)
	for _, testEntry := range fbsd82Archive {
		hdr, err := r.Next()
		if err != nil {
			t.Error(err)
			return
		}
		if !headerCmp(hdr, testEntry.hdr) {
			t.Fatalf("header mismatch:\nread = %v\norig = %v", hdr, testEntry.hdr)
		}
	}
}
