package seed

import (
	"bytes"
	"context"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"yorgut/internal/storage"
)

// The seed pictures are generated, so no real photos live in the repo.

const (
	profilePictureSide = 240
	identiconCells     = 5
	communityWidth     = 260
	communityHeight    = 208
)

// randomColor picks a color around a random hue; light ones make
// backgrounds, saturated ones the shapes on top.
func randomColor(light bool) color.RGBA {
	hue := float64(RandomInt(360))
	if light {
		return hsl(hue, 0.35, 0.90)
	}
	return hsl(hue, 0.55+float64(RandomInt(30))/100, 0.40+float64(RandomInt(20))/100)
}

func hsl(h, s, l float64) color.RGBA {
	c := (1 - abs(2*l-1)) * s
	hp := h / 60
	x := c * (1 - abs(mod2(hp)-1))
	var r, g, b float64
	switch {
	case hp < 1:
		r, g = c, x
	case hp < 2:
		r, g = x, c
	case hp < 3:
		g, b = c, x
	case hp < 4:
		g, b = x, c
	case hp < 5:
		r, b = x, c
	default:
		r, b = c, x
	}
	m := l - c/2
	return color.RGBA{uint8((r + m) * 255), uint8((g + m) * 255), uint8((b + m) * 255), 255}
}

func abs(v float64) float64 {
	if v < 0 {
		return -v
	}
	return v
}

func mod2(v float64) float64 {
	for v >= 2 {
		v -= 2
	}
	return v
}

// identicon draws a mirrored grid of blocks, like the default avatars of
// many sites.
func identicon() image.Image {
	img := image.NewRGBA(image.Rect(0, 0, profilePictureSide, profilePictureSide))
	draw.Draw(img, img.Bounds(), &image.Uniform{randomColor(true)}, image.Point{}, draw.Src)

	fg := &image.Uniform{randomColor(false)}
	margin := profilePictureSide / 10
	cell := (profilePictureSide - 2*margin) / identiconCells
	for row := 0; row < identiconCells; row++ {
		for col := 0; col <= identiconCells/2; col++ {
			if !RandomBoolean() {
				continue
			}
			for _, c := range []int{col, identiconCells - 1 - col} {
				x, y := margin+c*cell, margin+row*cell
				draw.Draw(img, image.Rect(x, y, x+cell, y+cell), fg, image.Point{}, draw.Src)
			}
		}
	}
	return img
}

// shapes draws a few random rectangles and circles.
func shapes() image.Image {
	img := image.NewRGBA(image.Rect(0, 0, communityWidth, communityHeight))
	draw.Draw(img, img.Bounds(), &image.Uniform{randomColor(true)}, image.Point{}, draw.Src)

	for range 3 + RandomInt(4) {
		fg := &image.Uniform{randomColor(false)}
		size := 30 + RandomInt(90)
		x, y := RandomInt(communityWidth-size/2), RandomInt(communityHeight-size/2)
		if RandomBoolean() {
			draw.Draw(img, image.Rect(x, y, x+size, y+size*(1+RandomInt(2))/2), fg, image.Point{}, draw.Src)
			continue
		}
		r := size / 2
		for py := y - r; py <= y+r; py++ {
			for px := x - r; px <= x+r; px++ {
				if (px-x)*(px-x)+(py-y)*(py-y) <= r*r && image.Pt(px, py).In(img.Bounds()) {
					img.Set(px, py, fg.C)
				}
			}
		}
	}
	return img
}

// uploadGenerated stores a generated picture and returns its key.
func uploadGenerated(store storage.ObjectStore, img image.Image) *string {
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		panic(err)
	}
	key, err := storage.PutImage(context.Background(), store, buf.Bytes())
	if err != nil {
		panic(err)
	}
	return &key
}

// randomProfilePicture uploads a new identicon for a user.
func randomProfilePicture(store storage.ObjectStore) *string {
	return uploadGenerated(store, identicon())
}

// randomCommunityPicture uploads a new picture of shapes for a community.
func randomCommunityPicture(store storage.ObjectStore) *string {
	return uploadGenerated(store, shapes())
}
