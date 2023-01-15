package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"syscall"
	"unsafe"
)

func main() {
	mem := mmapRegs32(0, 0x200)
	
	for i, v := range mem {
		fmt.Printf("%08x: %08x\n", i*4, v)
	}
	fmt.Println("update data")
	mem[0] = 0x39393939
	for i, v := range mem {
		fmt.Printf("%08x: %08x\n", i*4, v)
	}

}

func mmapRegs32(addr uint32, n int) []uint32 {
	file := "testdev0"
	f, err := os.OpenFile(file,os.O_RDWR | os.O_CREATE,0777)
	if err != nil {
		log.Fatal(err)
	}

	offset := int64(addr) &^ 0xfff
	data, err := syscall.Mmap(int(f.Fd()), offset, (n+0xfff)&^0xfff, syscall.PROT_READ | syscall.PROT_WRITE, syscall.MAP_SHARED)
	if err != nil {
		log.Fatalf("mmap %s: %v", f.Name(), err)
	}
	f.Close()

	map_array := (*[math.MaxInt32 / 4]uint32)(unsafe.Pointer(&data[0]))
	return map_array[(addr&0xfff)/4 : ((int(addr)&0xfff)+n)/4]
}
