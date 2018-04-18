package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"path"
	"unsafe"
)

type MyData struct {
	_ [1]byte
	Y int32
	X int32
	Z int32
}

func struct_write() {
	fp, _ := os.Create(path.Join("bin", "struct.binary"))
	defer fp.Close()

	// 将结构体转成bytes, 按照字段声明顺序，但是"_"放在最后．
	data := &MyData{X: 1, Y: 2, Z: 3}
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, data)

	// 将bytes写入文件
	fp.Write(buf.Bytes())
	fp.Sync()
}

func struct_read() {
	fp, _ := os.Open(path.Join("bin", "struct.binary"))
	defer fp.Close()

	fmt.Println(unsafe.Sizeof(MyData{}), binary.Size(MyData{}))

	// dataBytes := make([]byte, unsafe.Sizeof(MyData{}))
	dataBytes := make([]byte, binary.Size(MyData{}))
	data := MyData{}
	n, _ := fp.Read(dataBytes)
	dataBytes = dataBytes[:n]

	binary.Read(bytes.NewBuffer(dataBytes), binary.LittleEndian, &data)
	fmt.Println(data)
}

func main() {
	// struct_write()
	struct_read()
}
