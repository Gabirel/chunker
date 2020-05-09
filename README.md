# chunker

## What does this project do?

This is a universal wrapper for application-level to use chunking in Go.
The original chunker module is designed for [restic backup program](https://github.com/restic/restic).
To be honest, it is really hard to use it **RIGHT AWAY**. You need to take time to find it out how does it work and how to apply it into your code. It's not worth that.

So I make this.
You can find the chunker module from here: [restic/chunker](https://github.com/restic/chunker). 



## What is the default setting then?

Current setting for chunking are as follows:

- minimum chunk size : **2KB**
- maximum chunk size : **16KB**
- average chunk size : **8KB**

You can change this setting via changing bits:

- 2KB is `1 << 11`
- 16KB is `1 << 14`
- 8KB is `1 << 13`

Change as you want. Corresponding code is: 

```go
NewWithBoundaries(rd io.Reader, minBits, maxBits, averageBits int);
```

If you doesn't care about the default setting, just use:

```go
chnker := chunker.New(reader)
```



## Usage/Example

You can use chunker as bellow:

```go
package main

import (
    "fmt"
    chunk "github.com/Gabirel/chunker"
    "io"
    "os"
)

func main() {
    fmt.Println("Starting to open file...")

    filename := "test.txt"
    reader, err := os.Open(filename)
    if err != nil {
        fmt.Println("Read error: ", err)
        _ = reader.Close()
        return
    }

    chunker := chunk.New(reader)

    var totalLength uint = 0
    var averageLength uint = 0
    var chunkNum uint = 0

    for {
        chunkInfo, err := chunker.Chunking()
        if err == io.EOF {
            fmt.Println("\nEOF. Exiting...")
            fmt.Println("Total length = ", totalLength)
            fmt.Println("Chunk Number = ", chunkNum)
            if chunkNum != 0 {
                averageLength = totalLength / chunkNum
                fmt.Println("Average length = ", averageLength)
            }
            break
        }

        if chunkInfo == nil {
            panic("chunkInfo is nil!")
        }
        fmt.Printf("%d:\t%x\n", chunkInfo.Chunk.Length, chunkInfo.Digest)
        totalLength += chunkInfo.Chunk.Length
        chunkNum += 1
    }
}
```

