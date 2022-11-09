package util 

import (
  "io"
	"os"
  "image/jpeg"
	"github.com/adrium/goheif"
)

func CovertHeicToJpg(inputPath, outputPath string) error {
  fileInput, err := os.Open(inputPath)
  if err != nil {
    return err
  }
  defer fileInput.Close()
  
  exif, err := goheif.ExtractExif(fileInput)
	if err != nil {
		return err
	}
  img, err := goheif.Decode(fileInput)
  if err != nil {
    return err
  }
  fileOutput, err := os.OpenFile(outputPath, os.O_RDWR | os.O_CREATE, 0644)
  if err != nil {
    return err
  }
  defer fileOutput.Close()
  
  w, _ := newWriterExif(fileOutput, exif)
  err = jpeg.Encode(w, img, nil)
  if err != nil {
    return err
  }
  return nil
}

type writerSkipper struct {
	w           io.Writer
	bytesToSkip int
}

func newWriterExif(w io.Writer, exif []byte) (io.Writer, error) {
	writer := &writerSkipper{w, 2}
	soi := []byte{0xff, 0xd8}
	if _, err := w.Write(soi); err != nil {
		return nil, err
	}

	if exif != nil {
		app1Marker := 0xe1
		markerlen := 2 + len(exif)
		marker := []byte{0xff, uint8(app1Marker), uint8(markerlen >> 8), uint8(markerlen & 0xff)}
		if _, err := w.Write(marker); err != nil {
			return nil, err
		}

		if _, err := w.Write(exif); err != nil {
			return nil, err
		}
	}

	return writer, nil
}

func (w *writerSkipper) Write(data []byte) (int, error) {
	if w.bytesToSkip <= 0 {
		return w.w.Write(data)
	}

	if dataLen := len(data); dataLen < w.bytesToSkip {
		w.bytesToSkip -= dataLen
		return dataLen, nil
	}

	if n, err := w.w.Write(data[w.bytesToSkip:]); err == nil {
		n += w.bytesToSkip
		w.bytesToSkip = 0
		return n, nil
	} else {
		return n, err
	}
}

