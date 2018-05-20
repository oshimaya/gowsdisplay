// Copyright (c) 2018 Yasushi Oshima All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions
// are met:
// 1. Redistributions of source code must retain the above copyright
//    notice, this list of conditions and the following disclaimer.
// 2. Redistributions in binary form must reproduce the above copyright
//    notice, this list of conditions and the following disclaimer in the
//
//    documentation and/or other materials provided with the distribution.
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

// WsDisplay manager; wsdisplay(4) wrapper for golang
//   from dev/wscons/wsconsio.h

// +build netbsd

package gowsdisplay

import (
        "syscall"
        "unsafe"
	"errors"
)

// ioctl number for wsdisplayio
const (
        FBGTYPE           = 0x40045740
        FBGINFO           = 0x40105741
        FBGETCMAP         = 0x80205742
        FBPUTCMAP         = 0x80205743
        FBGVIDEO          = 0x40045744
        FBSVIDEO          = 0x80045745
        FBGCURPOS         = 0x40085746
        FBSCURPOS         = 0x80085747
        FBGCURMAX         = 0x40085748
        FBGCURSOR         = 0xc0505749
        FBSCURSOR         = 0x8050574a
        FBGMODE           = 0x4004574b
        FBSMODE           = 0x8004574c
        FBLDFONT          = 0x8030574d
        FBADDSCREEN       = 0x8018574e
        FBDELSCREEN       = 0x8008574f
        FBSFONT           = 0x80085750
        FBSETKEYBOARD     = 0xc0085751
        FBGETPARAM        = 0xc0205752
        FBSETPARAM        = 0xc0205753
        FBGETACTIVESCREEN = 0x40045754
        FBGETWSCHAR       = 0xc0105755
        FBPUTWSCHAR       = 0xc0105756
        FBDGSCROLL        = 0x400c5757
        FBDSSCROLL        = 0x800c5758
        FBGMSGATTRS       = 0x40185759
        FBSMSGATTRS       = 0x8018575a
        FBGBORDER         = 0x4004575b
        FBSBORDER         = 0x8004575c
        FBSSPLASH         = 0x8004575d
        FBSPROGRESS       = 0x8004575e
        FBLINEBYTES       = 0x4004575f
        FBSETVERSION      = 0x80045760
        FBGETBUSID        = 0x40245765
        FBGETEDID         = 0xc0105766
        FBSETPOLLING      = 0x80045767
        FBGETFBINFO       = 0xc0485768
        FBDOBLIT          = 0xc0245769
        FBWAITBLIT        = 0xc024576a
)


// for WSDISPLAYIO_[GS]MODE 
const (
        FBMODE_EMUL   = 0	// text emulation
        FBMODE_MAPPED = 1	// mapped graphics
        FBMODE_DUMBFB = 2	// mapped graphics frambuffer 
)

// Display Type
const (
        FBRGB       = 0		// RGB color
        FBCI        = 1		// indexed color
        FBGREYSCALE = 2		// grayscale
        FBYUV       = 3		// YUV color
)

// FBI type
const (
        FBVRAM_IS_RAM = 1	// not shadow
        FBVRAM_IS_SPLIT = 2	// 
)


type WsDisplay struct {
	fd	int
	info	FBinfo		// fbinfo struct 
	addr	[]byte		// display memory (VRAM) address for mmap
	dev	string		// device name
}

type RGBmask struct {
        Red_offset   uint32	// Red offset bits from the right
        Red_size     uint32	// Red size in bits
        Green_offset uint32	// Green offset bits from the right
        Green_size   uint32	// Green size in bits
        Blue_offset  uint32	// Blue offset bits from the right
        Blue_size    uint32	// Blue size in bits
        Alpha_offset uint32	// Alpha offset bits from the right
        Alpha_size   uint32	// Alpha size in bits
}

// struct wsdisplay_fbinfo
// XXX: how to define for non RGB type?
type FBinfo struct {
	Size uint64	// fb size in bytes
	Offset uint64	// start of visible fb in bytes
	Width uint32	// screen width in pixels
	Height uint32	// screen height in lines
	Stride uint32	// stride, bytes of one line 
	Bitsperpixel uint32	// size of one pixel in bits 
	Pixeltype uint32	// pixel type  (RGB/CI/GRAY..)
	Rgbmask RGBmask		// RGB masks for RGB type
	Flags uint32		// flags
}

// Create New Display
//   ex. NewWsDisplay("/dev/ttyE2")
func NewWsDisplay(dev string) *WsDisplay {
	wsd := new(WsDisplay)
	wsd.dev = dev
	return wsd
}


// Open display and get fbinfo 
func (wsd *WsDisplay) Open() (err error) {
	
	wsd.fd, err = syscall.Open(wsd.dev, syscall.O_RDWR, 0)
	if err != nil {
		return err
	}
	err = wsd.getFBinfo()
	
	return  err 
}

// Close display and set to text emul mode
func (wsd *WsDisplay) Close() error {
	if wsd.addr != nil {
		wsd.unmapFB()
	}
	wsd.setMode(FBMODE_EMUL)
	return syscall.Close(wsd.fd)
}


func (wsd *WsDisplay) GetFBinfo() FBinfo {
	return wsd.info
}

func (wsd *WsDisplay) GetRGBmask() RGBmask {
	return wsd.info.Rgbmask
}

func (wsd *WsDisplay) GetBufferAddr() *byte {
	return &wsd.addr[0]
}

func (wsd *WsDisplay) getFBinfo() error {
	_,_, errno := syscall.Syscall(syscall.SYS_IOCTL, uintptr(wsd.fd),
		FBGETFBINFO, uintptr(unsafe.Pointer(&wsd.info)))
	if errno == 0 {
		return nil
	}
	return errno
}

func (wsd *WsDisplay) mapFB() (err error) {
	if wsd.info.Size > 0 {
		len := int(wsd.info.Size)
		wsd.addr, err = syscall.Mmap(wsd.fd, 0, len,
			syscall.PROT_READ|syscall.PROT_WRITE,0)
		return err
	}
	err = errors.New("Maybe the framebuffer is uninitialized")
	return err
}

func (wsd *WsDisplay) unmapFB() error {
	return syscall.Munmap(wsd.addr)
}

func (wsd *WsDisplay) setMode(mode int) error {
	_,_, errno := syscall.Syscall(syscall.SYS_IOCTL, uintptr(wsd.fd),
		FBSMODE, uintptr(unsafe.Pointer(&mode)))
	if errno == 0 {
		return nil
	}
	return errno
}


// set to dumbfb mode and mmap framebuffer
func (wsd *WsDisplay) InitGraphics() error {
	err :=  wsd.setMode(FBMODE_DUMBFB)
	if err != nil {
		return err
	}
	err = wsd.mapFB()
	return err
}

