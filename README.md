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
interface for PIXEL*ARRAY, for Image operation
- PIXEL32ARRAY  
- PIXEL24ARRAY  
- PIXEL16ARRAY  
### Access to framebuffer memory

#### Minimum:

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
##### ImageOperation:

```go
	// Create PixelAarray
	p, err := wsd.NewPixelArray()
	// Convert and set image data to PixelArray
	p.StoreImage(img, wsd.GetRGBmask())

	// Draw image to wsdisplay framebuffer at (x,y)
	wsd.PutPixelArray(x,y, p)
```

### Terminate  

```go
	wsd.Close()
```
Typically, call this in defer.

## Examples

- [examples/wspicview]


## License

 BSD

## Author
 Yasushi Oshima

