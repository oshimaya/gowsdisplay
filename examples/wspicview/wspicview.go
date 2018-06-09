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

// wspicview - GoWsDisplay example
//
// Usage: 
//    wspicvew {picturefile ...}
//

package main

import (
	"fmt"
	"github.com/oshimaya/gowsdisplay"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"
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
	for i, argv := range os.Args {
		if i == 0 {
			continue
		}
		img, err := loadPic(argv)
		if err != nil {
			fmt.Println("Image load: ", err)
			continue
		}

		// Create Pixel Array (for pattern data)

		p, err := wsd.NewPixelArray()

		// Store image data to Pixel Array

		p.StoreImage(img, wsd.GetRGBmask())

		// Display Pixel Array to wsdisplay

		wsd.Clear()
		wsd.PutPixelArray(0, 0, p)

		time.Sleep(time.Second * 5)
	}
}

func loadPic(fname string) (image.Image, error) {
	f, err := os.Open(fname)
	if err != nil {
		fmt.Println("Open:", err)
		return nil, err
	}
	defer f.Close()
	img, _, err := image.Decode(f)
	return img, err
}
