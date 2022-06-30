# imgproxy-go

A Go client library to generate urls for imgproxy services.

Based on https://docs.imgproxy.net/

**Note:** Is not feature complete.

## Usage

```go
  ip, err := NewImgproxy(Config{
    BaseURL:       "http://localhost",
    SignatureSize: 15,
    Key:           hex.EncodeToString([]byte("key")),
    Salt:          hex.EncodeToString([]byte("salt")),
    Encode:        false,
  })
  if err != nil {
    panic(err)
  }

  url, err := ip.
    Builder().
    Resize(ResizingTypeFill, 123, 456, true).
    DPR(10).
    Format("png").
    Generate("path/to/my/image.jpg")
  if err != nil {
    panic(err)
  }

  fmt.Println(url) // http://localhost/QMYScvaF2YUPdA-NJl8E/dpr:10/f:png/rs:fill:123:456:1/plain/path/to/my/image.jpg
```

## Tests

```bash
  $ go test
```

## License

This project is licensed under the [MIT License](LICENSE.md).
