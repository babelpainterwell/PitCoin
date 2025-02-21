package main

import (
	"bytes"
	"fmt"
	"testing"
)

// func BenchmarkByteSlice(b *testing.B) {
//     var data []byte
//     for i := 0; i < b.N; i++ {
//         data = append(data, make([]byte, 1000)...)
//     }
// }

func BenchmarkByteBuffer(b *testing.B) {
    var buf bytes.Buffer
    for i := 0; i < b.N; i++ {
        buf.Write(make([]byte, 1000))
		fmt.Printf("%d\n", buf.Available())
    }
}

func main() {
    // Run both in the same process
    // fmt.Println(testing.Benchmark(BenchmarkByteSlice))
    fmt.Println(testing.Benchmark(BenchmarkByteBuffer))
}
