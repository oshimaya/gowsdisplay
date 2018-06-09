// Copyright (c) 2018 Yasushi Oshima All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions
// are met:
// 1. Redistributions of source code must retain the above copyright
//    notice, this list of conditions and the following disclaimer.
// 2. Redistributions in binary form must reproduce the above copyright
//    notice, this list of conditions and the following disclaimer in the
//    documentation and/or other materials provided with the distribution.
//
// THIS SOFTWARE IS PROVIDED BY THE AUTHOR AND CONTRIBUTORS ``AS IS'' AND
// ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
// IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
// ARE DISCLAIMED.  IN NO EVENT SHALL THE AUTHOR OR CONTRIBUTORS BE LIABLE
// FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
// DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS
// OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION)
// HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT
// LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY
// OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF
// SUCH DAMAGE.

// Pixel data types for accessing to WsDisplay frame buffer

// +build netbsd

package gowsdisplay

import (
	"image"
	"image/color"
	"unsafe"
)

// Interface depth independent

type PIXEL interface {
	SetColor(c color.Color, rgbmask RGBmask)
}

// 32bit per pixel, ex RGBA(8:8:8:8)
type PIXEL32 [4]uint8

func NewRGB32(c color.Color, rgbmask RGBmask) (p PIXEL32) {
	p.SetColor(c, rgbmask)
	return
}

// Set 32bit-color data
func (p *PIXEL32) SetColor(c color.Color, rgbmask RGBmask) {
	r, g, b, a := c.RGBA()
	r >>= 8
	g >>= 8
	b >>= 8
	a >>= 8
	d := (r<<rgbmask.Red_size-1)/255<<rgbmask.Red_offset |
		(g<<rgbmask.Green_size-1)/255<<rgbmask.Green_offset |
		(b<<rgbmask.Blue_size-1)/255<<rgbmask.Blue_offset
		//
		// Probably alpha bit is not used by any fb driver now because
	// alpha_offset and alpha_size is always 0 in wsdisplayio_get_fbinfo()
	// in sys/dev/wscons/wsdisplay_util.c.
	// Howerver check it heare for sure.
	//
	if rgbmask.Alpha_size > 0 {
		d |= (a<<rgbmask.Alpha_size - 1) / 255 << rgbmask.Alpha_offset
	}
	for i := range p {
		p[i] = (*PIXEL32)(unsafe.Pointer(&d))[i]
	}
}

// 24bit per pixel, ex RGB(8:8:8)
type PIXEL24 [3]uint8

func NewRGB24(c color.Color, rgbmask RGBmask) (p PIXEL24) {
	p.SetColor(c, rgbmask)
	return
}

//
// Set 24bit-color data
//
func (p *PIXEL24) SetColor(c color.Color, rgbmask RGBmask) {
	r, g, b, a := c.RGBA()
	r >>= 8
	g >>= 8
	b >>= 8
	a >>= 8
	d := (r<<rgbmask.Red_size-1)/255<<rgbmask.Red_offset |
		(g<<rgbmask.Green_size-1)/255<<rgbmask.Green_offset |
		(b<<rgbmask.Blue_size-1)/255<<rgbmask.Blue_offset
	//
	// all-rgbmask for checking the valid byte
	//
	m := (1<<rgbmask.Red_size-1)<<rgbmask.Red_offset |
		(1<<rgbmask.Green_size-1)<<rgbmask.Green_offset |
		(1<<rgbmask.Blue_size-1)<<rgbmask.Blue_offset
		// maybe alpha bit is nothing but check it for safe
	if rgbmask.Alpha_size > 0 {
		d |= (a<<rgbmask.Alpha_size - 1) / 255 << rgbmask.Alpha_offset
		m |= (1<<rgbmask.Alpha_size - 1) << rgbmask.Alpha_offset
	}
	mp := (*PIXEL24)(unsafe.Pointer(&m))
	// convert to uint8 data in PIXEL
	j := 0
	for i := 0; i < len(mp); i++ {
		if mp[i] != 0 {
			p[j] = (*PIXEL24)(unsafe.Pointer(&d))[i]
			j++
		}
	}
}

// 16bit per pixel, ex RGB(5:6:6) or YUV(4:2:2)
type PIXEL16 [2]uint8

func NewRGB16(c color.Color, rgbmask RGBmask) (p PIXEL16) {
	p.SetColor(c, rgbmask)
	return
}

//
// Set 16bit color Data, ex RGB=565, RGBA=5551
//
func (p *PIXEL16) SetColor(c color.Color, rgbmask RGBmask) {
	r, g, b, a := c.RGBA()
	r >>= 8
	g >>= 8
	b >>= 8
	a >>= 8
	d := (r<<rgbmask.Red_size-1)/255<<rgbmask.Red_offset |
		(g<<rgbmask.Green_size-1)/255<<rgbmask.Green_offset |
		(b<<rgbmask.Blue_size-1)/255<<rgbmask.Blue_offset
	if rgbmask.Alpha_size > 0 {
		d |= (a<<rgbmask.Alpha_size - 1) / 255 << rgbmask.Alpha_offset
	}
	//
	// convert to int8 data in PIXEL
	//
	for i := range p {
		p[i] = (*PIXEL16)(unsafe.Pointer(&d))[i]
	}
}

// 8bit per pixel, ex Gray8 or Color Indexed
type PIXEL8 [1]uint8

func (p *PIXEL8) SetColor(c color.Color, mask RGBmask) {
	r, g, b, a := c.RGBA()
	r >>= 8
	g >>= 8
	b >>= 8
	a >>= 8
	p[0] = uint8((r*299 + g*587 + b*114) / 1000)
}

type PIXELARRAY interface {
	StoreImage(src image.Image, rgbmask RGBmask)
	GetWidth() int
	GetHeight() int
	PutPixelPat(x int, y int, pix PIXELARRAY)
}

type pixelarray struct {
	width  int
	height int
	mask   []bool
}

func (p *pixelarray) GetWidth() int {
	return p.width
}

func (p *pixelarray) GetHeight() int {
	return p.height
}

type PIXEL32ARRAY struct {
	pixelarray
	pix []PIXEL32
}

type PIXEL24ARRAY struct {
	pixelarray
	pix []PIXEL24
}

type PIXEL16ARRAY struct {
	pixelarray
	pix []PIXEL16
}

type PIXEL8ARRAY struct {
	pixelarray
	pix []PIXEL8
}

func (p *PIXEL32ARRAY) StoreImage(img image.Image, rgbmask RGBmask) {

	w := img.Bounds().Max.X - img.Bounds().Min.X
	h := img.Bounds().Max.Y - img.Bounds().Min.Y
	min_x := img.Bounds().Min.X
	min_y := img.Bounds().Min.Y

	p.width = w
	p.height = h
	p.pix = make([]PIXEL32, w*h)
	p.mask = make([]bool, w*h)
	for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ {
		for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
			c := img.At(x, y)
			p.pix[x-min_x+(y-min_y)*w].SetColor(c, rgbmask)
			_, _, _, a := c.RGBA()
			if a > 0 {
				p.mask[x-min_x+(y-min_y)*w] = true
			} else {
				p.mask[x-min_x+(y-min_y)*w] = false
			}
		}
	}
}

func (p *PIXEL32ARRAY) PutPixelPat(dest_x int, dest_y int, pix PIXELARRAY) {
	switch pix.(type) {
	case *PIXEL32ARRAY:
		src := pix.(*PIXEL32ARRAY)
		for y := 0; y < src.height; y++ {
			for x := 0; x < src.width; x++ {
				if dest_x < 0 || dest_x >= p.width ||
					dest_y < 0 || dest_y >= p.height ||
					!src.mask[x+y*src.width] {
					return
				}
				p.pix[dest_x+x+(dest_y+y)*p.width] =
					src.pix[x+y*src.width]
				p.mask[dest_x+x+(dest_y+y)*p.width] = true
			}
		}
	}
}

func (p *PIXEL24ARRAY) StoreImage(img image.Image, rgbmask RGBmask) {

	w := img.Bounds().Max.X - img.Bounds().Min.X
	h := img.Bounds().Max.Y - img.Bounds().Min.Y
	min_x := img.Bounds().Min.X
	min_y := img.Bounds().Min.Y

	p.width = w
	p.height = h
	p.pix = make([]PIXEL24, w*h)
	for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ {
		for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
			c := img.At(x, y)
			p.pix[x-min_x+(y-min_y)*w].SetColor(c, rgbmask)
			_, _, _, a := c.RGBA()
			if a > 0 {
				p.mask[x-min_x+(y-min_y)*w] = true
			} else {
				p.mask[x-min_x+(y-min_y)*w] = false
			}
		}
	}
}

func (p *PIXEL24ARRAY) PutPixelPat(dest_x int, dest_y int, pix PIXELARRAY) {
	switch pix.(type) {
	case *PIXEL24ARRAY:
		src := pix.(*PIXEL24ARRAY)
		for y := 0; y < src.height; y++ {
			for x := 0; x < src.width; x++ {
				if dest_x < 0 || dest_x >= p.width ||
					dest_y < 0 || dest_y >= p.height ||
					!src.mask[x+y*src.width] {
					return
				}
				p.pix[dest_x+x+(dest_y+y)*p.width] =
					src.pix[x+y*src.width]
				p.mask[dest_x+x+(dest_y+y)*p.width] = true
			}
		}
	}
}

func (p *PIXEL16ARRAY) StoreImage(img image.Image, rgbmask RGBmask) {

	w := img.Bounds().Max.X - img.Bounds().Min.X
	h := img.Bounds().Max.Y - img.Bounds().Min.Y
	min_x := img.Bounds().Min.X
	min_y := img.Bounds().Min.Y

	p.width = w
	p.height = h
	p.pix = make([]PIXEL16, w*h)
	for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ {
		for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
			c := img.At(x, y)
			p.pix[x-min_x+(y-min_y)*w].SetColor(c, rgbmask)
			_, _, _, a := c.RGBA()
			if a > 0 {
				p.mask[x-min_x+(y-min_y)*w] = true
			} else {
				p.mask[x-min_x+(y-min_y)*w] = false
			}
		}
	}
}

func (p *PIXEL16ARRAY) PutPixelPat(dest_x int, dest_y int, pix PIXELARRAY) {
	switch pix.(type) {
	case *PIXEL16ARRAY:
		src := pix.(*PIXEL16ARRAY)
		for y := 0; y < src.height; y++ {
			for x := 0; x < src.width; x++ {
				if dest_x < 0 || dest_x >= p.width ||
					dest_y < 0 || dest_y >= p.height ||
					!src.mask[x+y*src.width] {
					return
				}
				p.pix[dest_x+x+(dest_y+y)*p.width] =
					src.pix[x+y*src.width]
				p.mask[dest_x+x+(dest_y+y)*p.width] = true
			}
		}
	}
}

func (p *PIXEL8ARRAY) StoreImage(img image.Image, rgbmask RGBmask) {

	w := img.Bounds().Max.X - img.Bounds().Min.X
	h := img.Bounds().Max.Y - img.Bounds().Min.Y
	min_x := img.Bounds().Min.X
	min_y := img.Bounds().Min.Y

	p.width = w
	p.height = h
	p.pix = make([]PIXEL8, w*h)
	for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ {
		for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
			c := img.At(x, y)
			p.pix[x-min_x+(y-min_y)*w].SetColor(c, rgbmask)
			_, _, _, a := c.RGBA()
			if a > 0 {
				p.mask[x-min_x+(y-min_y)*w] = true
			} else {
				p.mask[x-min_x+(y-min_y)*w] = false
			}
		}
	}
}

func (p *PIXEL8ARRAY) PutPixelPat(dest_x int, dest_y int, pix PIXELARRAY) {
	switch pix.(type) {
	case *PIXEL8ARRAY:
		src := pix.(*PIXEL8ARRAY)
		for y := 0; y < src.height; y++ {
			for x := 0; x < src.width; x++ {
				if dest_x < 0 || dest_x >= p.width ||
					dest_y < 0 || dest_y >= p.height ||
					!src.mask[x+y*src.width] {
					return
				}
				p.pix[dest_x+x+(dest_y+y)*p.width] =
					src.pix[x+y*src.width]
				p.mask[dest_x+x+(dest_y+y)*p.width] = true
			}
		}
	}
}
