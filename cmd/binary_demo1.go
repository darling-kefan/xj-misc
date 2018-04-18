package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"math/rand"
	"os"
	"path"
	"time"
)

func bin_write() {
	t := time.Now().Nanosecond()
	fmt.Println(path.Join("bin", "numbers.binary"))
	fp, err := os.Create(path.Join("bin", "numbers.binary"))
	if err != nil {
		fmt.Println(err)
	}
	defer fp.Close()

	rand.Seed(int64(t))

	buf := new(bytes.Buffer)
	for i := 0; i < 10; i++ {
		binary.Write(buf, binary.LittleEndian, int32(i))
		fmt.Printf("% x\n", buf.Bytes())
		fp.Write(buf.Bytes())
		buf.Reset()
	}
}

func bin_read() {
	fp, _ := os.Open(path.Join("bin", "numbers.binary"))
	defer fp.Close()

	data := make([]byte, 4)
	var k int32
	for {
		// data = data[:cap(data)]

		n, err := fp.Read(data)
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println(err)
			break
		}

		data = data[:n]
		binary.Read(bytes.NewBuffer(data), binary.LittleEndian, &k)
		fmt.Println(k)
	}
}

func main() {
	//bin_write()
	bin_read()
}
