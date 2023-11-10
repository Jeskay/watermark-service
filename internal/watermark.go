package internal

import (
	"image"
	"image/color"

	"golang.org/x/image/draw"
	"golang.org/x/image/math/fixed"

	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
)

type Position string

const (
	LeftTop     Position = "left_top"
	LeftBottom  Position = "left_bottom"
	RightTop    Position = "right_top"
	RightBottom Position = "right_bottom"
)

func PositionFromString(text string) Position {
	switch text {
	case "left_top":
		return LeftTop
	case "left_bottom":
		return LeftBottom
	case "right_top":
		return RightTop
	case "right_bottom":
		return RightBottom
	}
	return LeftTop
}

func CombineTextWithLogo(logo image.Image, text string) image.Image {
	if text != "" {
		var logo_new *image.RGBA
		var text_x, text_y fixed.Int26_6

		space_between := 20

		text_width := font.MeasureString(basicfont.Face7x13, text).Ceil()
		text_height := basicfont.Face7x13.Metrics().Ascent.Ceil()

		if logo != nil {
			logo_rect := logo.Bounds()
			//scaling logo
			adjacent_h := text_height * 4
			multiplier := float64(adjacent_h) / float64(logo_rect.Dy())
			logo_img := image.NewRGBA(image.Rect(0, 0, int(float64(logo_rect.Dx())*multiplier), adjacent_h))
			draw.ApproxBiLinear.Scale(logo_img, logo_img.Rect, logo, logo_rect, draw.Over, nil)
			//adding extra space for text
			logo_new = image.NewRGBA(image.Rect(0, 0, logo_img.Rect.Dx()+space_between+text_width, logo_img.Rect.Dy()))
			draw.Draw(logo_new, logo_img.Rect, logo_img, image.Point{0, 0}, draw.Over)
			text_x, text_y = fixed.I(logo_img.Rect.Dx()), fixed.I((logo_img.Rect.Dy()+text_height)/2)
		} else {
			logo_new = image.NewRGBA(image.Rect(0, 0, text_width+space_between, text_height+space_between))
			text_x, text_y = fixed.I(0), fixed.I(text_height)
		}
		// inserting text
		col := color.RGBA{200, 100, 0, 255}
		point := fixed.Point26_6{X: text_x, Y: text_y}
		d := font.Drawer{
			Dst:  logo_new,
			Src:  image.NewUniform(col),
			Face: basicfont.Face7x13,
			Dot:  point,
		}
		d.DrawString(text)
		return logo_new
	} else if logo != nil {
		return logo
	} else {
		return nil
	}
}

func AddWatermarkToImage(watermark image.Image, src image.Image, pos Position) draw.Image {
	src_rect := src.Bounds()
	wtm_rect := watermark.Bounds()

	var offset image.Point
	switch pos {
	case LeftTop:
		offset = image.Pt(0, 0)
	case RightTop:
		offset = image.Pt(src_rect.Dx()-wtm_rect.Bounds().Dx(), 0)
	case LeftBottom:
		offset = image.Pt(0, src_rect.Dy()-wtm_rect.Bounds().Dy())
	case RightBottom:
		offset = image.Pt(src_rect.Dx()-wtm_rect.Bounds().Dx(), src_rect.Dy()-wtm_rect.Bounds().Dy())
	default:
		offset = image.Pt(0, 0)
	}

	bg := image.NewRGBA(image.Rect(0, 0, src_rect.Dx(), src_rect.Dy()))
	draw.Draw(bg, src_rect, src, image.Point{0, 0}, draw.Over)
	//applying opacity mask to watermark
	mask := image.NewUniform(color.Alpha{96})
	draw.DrawMask(bg, src_rect.Add(offset), watermark, image.Point{0, 0}, mask, image.Point{0, 0}, draw.Over)
	return bg
}
