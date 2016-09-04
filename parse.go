// Copyright 2016 Albert Nigmatzianov. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package id3v2

import (
	"errors"
	"io"
	"os"

	"github.com/bogem/id3v2/bbpool"
	"github.com/bogem/id3v2/util"
)

const frameHeaderSize = 10

type frameHeader struct {
	ID        string
	FrameSize int64
}

func parseTag(file *os.File) (*Tag, error) {
	if file == nil {
		err := errors.New("Invalid file: file is nil")
		return nil, err
	}
	header, err := parseHeader(file)
	if err != nil {
		err = errors.New("Trying to parse tag header: " + err.Error())
		return nil, err
	}
	if header == nil {
		return newTag(file, 0, 4), nil
	}
	if header.Version < 3 {
		err = errors.New("Unsupported version of ID3 tag")
		return nil, err
	}

	t := newTag(file, tagHeaderSize+header.FramesSize, header.Version)
	err = t.parseAllFrames()

	return t, err
}

func newTag(file *os.File, originalSize int64, version byte) *Tag {
	t := &Tag{
		frames:    make(map[string]Framer),
		sequences: make(map[string]sequencer),

		file:         file,
		originalSize: originalSize,
		version:      version,
	}

	if version == 3 {
		t.ids = V23IDs
	} else {
		t.ids = V24IDs
	}

	return t
}

func (t *Tag) parseAllFrames() error {
	size := t.originalSize - tagHeaderSize // Size of all frames = Size of tag - tag header
	f := t.file

	// Initial position of read - beginning of first frame
	if _, err := f.Seek(tagHeaderSize, os.SEEK_SET); err != nil {
		return err
	}

	limFile := &io.LimitedReader{R: f}
	for size > 0 {
		limFile.N = frameHeaderSize
		header, err := parseFrameHeader(limFile)
		if err != nil {
			return err
		}

		parseFunc := t.findParseFunc(header.ID)

		limFile.N = header.FrameSize
		frame, err := parseFunc(limFile)
		if err != nil {
			return err
		}

		t.AddFrame(header.ID, frame)

		size -= frameHeaderSize + header.FrameSize
	}

	return nil
}

func parseFrameHeader(rd io.Reader) (*frameHeader, error) {
	fhBuf := bbpool.Get()
	defer bbpool.Put(fhBuf)

	n, err := fhBuf.ReadFrom(rd)
	if err != nil {
		return nil, err
	}
	if n != frameHeaderSize {
		return nil, errors.New("Unexpected frame header size")
	}

	byteFH := fhBuf.Bytes()

	header := &frameHeader{
		ID:        string(byteFH[:4]),
		FrameSize: util.ParseSize(byteFH[4:8]),
	}

	return header, nil

}

func (t Tag) findParseFunc(id string) func(io.Reader) (Framer, error) {
	if id[0] == 'T' {
		return parseTextFrame
	}

	switch id {
	case t.ID("Attached picture"):
		return parsePictureFrame
	case t.ID("Comments"):
		return parseCommentFrame
	case t.ID("Unsynchronised lyrics/text transcription"):
		return parseUnsynchronisedLyricsFrame
	}
	return parseUnknownFrame
}