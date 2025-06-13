# go-metis

Go bindings for [METIS](https://github.com/KarypisLab/METIS) - Serial Graph Partitioning and Fill-reducing Matrix Ordering

[![Go Reference](https://pkg.go.dev/badge/github.com/notargets/go-metis.svg)](https://pkg.go.dev/github.com/notargets/go-metis)
[![Go Report Card](https://goreportcard.com/badge/github.com/notargets/go-metis)](https://goreportcard.com/report/github.com/notargets/go-metis)
[![CI](https://github.com/notargets/go-metis/actions/workflows/test.yml/badge.svg)](https://github.com/notargets/go-metis/actions/workflows/test.yml)

## Installation

```bash
go get github.com/notargets/go-metis
```

## Requirements

- METIS library (5.1.0 or later)
- CGO-enabled Go installation

### Installing METIS

#### Ubuntu/Debian
```bash
# Install dependencies
sudo apt-get install build-essential cmake

# Build from source
git clone https://github.com/KarypisLab/GKlib.git
git clone https://github.com/KarypisLab/METIS.git

cd GKlib
make config prefix=/usr/local
sudo make install

cd ../METIS
make config shared=1 prefix=/usr/local gklib_path=/usr/local
sudo make install
sudo ldconfig
```

#### macOS
```bash
brew install metis
```

## Usage

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/notargets/go-metis"
)

func main() {
    // Simple 4-node graph (square)
    //  0 -- 1
    //  |    |
    //  2 -- 3
    xadj := []int32{0, 2, 4, 6, 8}
    adjncy := []int32{1, 2, 0, 3, 0, 3, 1, 2}
    
    // Partition into 2 parts
    part, edgeCut, err := metis.PartitionGraph(xadj, adjncy, 2)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Partitioning: %v\n", part)
    fmt.Printf("Edge cut: %d\n", edgeCut)
}
```

## Features

- [x] Graph partitioning (k-way and recursive)
- [x] Mesh partitioning
- [x] Matrix reordering
- [x] Graph coarsening
- [ ] Nested dissection (coming soon)

## Documentation

See [GoDoc](https://pkg.go.dev/github.com/notargets/go-metis) for detailed API documentation.

## Examples

Check the [examples](examples/) directory for more usage examples.

## Testing

```bash
go test -v ./...
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

Note: METIS itself is licensed under the Apache License 2.0.
