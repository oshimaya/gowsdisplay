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

// +build netbsd

// wsdraw - GoWsDisplay example
//
// Usage:
//    wsdraw
//

package main

import (
	"fmt"
	"github.com/oshimaya/gowsdisplay"
	"image"
	"image/color"
	"time"
)

func main() {
	wsd := gowsdisplay.NewWsDisplay("/dev/ttyE1")
	err := wsd.Open()
	if err != nil {
		fmt.Println("Open: ", err)
		return
	}
	err = wsd.InitGraphics()
	if err != nil {
		fmt.Println("Initialize: ", err)
		return
	}

	defer wsd.Close()

	if wsd.GetPixelType() != gowsdisplay.FBRGB {
		fmt.Println("Sorry only support RGB type framebuffer.")
		return
	}

	wsd.Clear()

	w := wsd.GetWidth()
	h := wsd.GetHeight()
	p := wsd.NewPixel(color.RGBA{255, 255, 255, 255})

	for i := 0; i < h; i++ {
		wsd.SetPixel(i, i, p)
	}
	for i := 0; i < h/2; i += 5 {
		b := image.Rect(w/2-i, h/2-i, w/2+i, h/2+i)
		wsd.DrawBox(b, p)
	}
	for i := 0; i < 256; i++ {
		b := image.Rect(0, 0, 256-i, 256-i)
		c := color.RGBA{uint8(i), uint8(i), uint8(i), 255}
		p.SetColor(c, wsd.GetRGBmask())
		wsd.FillBox(b, p)
	}

	for i := 0; i < h; i += 8 {
		c := color.RGBA{uint8(i%256) / 4, uint8(i%256) / 2, uint8(i % 256), 255}
		p.SetColor(c, wsd.GetRGBmask())
		wsd.DrawLine(image.Pt(i, 0), image.Pt(799-i, 599), p)
	}

	c := color.RGBA{255, 0, 0, 255}
	p.SetColor(c, wsd.GetRGBmask())
	wsd.DrawCircle(w/2, h/2, 150, p)
	wsd.DrawCircle(w/2, h/2, 100, p)
	c = color.RGBA{0, 255, 0, 255}
	p.SetColor(c, wsd.GetRGBmask())
	wsd.FillCircle(w/2, h/2, 50, p)

	time.Sleep(5 * time.Second)
}
