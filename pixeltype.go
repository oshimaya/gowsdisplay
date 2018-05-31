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
	RGBA(c color.Color, mask RGBmask)
}


// 32bit per pixel, ex RGBA(8:8:8:8)
type PIXEL32 [4]uint8

func NewRGB32(c color.Color, mask RGBmask) (p PIXEL32) {
	p.RGBA(c, mask)
	return
}

// Set 32bit-color data
func (p *PIXEL32) RGBA(c color.Color, mask RGBmask) {
	r,g,b,a := c.RGBA()
	r >>= 8
	g >>= 8
	b >>= 8
	a >>= 8
	d := (r<<mask.Red_size-1)/255<<mask.Red_offset |
		(g<<mask.Green_size-1)/255<<mask.Green_offset |
		(b<<mask.Blue_size-1)/255<<mask.Blue_offset
		//
		// Probably alpha bit is not used by any fb driver now because
	// alpha_offset and alpha_size is always 0 in wsdisplayio_get_fbinfo()
	// in sys/dev/wscons/wsdisplay_util.c.
	// Howerver check it heare for sure.
	//
	if mask.Alpha_size > 0 {
		d |= (a<<mask.Alpha_size - 1) / 255 << mask.Alpha_offset
	}
	for i := range p {
		p[i] = (*PIXEL32)(unsafe.Pointer(&d))[i]
	}
}

// 24bit per pixel, ex RGB(8:8:8)
type PIXEL24 [3]uint8

func NewRGB24(c color.Color, mask RGBmask) (p PIXEL24) {
	p.RGBA(c, mask)
	return
}
//
// Set 24bit-color data
//
func (p *PIXEL24) RGBA(c color.Color, mask RGBmask) {
	r,g,b,a := c.RGBA()
	r >>= 8
	g >>= 8
	b >>= 8
	a >>= 8
	d := (r<<mask.Red_size-1)/255<<mask.Red_offset |
		(g<<mask.Green_size-1)/255<<mask.Green_offset |
		(b<<mask.Blue_size-1)/255<<mask.Blue_offset
	//
	// all-mask for checking the valid byte
	//
	m := (1<<mask.Red_size-1)<<mask.Red_offset |
		(1<<mask.Green_size-1)<<mask.Green_offset |
		(1<<mask.Blue_size-1)<<mask.Blue_offset
		// maybe alpha bit is nothing but check it for safe
	if mask.Alpha_size > 0 {
		d |= (a<<mask.Alpha_size-1)/255 << mask.Alpha_offset
		m |= (1<<mask.Alpha_size - 1) << mask.Alpha_offset
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

func NewRGB16(c color.Color, mask RGBmask) (p PIXEL16) {
	p.RGBA(c, mask)
	return
}

//
// Set 16bit color Data, ex RGB=565, RGBA=5551
//
func (p *PIXEL16) RGBA(c color.Color, mask RGBmask) {
	r,g,b,a := c.RGBA()
	r >>= 8
	g >>= 8
	b >>= 8
	a >>= 8
	d := (r<<mask.Red_size-1)/255<<mask.Red_offset |
		(g<<mask.Green_size-1)/255<<mask.Green_offset |
		(b<<mask.Blue_size-1)/255<<mask.Blue_offset
	if mask.Alpha_size > 0 {
		d |= (a<<mask.Alpha_size-1)/255 << mask.Alpha_offset
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

func (p *PIXEL8) RGBA(c color.Color, mask RGBmask) {
	r,g,b,a := c.RGBA()
	r >>= 8
	g >>= 8
	b >>= 8
	a >>= 8
	p[0] = uint8((r * 299 + g* 587 + b*114)/1000)
}

type PIXELARRAY interface {
	StoreImage(src image.Image, mask RGBmask)
}

type PIXEL32ARRAY struct {
	Width int
	Height int
	Pix []PIXEL32
}

type PIXEL24ARRAY struct {
	Width int
	Height int
	Pix []PIXEL24
}

type PIXEL16ARRAY struct {
	Width int
	Height int
	Pix []PIXEL16
}

type PIXEL8ARRAY struct {
	Width int
	Height int
	Pix []PIXEL8
}

func (p *PIXEL32ARRAY) StoreImage (img image.Image, mask RGBmask) {

	w:=img.Bounds().Max.X-img.Bounds().Min.X
	h:=img.Bounds().Max.Y-img.Bounds().Min.Y

	p.Width = w
	p.Height = h
	p.Pix = make([]PIXEL32, w*h)
	for y:=img.Bounds().Min.Y; y<img.Bounds().Max.Y; y++ {
		for x:=img.Bounds().Min.X; x<img.Bounds().Max.X; x++ {
			p.Pix[x+y*w].RGBA(img.At(x,y), mask)
		}
	}
}

func (p *PIXEL24ARRAY) StoreImage (img image.Image, mask RGBmask) {

	w:=img.Bounds().Max.X-img.Bounds().Min.X
	h:=img.Bounds().Max.Y-img.Bounds().Min.Y

	p.Width = w
	p.Height = h
	p.Pix = make([]PIXEL24, w*h)
	for y:=img.Bounds().Min.Y; y<img.Bounds().Max.Y; y++ {
		for x:=img.Bounds().Min.X; x<img.Bounds().Max.X; x++ {
			p.Pix[x+y*w].RGBA(img.At(x,y), mask)
		}
	}
}
func (p *PIXEL16ARRAY) StoreImage (img image.Image, mask RGBmask) {

	w:=img.Bounds().Max.X-img.Bounds().Min.X
	h:=img.Bounds().Max.Y-img.Bounds().Min.Y

	p.Width = w
	p.Height = h
	p.Pix = make([]PIXEL16, w*h)
	for y:=img.Bounds().Min.Y; y<img.Bounds().Max.Y; y++ {
		for x:=img.Bounds().Min.X; x<img.Bounds().Max.X; x++ {
			p.Pix[x+y*w].RGBA(img.At(x,y), mask)
		}
	}
}

func (p *PIXEL8ARRAY) StoreImage (img image.Image, mask RGBmask) {

	w:=img.Bounds().Max.X-img.Bounds().Min.X
	h:=img.Bounds().Max.Y-img.Bounds().Min.Y

	p.Width = w
	p.Height = h
	p.Pix = make([]PIXEL8, w*h)
	for y:=img.Bounds().Min.Y; y<img.Bounds().Max.Y; y++ {
		for x:=img.Bounds().Min.X; x<img.Bounds().Max.X; x++ {
			p.Pix[x+y*w].RGBA(img.At(x,y), mask)
		}
	}
}
