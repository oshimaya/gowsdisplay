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

### Check info fb's depth, size, type, and so on.

```go
	...
	info := wsd.GetFBinfo()
	depth := int(info.Bitsperpixel) / 8
	view_width := int(info.Width)
	view_height := int(info.Height)
	stride := int(info.Stride)
	...
```

Maybe some accessor function will be added in the future...

### Data type
- 1pixel format
 - PIXEL32  
1pixel = 32bit, [4]uint8. Typically RGBA8:8:8:8, but the order or bit format is not specified this.
 - PIXEL24  
1pixel = 24bit, [3]uint8. Typically RGB8:8:8, but the order or bit format is not specified this.
 - PIXEL16  
1pixel = 16bit, [2]uint8. Typically RGB5:6/5:5 or YUV422, but the order or bit format is not specified this.
 - PIXEL8  
1pixel = 8bit, [1]uint8. Typically Gray or Indexed color.

- Apply the mask and order to Color value  
T.B.D.

### Access to framebuffer memory

Minimum:

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

  T.B.D.

## License

 BSD

## Author
 Yasushi Oshima

