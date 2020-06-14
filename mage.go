package main

// `cargo build --release --features ffi` to use this!
// On Linux LD_LIBRARY_PATH=target/release env is required!

// #cgo LDFLAGS: -Ltarget/release -lmage
// #cgo CFLAGS: -Itarget/release
// #include "target/release/mage.h"
import "C"

import (
	"fmt"
	"unsafe"
)

type Listener struct {
	index C.ulong
}

func Listen(addr string) *Listener {
	addrPtr := C.CString(addr)
	addrCharPtr := (*C.schar)(addrPtr)

	i := C.ffi_listen_opt(addrCharPtr)

	C.free(unsafe.Pointer(addrPtr))

	return &Listener{index: i}
}

func (l *Listener) Accept(seed [32]byte, key [32]byte) *StreamChanneled {
	seedPtr := unsafe.Pointer(&seed[0])
	seedCharPtr := (*C.uchar)(seedPtr)
	keyPtr := unsafe.Pointer(&key[0])
	keyCharPtr := (*C.uchar)(keyPtr)

	i := C.ffi_accept_opt(l.index, 1, seedCharPtr, keyCharPtr)
	fmt.Printf("[Go] New accept: %d\n", uint64(i))

	return &StreamChanneled{index: i}
}

type StreamChanneled struct {
	index C.ulong
}

func Connect(addr string, seed [32]byte, key [32]byte) *StreamChanneled {
	addrPtr := C.CString(addr)
	addrCharPtr := (*C.schar)(addrPtr)
	seedPtr := unsafe.Pointer(&seed[0])
	seedCharPtr := (*C.uchar)(seedPtr)
	keyPtr := unsafe.Pointer(&key[0])
	keyCharPtr := (*C.uchar)(keyPtr)

	i := C.ffi_connect_opt(addrCharPtr, 0, seedCharPtr, keyCharPtr)
	fmt.Printf("[Go] New connect: %d\n", uint64(i))

	C.free(unsafe.Pointer(addrPtr))

	return &StreamChanneled{index: i}
}

func (c *StreamChanneled) GetChannel(i byte) *Channel {
	r := C.ffi_get_channel(c.index, C.uchar(i))

	return &Channel{index: r}
}

func (c *StreamChanneled) ChannelLoop() {
	C.ffi_channel_loop(c.index)
}

func (c *StreamChanneled) Read(buffer []byte) (int, error) {
	bufferLen := C.ulong(len(buffer))
	bufferPtr := unsafe.Pointer(&buffer[0])

	r := C.ffi_recv(c.index, bufferPtr, bufferLen)

	return int(r), nil
}

func (c *StreamChanneled) Write(buffer []byte) (int, error) {
	bufferLen := C.ulong(len(buffer))
	bufferPtr := unsafe.Pointer(&buffer[0])

	r := C.ffi_send(c.index, bufferPtr, bufferLen)

	return int(r), nil
}

type Channel struct {
	index C.ulong
}

func (c *Channel) Read(buffer []byte) (int, error) {
	bufferLen := C.ulong(len(buffer))
	bufferPtr := unsafe.Pointer(&buffer[0])

	r := C.ffi_recv_channel(c.index, bufferPtr, bufferLen)

	return int(r), nil
}

func (c *Channel) Write(buffer []byte) (int, error) {
	bufferLen := C.ulong(len(buffer))
	bufferPtr := unsafe.Pointer(&buffer[0])

	r := C.ffi_send_channel(c.index, bufferPtr, bufferLen)
	return int(r), nil
}

func main() {
}
