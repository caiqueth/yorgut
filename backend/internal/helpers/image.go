package helpers

import (
	"bytes"
	"encoding/binary"
	"errors"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"

	"golang.org/x/image/draw"
)

// ImageDimensions reads the width and height of a PNG or JPEG from its header,
// without decoding the whole image.
func ImageDimensions(data []byte) (int, int, error) {
	config, _, err := image.DecodeConfig(bytes.NewReader(data))
	if err != nil {
		return 0, 0, err
	}
	return config.Width, config.Height, nil
}

// JPEG qualities tried in order, and how much a picture shrinks when even
// the lowest one is too big.
var (
	jpegQualities = []int{90, 85, 80, 75, 70, 60}
	shrinkStep    = 0.75
)

// ShrinkImage turns an uploaded PNG or JPEG into the picture we store: at most
// maxSide pixels on its longest side and maxBytes long.
//
//   - A PNG that already fits is kept as is, transparency included.
//   - Anything else is resized to fit maxSide and saved as a JPEG, at the best
//     quality that fits maxBytes; if none does, it shrinks further.
//   - JPEGs are always re-encoded: that applies the EXIF orientation phones
//     set (instead of rotating the pixels) and drops the EXIF metadata, GPS
//     location included.
//
// Callers check the dimensions first (ImageDimensions): decoding a huge image
// takes a lot of memory.
func ShrinkImage(data []byte, maxSide int, maxBytes int) ([]byte, error) {
	src, format, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	bounds := src.Bounds()

	if format == "png" && len(data) <= maxBytes && bounds.Dx() <= maxSide && bounds.Dy() <= maxSide {
		return data, nil
	}

	img := resizeToFit(src, maxSide)
	if format == "jpeg" {
		img = applyOrientation(img, jpegOrientation(data))
	}

	// A resized PNG stays a PNG (and keeps its transparency) if it fits.
	if format == "png" {
		var buf bytes.Buffer
		if err := png.Encode(&buf, img); err != nil {
			return nil, err
		}
		if buf.Len() <= maxBytes {
			return buf.Bytes(), nil
		}
	}

	img = onWhite(img)
	for {
		for _, quality := range jpegQualities {
			var buf bytes.Buffer
			if err := jpeg.Encode(&buf, img, &jpeg.Options{Quality: quality}); err != nil {
				return nil, err
			}
			if buf.Len() <= maxBytes {
				return buf.Bytes(), nil
			}
		}
		side := int(float64(max(img.Bounds().Dx(), img.Bounds().Dy())) * shrinkStep)
		if side < 64 {
			return nil, errors.New("image: cannot fit the size limit")
		}
		img = resizeToFit(img, side)
	}
}

// resizeToFit scales img down so its longest side is at most maxSide, keeping
// its proportions. Smaller images are returned as they are.
func resizeToFit(img image.Image, maxSide int) image.Image {
	width, height := img.Bounds().Dx(), img.Bounds().Dy()
	if width <= maxSide && height <= maxSide {
		return img
	}
	if width >= height {
		height = max(1, height*maxSide/width)
		width = maxSide
	} else {
		width = max(1, width*maxSide/height)
		height = maxSide
	}
	dst := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.CatmullRom.Scale(dst, dst.Bounds(), img, img.Bounds(), draw.Src, nil)
	return dst
}

// onWhite flattens transparency onto white, since JPEG has no alpha (a
// transparent PNG would otherwise turn black).
func onWhite(img image.Image) image.Image {
	dst := image.NewRGBA(image.Rect(0, 0, img.Bounds().Dx(), img.Bounds().Dy()))
	draw.Draw(dst, dst.Bounds(), image.NewUniform(color.White), image.Point{}, draw.Src)
	draw.Draw(dst, dst.Bounds(), img, img.Bounds().Min, draw.Over)
	return dst
}

// applyOrientation turns the pixels the way the EXIF orientation (1-8) says
// the image should be shown. 1 (or anything unknown) leaves it as it is.
func applyOrientation(img image.Image, orientation int) image.Image {
	if orientation < 2 || orientation > 8 {
		return img
	}
	b := img.Bounds()
	w, h := b.Dx(), b.Dy()
	// 5 to 8 swap width and height.
	dw, dh := w, h
	if orientation >= 5 {
		dw, dh = h, w
	}
	dst := image.NewRGBA(image.Rect(0, 0, dw, dh))
	for y := 0; y < dh; y++ {
		for x := 0; x < dw; x++ {
			// The source pixel shown at (x, y).
			var sx, sy int
			switch orientation {
			case 2: // flip left-right
				sx, sy = w-1-x, y
			case 3: // turn 180°
				sx, sy = w-1-x, h-1-y
			case 4: // flip top-bottom
				sx, sy = x, h-1-y
			case 5: // flip along the top-left diagonal
				sx, sy = y, x
			case 6: // turn 90° clockwise (most phone photos taken upright)
				sx, sy = y, h-1-x
			case 7: // flip along the top-right diagonal
				sx, sy = w-1-y, h-1-x
			case 8: // turn 90° counter-clockwise
				sx, sy = w-1-y, x
			}
			dst.Set(x, y, img.At(b.Min.X+sx, b.Min.Y+sy))
		}
	}
	return dst
}

// jpegOrientation reads the EXIF orientation tag (0x0112) of a JPEG, or 1
// when there is none or it can't be read.
func jpegOrientation(data []byte) int {
	if len(data) < 4 || data[0] != 0xFF || data[1] != 0xD8 {
		return 1
	}
	// Walk the segments after the start-of-image marker up to the image data.
	for i := 2; i+4 <= len(data); {
		if data[i] != 0xFF {
			return 1
		}
		marker := data[i+1]
		if marker == 0xDA { // start of scan: no EXIF before the pixels
			return 1
		}
		length := int(binary.BigEndian.Uint16(data[i+2:]))
		segment := i + 4
		end := i + 2 + length
		if length < 2 || end > len(data) {
			return 1
		}
		if marker == 0xE1 && end-segment > 6 && string(data[segment:segment+6]) == "Exif\x00\x00" {
			return tiffOrientation(data[segment+6 : end])
		}
		i = end
	}
	return 1
}

// tiffOrientation reads the orientation from the TIFF block inside EXIF.
func tiffOrientation(tiff []byte) int {
	if len(tiff) < 8 {
		return 1
	}
	var order binary.ByteOrder
	switch string(tiff[:2]) {
	case "II":
		order = binary.LittleEndian
	case "MM":
		order = binary.BigEndian
	default:
		return 1
	}
	ifd := int(order.Uint32(tiff[4:]))
	if ifd+2 > len(tiff) {
		return 1
	}
	entries := int(order.Uint16(tiff[ifd:]))
	for n := 0; n < entries; n++ {
		entry := ifd + 2 + n*12
		if entry+12 > len(tiff) {
			return 1
		}
		if order.Uint16(tiff[entry:]) == 0x0112 {
			return int(order.Uint16(tiff[entry+8:]))
		}
	}
	return 1
}
