# WsDisplay for Golang

## Whatis

NetBSD wsdisplay(4) wrapper for Golang.

## Basic use

- Open device


```go
	wsd := NewWsDisplay("/dev/ttyE1")
	wsd.Open()
```

- Init and set to framebuffer mode 

```go
	wsd.InitGraphics()
```

- Check info fb's depth, size, type, and so on.

```go
	...
	info := wsd.GetFBinfo()
	depth := int(info.Bitsperpixel) / 8
	view_width := int(info.Width)
	view_height := int(info.Height)
	stride := int(info.Stride)
	...
```

- Data type

Currently, only RGBtype is supported.

- Access to framebuffer memory


## Examples

  T.B.D.

## License

 BSD

## Author
 Yasushi Oshima

