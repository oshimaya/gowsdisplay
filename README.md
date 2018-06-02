# WsDisplay for Golang

## Whatis

NetBSD wsdisplay(4) wrapper for Golang.

## Basic use

### Open device

```go
	wsd := gowsdisplay.NewWsDisplay("/dev/ttyE?")
	wsd.Open()
```
Note: Requires RW access permission to /dev/ttyE?

### Init and set to framebuffer mode 

```go
	wsd.InitGraphics()
```
Note: This changes ttyE? to framebuffer mode, it's better to operation from network remotely.

### Check fb's depth, size, type, and so on.

```go
	...
	depth := wsd.GetDepth()
	type := wsd.GetPixelType()
	view_width := wsd.GetWidth()
	view_height := wsd.GetHeight()
	stride := wsd.GetStride()
	pixel_stride := wsd.GetPixelStride()
	...
```

### Data type
1pixel format:
- PIXEL  
Interface for PIXEL*
  - SetColor(color.Color, RGBmask)  
Set/Convert color data to PIXEL format with specified mask
- PIXEL32  
1pixel = 32bit, [4]uint8. Typically RGBA8:8:8:8, but the order or bit format is not specified this.
- PIXEL24  
1pixel = 24bit, [3]uint8. Typically RGB8:8:8, but the order or bit format is not specified this.
- PIXEL16  
1pixel = 16bit, [2]uint8. Typically RGB5:6/5:5 or YUV422, but the order or bit format is not specified this.
- PIXEL8  
1pixel = 8bit, [1]uint8. Typically Gray or Indexed color.

PixelArray format:

- PIXELARRAY  
Interface for PIXEL*ARRAY, for Image operation
  - StoreImage(image.Image, RGBmask)  
Set image data from image.Image 
  - GetWidth()  
Get width of the image stored in this PIXELARRAY
 - GetHeight()  
Get height of the image stored in this PIXELARRAY
- PIXEL32ARRAY  
Array of PIXEL32
- PIXEL24ARRAY  
Array of PIXEL24
- PIXEL16ARRAY  
Array of PIXEL16
### Access to framebuffer memory

##### Drawing Operaion::

- Prepear PIXEL

```go
	var c color.Color
	...
	pix := wsd.NewPixel( c )
	...
```

- SetPixel

```go
	wsd.SetPixel(x, y, pix)
```
- DrawLine

```go
	var start, end image.Point 
	...
	wsd.DrawLine(start, end, pix)
```

- DrawBox, FillBox

```go
	var area image.Rectangle 
	...
	wsd.DrawBox(area, pix)
	wsd.FillBox(area, pix)
```

- DrawCircle, FillCircle

```go
	wsd.DrawCircle(x, y, r, pix)
	wsd.FillCircle(x, y, r, pix)
```
##### Draw Image 

```go
	// Create PixelAarray
	p, err := wsd.NewPixelArray()
	// Convert and set image data to PixelArray
	p.StoreImage(img, wsd.GetRGBmask())

	// Draw image to wsdisplay framebuffer at (x,y)
	wsd.PutPixelArray(x,y, p)
```

#### Set Raw Data:

```go
	pix := ws.GetBufferAsPixel32()	// Get framebuffer data as []PIXEL32 slice

	var rawdata PIXEL32
	...
	// prepar 1 pixel raw data
	rawdata[0] = ...
	rawdata[1] = ...
	rawdata[2] = ...
	rawdata[3] = ...
	// Write 1 pixel
	pix[index] = rawdata
```
### Terminate  

```go
	wsd.Close()
```
Typically, call this in defer.

## Examples

- examples/wsdraw  
Some drawing operation
- examples/wspcview  
View picture from JPEG/PNG file

## License

 BSD

## Author
 Yasushi Oshima

