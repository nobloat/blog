//go:build image

package main

import (
	"image"
	"image/color"
	_ "image/jpeg"
	"image/png"
	"log"
	"math"
	"os"
	"path"
	"path/filepath"
	"strings"
)

const (
	maxLongEdge     = 400
	unsharpSigma    = 1.5
	unsharpAmount   = 1.5
	sigmoidContrast = 5.0
	sigmoidMidpoint = 0.5
	stretchBlack    = 0.02
	stretchWhite    = 0.98
	ditherThreshold = 127
)

func runImageCommand(args []string) {
	if len(args) < 1 {
		log.Fatal("Usage: go run main.go image <input> [output]")
	}

	in := args[0]
	out := path.Join("public", "images", strings.TrimSuffix(filepath.Base(in), filepath.Ext(in))+".png")

	inStat, err := os.Stat(in)
	if err != nil {
		log.Fatal(err)
	}
	inSize := inStat.Size()

	f, err := os.Open(in)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		log.Fatal(err)
	}

	img = resizeLongEdge(img, maxLongEdge)
	gray := toGrayscale(img)
	bw := dither(gray)

	o, err := os.Create(out)
	if err != nil {
		log.Fatal(err)
	}
	defer o.Close()

	enc := &png.Encoder{CompressionLevel: png.BestCompression}
	enc.Encode(o, bw)
	o.Close()

	outStat, err := os.Stat(out)
	if err != nil {
		log.Fatal(err)
	}
	outSize := outStat.Size()

	log.Printf("Input: %s (%.2f MB)", in, float64(inSize)/(1024*1024))
	log.Printf("Output: %s (%.2f MB)", out, float64(outSize)/(1024*1024))
	log.Printf("Reduction: %.1f%%", 100.0*(1.0-float64(outSize)/float64(inSize)))
}

func toGrayscale(img image.Image) *image.Gray {
	b := img.Bounds()
	g := image.NewGray(b)

	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			r, gr, bl, _ := img.At(x, y).RGBA()
			luma := uint8((0.2126*float64(r) + 0.7152*float64(gr) + 0.0722*float64(bl)) / 256)
			g.SetGray(x, y, color.Gray{Y: luma})
		}
	}

	g = unsharp(g, unsharpSigma, unsharpAmount)
	g = sigmoid(g, sigmoidContrast, sigmoidMidpoint)
	g = stretch(g, stretchBlack, stretchWhite)
	return g
}

func unsharp(img *image.Gray, sigma, amt float64) *image.Gray {
	b := img.Bounds()
	blur := boxBlur(img, int(sigma*2))
	out := image.NewGray(b)

	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			orig := float64(img.GrayAt(x, y).Y)
			blr := float64(blur.GrayAt(x, y).Y)
			out.SetGray(x, y, color.Gray{Y: clamp(orig + amt*(orig-blr))})
		}
	}
	return out
}

func boxBlur(img *image.Gray, r int) *image.Gray {
	b := img.Bounds()
	out := image.NewGray(b)

	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			var sum, n float64
			for dy := -r; dy <= r; dy++ {
				for dx := -r; dx <= r; dx++ {
					nx, ny := x+dx, y+dy
					if nx >= b.Min.X && nx < b.Max.X && ny >= b.Min.Y && ny < b.Max.Y {
						sum += float64(img.GrayAt(nx, ny).Y)
						n++
					}
				}
			}
			out.SetGray(x, y, color.Gray{Y: uint8(sum / n)})
		}
	}
	return out
}

func sigmoid(img *image.Gray, contrast, mid float64) *image.Gray {
	b := img.Bounds()
	out := image.NewGray(b)

	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			v := float64(img.GrayAt(x, y).Y) / 255.0
			adj := 1.0 / (1.0 + math.Exp(contrast*(mid-v)))
			out.SetGray(x, y, color.Gray{Y: uint8(adj * 255.0)})
		}
	}
	return out
}

func stretch(img *image.Gray, black, white float64) *image.Gray {
	b := img.Bounds()
	var min, max uint8 = 255, 0

	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			v := img.GrayAt(x, y).Y
			if v < min {
				min = v
			}
			if v > max {
				max = v
			}
		}
	}

	minF := float64(min) + float64(max-min)*black
	maxF := float64(min) + float64(max-min)*white
	rng := maxF - minF
	if rng == 0 {
		return img
	}

	out := image.NewGray(b)
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			v := float64(img.GrayAt(x, y).Y)
			out.SetGray(x, y, color.Gray{Y: clamp((v - minF) / rng * 255.0)})
		}
	}
	return out
}

func clamp(v float64) uint8 {
	if v < 0 {
		return 0
	}
	if v > 255 {
		return 255
	}
	return uint8(v)
}

func resizeLongEdge(img image.Image, maxLongEdge int) image.Image {
	b := img.Bounds()
	w, h := b.Dx(), b.Dy()

	var nw, nh int
	if w > h {
		nw = maxLongEdge
		nh = int(float64(h) * float64(maxLongEdge) / float64(w))
	} else {
		nh = maxLongEdge
		nw = int(float64(w) * float64(maxLongEdge) / float64(h))
	}

	out := image.NewRGBA(image.Rect(0, 0, nw, nh))

	for y := 0; y < nh; y++ {
		for x := 0; x < nw; x++ {
			sx := int(float64(x) * float64(w) / float64(nw))
			sy := int(float64(y) * float64(h) / float64(nh))
			out.Set(x, y, img.At(sx, sy))
		}
	}

	return out
}

func dither(img *image.Gray) *image.Gray {
	b := img.Bounds()
	out := image.NewGray(b)

	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			old := img.GrayAt(x, y).Y
			var new uint8
			if old > ditherThreshold {
				new = 255
			} else {
				new = 0
			}
			out.SetGray(x, y, color.Gray{Y: new})

			err := int(old) - int(new)

			if x+1 < b.Max.X {
				v := int(img.GrayAt(x+1, y).Y) + err*7/16
				img.SetGray(x+1, y, color.Gray{Y: clamp(float64(v))})
			}
			if y+1 < b.Max.Y {
				if x > b.Min.X {
					v := int(img.GrayAt(x-1, y+1).Y) + err*3/16
					img.SetGray(x-1, y+1, color.Gray{Y: clamp(float64(v))})
				}
				v := int(img.GrayAt(x, y+1).Y) + err*5/16
				img.SetGray(x, y+1, color.Gray{Y: clamp(float64(v))})
				if x+1 < b.Max.X {
					v := int(img.GrayAt(x+1, y+1).Y) + err*1/16
					img.SetGray(x+1, y+1, color.Gray{Y: clamp(float64(v))})
				}
			}
		}
	}

	return out
}
