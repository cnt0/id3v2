package id3v2

import (
	"bytes"
	"github.com/bogem/id3v2/frame"
	"os"
	"sync"
)

const (
	// Picture types
	PTOther                   = 0
	PTFileIcon                = 1
	PTOtherFileIcon           = 2
	PTFrontCover              = 3
	PTBackCover               = 4
	PTLeafletPage             = 5
	PTMedia                   = 6
	PTLeadArtistSoloist       = 7
	PTArtistPerformer         = 8
	PTConductor               = 9
	PTBandOrchestra           = 10
	PTComposer                = 11
	PTLyricistTextWriter      = 12
	PTRecordingLocation       = 13
	PTDuringRecording         = 14
	PTDuringPerformance       = 15
	PTMovieScreenCapture      = 16
	PTBrightColouredFish      = 17
	PTIllustration            = 18
	PTBandArtistLogotype      = 19
	PTPublisherStudioLogotype = 20
)

var bytesBufPool = sync.Pool{
	New: func() interface{} { return new(bytes.Buffer) },
}

func Open(name string) (*Tag, error) {
	file, err := os.Open(name)
	if err != nil {
		return nil, err
	}

	return ParseTag(file)
}

func NewAttachedPicture() *frame.PictureFrame {
	return new(frame.PictureFrame)
}
