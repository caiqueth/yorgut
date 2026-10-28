package helpers

import (
	"bytes"
	"encoding/binary"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"math/rand"
	"testing"
)

// noise is an image JPEG and PNG can't compress well, to force big files.
func noise(width, height int) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	rng := rand.New(rand.NewSource(1))
	rng.Read(img.Pix)
	for i := 3; i < len(img.Pix); i += 4 {
		img.Pix[i] = 255
	}
	return img
}

func encodePNG(t *testing.T, img image.Image) []byte {
	t.Helper()
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		t.Fatal(err)
	}
	return buf.Bytes()
}

func encodeJPEG(t *testing.T, img image.Image) []byte {
	t.Helper()
	var buf bytes.Buffer
	if err := jpeg.Encode(&buf, img, &jpeg.Options{Quality: 95}); err != nil {
		t.Fatal(err)
	}
	return buf.Bytes()
}

// withOrientation inserts an EXIF block with the given orientation right
// after the JPEG's start-of-image marker, like a phone camera does.
func withOrientation(data []byte, orientation uint16) []byte {
	tiff := []byte("II*\x00\x08\x00\x00\x00")                  // little endian, first IFD at 8
	tiff = binary.LittleEndian.AppendUint16(tiff, 1)           // one entry
	tiff = binary.LittleEndian.AppendUint16(tiff, 0x0112)      // orientation
	tiff = binary.LittleEndian.AppendUint16(tiff, 3)           // SHORT
	tiff = binary.LittleEndian.AppendUint32(tiff, 1)           // one value
	tiff = binary.LittleEndian.AppendUint16(tiff, orientation) // the value
	tiff = append(tiff, 0, 0, 0, 0, 0, 0)                      // padding, no next IFD
	payload := append([]byte("Exif\x00\x00"), tiff...)

	segment := []byte{0xFF, 0xE1}
	segment = binary.BigEndian.AppendUint16(segment, uint16(len(payload)+2))
	segment = append(segment, payload...)

	out := append([]byte{}, data[:2]...)
	out = append(out, segment...)
	return append(out, data[2:]...)
}

func decode(t *testing.T, data []byte) (image.Image, string) {
	t.Helper()
	img, format, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		t.Fatal(err)
	}
	return img, format
}

func TestShrinkImage(t *testing.T) {
	const maxSide, maxBytes = 1024, 1 << 20

	t.Run("a big photo fits both limits and keeps its proportions", func(t *testing.T) {
		got, err := ShrinkImage(encodeJPEG(t, noise(3000, 2000)), maxSide, maxBytes)
		if err != nil {
			t.Fatal(err)
		}
		img, format := decode(t, got)
		if format != "jpeg" || len(got) > maxBytes {
			t.Fatalf("got %s of %d bytes, want a JPEG of at most %d", format, len(got), maxBytes)
		}
		if b := img.Bounds(); b.Dx() != 1024 || b.Dy() != 682 {
			t.Errorf("got %dx%d, want 1024x682", b.Dx(), b.Dy())
		}
	})

	t.Run("a PNG that fits is kept byte for byte", func(t *testing.T) {
		original := encodePNG(t, image.NewRGBA(image.Rect(0, 0, 200, 100)))
		got, err := ShrinkImage(original, maxSide, maxBytes)
		if err != nil || !bytes.Equal(got, original) {
			t.Fatalf("got %d bytes (err %v), want the original %d", len(got), err, len(original))
		}
	})

	t.Run("a big PNG that compresses well stays a PNG", func(t *testing.T) {
		transparent := image.NewRGBA(image.Rect(0, 0, 2048, 1024))
		got, err := ShrinkImage(encodePNG(t, transparent), maxSide, maxBytes)
		if err != nil {
			t.Fatal(err)
		}
		img, format := decode(t, got)
		if b := img.Bounds(); format != "png" || b.Dx() != 1024 || b.Dy() != 512 {
			t.Fatalf("got a %dx%d %s, want a 1024x512 png", b.Dx(), b.Dy(), format)
		}
	})

	t.Run("a PNG too big even resized becomes a JPEG", func(t *testing.T) {
		got, err := ShrinkImage(encodePNG(t, noise(1024, 1024)), maxSide, maxBytes)
		if err != nil {
			t.Fatal(err)
		}
		if _, format := decode(t, got); format != "jpeg" || len(got) > maxBytes {
			t.Fatalf("got %s of %d bytes, want a JPEG of at most %d", format, len(got), maxBytes)
		}
	})

	t.Run("a photo taken upright on a phone comes out upright", func(t *testing.T) {
		// Stored sideways, 40x20: left half red, right half blue. Orientation
		// 6 means "turn 90° clockwise to show it", which puts red on top.
		stored := image.NewRGBA(image.Rect(0, 0, 40, 20))
		for y := 0; y < 20; y++ {
			for x := 0; x < 40; x++ {
				c := color.RGBA{R: 255, A: 255}
				if x >= 20 {
					c = color.RGBA{B: 255, A: 255}
				}
				stored.Set(x, y, c)
			}
		}
		got, err := ShrinkImage(withOrientation(encodeJPEG(t, stored), 6), maxSide, maxBytes)
		if err != nil {
			t.Fatal(err)
		}
		img, _ := decode(t, got)
		if b := img.Bounds(); b.Dx() != 20 || b.Dy() != 40 {
			t.Fatalf("got %dx%d, want 20x40", b.Dx(), b.Dy())
		}
		top, _, _, _ := img.At(10, 5).RGBA()
		_, _, bottom, _ := img.At(10, 35).RGBA()
		if top < 0xC000 || bottom < 0xC000 {
			t.Errorf("top red %x, bottom blue %x: want red on top and blue at the bottom", top, bottom)
		}
	})
}

func TestJPEGOrientation(t *testing.T) {
	plain := encodeJPEG(t, image.NewRGBA(image.Rect(0, 0, 8, 8)))
	if got := jpegOrientation(plain); got != 1 {
		t.Errorf("no EXIF: got %d, want 1", got)
	}
	for _, want := range []uint16{3, 6, 8} {
		if got := jpegOrientation(withOrientation(plain, want)); got != int(want) {
			t.Errorf("got %d, want %d", got, want)
		}
	}
	if got := jpegOrientation([]byte("not a jpeg")); got != 1 {
		t.Errorf("garbage: got %d, want 1", got)
	}
}
