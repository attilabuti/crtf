# crtf

> This code is ported from the Python implementation available at: [delimitry/compressed_rtf](https://github.com/delimitry/compressed_rtf)

The `crtf` package provides functionality for compressing and decompressing [Compressed Rich Text Format](https://msdn.microsoft.com/en-us/library/cc463890(v=exchg.80).aspx) (RTF) (also known as "LZFu" compression format).

## Installation

To install the package, use the `go get` command:

```bash
$ go get github.com/attilabuti/crtf@latest
```

## Usage

### Decompressing RTF

```go
import (
    "fmt"
    "github.com/attilabuti/crtf"
)

func main() {
    compressedData := []byte{/* your compressed RTF data */}
    decompressedData, err := crtf.Decompress(compressedData)
    if err != nil {
        fmt.Println("Error:", err)
        return
    }

    fmt.Println("Decompressed data:", string(decompressedData))
}
```

### Compressing RTF

```go
import (
    "fmt"
    "github.com/attilabuti/crtf"
)

func main() {
    data := []byte("Your RTF data")
    compressedData := crtf.Compress(data, true)  // Set to true for compression
    fmt.Println("Compressed data:", compressedData)

    // If you want to write the data uncompressed, set the second parameter to false
    uncompressedData := crtf.Compress(data, false)
    fmt.Println("Uncompressed data:", uncompressedData)
}
```

## Issues

Submit the [issues](https://github.com/attilabuti/crtf/issues) if you find any bug or have any suggestion.

## Contribution

Fork the [repo](https://github.com/attilabuti/crtf) and submit pull requests.

## License

This extension is licensed under the [MIT License](https://github.com/attilabuti/crtf/blob/main/LICENSE).