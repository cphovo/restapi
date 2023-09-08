package main

import (
	"bytes"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/math/fixed"
)

// HexToColor converts a 6-digit hexadecimal string to a color.RGBA value.
// The hexadecimal string can optionally start with a '#'.
// It returns an error if the input string is not a valid 6-digit hexadecimal color code.
func HexToColor(hex string) (color.RGBA, error) {
	if hex[0] == '#' {
		hex = hex[1:]
	}

	// Handle short 3-character format by duplicating each character
	if len(hex) == 3 {
		hex = string([]byte{
			hex[0], hex[0],
			hex[1], hex[1],
			hex[2], hex[2],
		})
	}

	rgb, err := strconv.ParseUint(hex, 16, 24)
	if err != nil {
		return color.RGBA{}, err
	}

	return color.RGBA{
		R: uint8((rgb & 0xFF0000) >> 16),
		G: uint8((rgb & 0x00FF00) >> 8),
		B: uint8((rgb & 0x0000FF) >> 0),
		A: 255,
	}, nil
}

func addLabel(img *image.RGBA, y int, label string, fontSize float64) {
	// Load the truetype font
	tt, err := truetype.Parse(goregular.TTF)
	if err != nil {
		panic(err)
	}

	// Create the font face with desired font size
	face := truetype.NewFace(tt, &truetype.Options{
		Size: fontSize,
	})

	// Measure the text's width
	textWidth := font.MeasureString(face, label).Round()

	// Calculate x to center the text
	x := (img.Rect.Max.X - textWidth) / 2

	// Calculate y to center the text vertically
	metrics := face.Metrics()
	ascent := metrics.Ascent.Round()
	descent := metrics.Descent.Round()
	lineHeight := ascent + descent
	y = y + ascent - lineHeight/2

	point := fixed.Point26_6{X: fixed.Int26_6(x * 64), Y: fixed.Int26_6(y * 64)}

	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(color.RGBA{60, 60, 60, 255}),
		Face: face,
		Dot:  point,
	}
	d.DrawString(label)
}

func GenerateImage(width, height int, hex string) ([]byte, error) {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	c, err := HexToColor(hex)
	if err != nil {
		return nil, err
	}

	draw.Draw(img, img.Bounds(), &image.Uniform{c}, image.Point{}, draw.Src)

	// Add label at the center
	label := strconv.Itoa(width) + " x " + strconv.Itoa(height)
	fontSize := float64(width) / 10.0 // Font size is 1/10th of image width
	y := height / 2                   // Center the label vertically
	addLabel(img, y, label, fontSize)

	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func main() {
	// GenerateImage(600, 600, "#FF5733")
	r := gin.Default()
	r.GET("/", func(ctx *gin.Context) { ctx.JSON(200, gin.H{"message": "Everything is ok!"}) })
	r.GET("/photos/:size/:color", func(ctx *gin.Context) {
		size, err := strconv.Atoi(ctx.Param("size"))
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid size parameter"})
			return
		}
		color := ctx.Param("color")
		imgBytes, err := GenerateImage(size, size, color)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid color parameter"})
			return
		}
		ctx.Data(http.StatusOK, "image/png", imgBytes)
	})
	r.Run(":8000")
}
