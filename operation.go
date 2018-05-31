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

package gowsdisplay

import (
	"errors"
)

func (wsd *WsDisplay) PutPxielArray(px int, py int, p PIXELARRAY) error {
	screenW := int(wsd.GetWidth())
	screenH := int(wsd.GetHeight())
	w := p.GetWidth()
	h := p.GetHeight()
	if px+w < 0 || py+h < 0 || px > screenW || py > screenH {
		// Nothing to do. All area is out of screen
		return nil
	}
	startx := 0
	starty := 0
	endx := w
	endy := h
	if px < 0 {
		startx = (0 - px)
		endx = endx + px
	}
	if w+px > screenW {
		endx = screenW - px
	}
	if py < 0 {
		starty = (0 - py)
		endy = endy + py
	}
	if h+py > screenH {
		endy = screenH - py
	}
	s := int(wsd.GetPixelStride())
	switch p.(type) {
	case *PIXEL32ARRAY:
		if wsd.GetDepth() != 32 {
			err := errors.New("Unmatch PixelDepth")
			return err
		}
		q := p.(*PIXEL32ARRAY)
		pix := wsd.GetBufferAsPixel32()
		for srcy := starty; srcy < endy; srcy++ {
			copy(
				pix[startx+px+(srcy+py)*s:startx+px+(srcy+py)*s+endx],
				q.Pix[startx+srcy*w:startx+srcy*w+endx])
		}
	case *PIXEL24ARRAY:
		if wsd.GetDepth() != 24 {
			err := errors.New("Unmatch PixelDepth")
			return err
		}
		q := p.(*PIXEL24ARRAY)
		pix := wsd.GetBufferAsPixel24()
		for srcy := starty; srcy < endy; srcy++ {
			copy(
				pix[startx+px+(srcy+py)*s:startx+px+(srcy+py)*s+endx],
				q.Pix[startx+srcy*w:startx+srcy*w+endx])
		}
	case *PIXEL16ARRAY:
		if wsd.GetDepth() != 16 {
			err := errors.New("Unmatch PixelDepth")
			return err
		}
		q := p.(*PIXEL16ARRAY)
		pix := wsd.GetBufferAsPixel16()
		for srcy := starty; srcy < endy; srcy++ {
			copy(
				pix[startx+px+(srcy+py)*s:startx+px+(srcy+py)*s+endx],
				q.Pix[startx+srcy*w:startx+srcy*w+endx])
		}
	default:
		err := errors.New("Unsupported PixelType")
		return err
	}
	return nil
}

func (wsd *WsDisplay) NewPixelArray() (p PIXELARRAY, err error) {
	switch wsd.GetDepth() {
	case 32:
		p = new(PIXEL32ARRAY)
	case 24:
		p = new(PIXEL24ARRAY)
	case 16:
		p = new(PIXEL16ARRAY)
	default:
		return nil, errors.New("Unspport Display Depth")
	}
	return p, nil
}

func (wsd *WsDisplay) Clear() {
	switch wsd.GetDepth() {
	case 32:
		pix := wsd.GetBufferAsPixel32()
		for i := range pix {
			pix[i] = PIXEL32{0, 0, 0, 0}
		}
	case 24:
		pix := wsd.GetBufferAsPixel24()
		for i := range pix {
			pix[i] = PIXEL24{0, 0, 0}
		}
	case 16:
		pix := wsd.GetBufferAsPixel16()
		for i := range pix {
			pix[i] = PIXEL16{0, 0}
		}
	}
}
