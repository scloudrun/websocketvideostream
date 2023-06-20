package main

import (
	"bytes"
	"encoding/binary"
	"errors"
	"image"
	"image/draw"
	"io/ioutil"
	"net"
	"os"
	"unsafe"
	"websocket/internal/encoders"
	"websocket/internal/rdisplay"

	"github.com/nfnt/resize"
)

//go build -tags "h264enc" cmd/h264_encode.go

var (
	Conn    net.Conn
	Encoder encoders.Encoder
)

func getEncode(v string) ([]byte, error) {
	size, err := Encoder.VideoSize()
	if err != nil {
		return nil, err
	}
	frame, err := getImage(v) //30ms - 50ms
	if err != nil {
		return nil, err
	}
	resized := resizeImage(frame, size)     //80ms-110ms
	payload, err := Encoder.Encode(resized) //50ms - 100ms
	return payload, err
}

func initEncoder() error {
	var (
		err error
		enc encoders.Service = &encoders.EncoderService{}
	)

	screen := rdisplay.Screen{Index: 0, Bounds: image.Rectangle{Min: image.Point{0, 0}, Max: image.Point{1920, 1080}}}
	sourceSize := image.Point{
		screen.Bounds.Dx(),
		screen.Bounds.Dy(),
	}

	Encoder, err = enc.NewEncoder(1, sourceSize, String2Int(frame))
	if err != nil {
		return err
	}

	return nil
}

func resizeImage(src *image.RGBA, target image.Point) *image.RGBA {
	return resize.Resize(uint(target.X), uint(target.Y), src, resize.Lanczos3).(*image.RGBA)
}

func getImage(filePath string) (*image.RGBA, error) {
	// Decoding JPEG into image.Image
	imgFile, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer imgFile.Close()
	imgFile.Seek(0, 0)

	//jpg 图片的开始是ffd8,结束是ffd9 //https://github.com/corkami/formats/blob/master/image/jpeg.md
	//xxd examples.jpeg |egrep "ffd9|ff d9"
	//xxd examples_flow.jpeg |egrep "ffd9|ff d9"

	contents, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	if len(contents) < 4 {
		return nil, errors.New("below 4 byte")
	}

	// Maybe wrong End-Of-Image.
	if !(contents[0] == '\xff' || contents[1] == '\xd8') {
		return nil, err
	}
	if !(contents[len(contents)-1] == '\xd9' && contents[len(contents)-2] == '\xff') {
		return nil, err
	}

	// Decode the JPEG data. If reading from file, create a reader with
	img, _, err := image.Decode(imgFile)
	if err != nil {
		return nil, err
	}

	rgba := imageToRGBA(img)
	return rgba, err
}

func imageToRGBA(src image.Image) *image.RGBA {
	// No conversion needed if image is an *image.RGBA.
	if dst, ok := src.(*image.RGBA); ok {
		return dst
	}

	// Use the image/draw package to convert to *image.RGBA.
	b := src.Bounds()
	dst := image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
	draw.Draw(dst, dst.Bounds(), src, b.Min, draw.Src)
	return dst
}

// BytesToUInt32 ...
func BytesToUInt32(buf []byte) uint32 {
	return uint32(binary.BigEndian.Uint32(buf))
}

// Bytes2Str ...
func Bytes2Str(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

// Int32ToBytes ...
func Int32ToBytes(i uint32) []byte {
	var buf = make([]byte, 4)
	binary.BigEndian.PutUint32(buf, uint32(i))
	return buf
}

func bytesMerge(pBytes ...[]byte) []byte {
	return bytes.Join(pBytes, []byte(""))
}

// 整形转换成字节
func IntToBytes(n int) []byte {
	x := int32(n)
	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.BigEndian, x)
	return bytesBuffer.Bytes()
}

// 字节转换成整形
func BytesToInt(b []byte) int {
	bytesBuffer := bytes.NewBuffer(b)

	var x int32
	binary.Read(bytesBuffer, binary.BigEndian, &x)

	return int(x)
}
