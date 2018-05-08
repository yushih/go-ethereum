package io

import "github.com/ethereum/go-ethereum/core/jvm/native"
import "github.com/ethereum/go-ethereum/core/jvm/rtda"

const fd = "java/io/FileDescriptor"

func init() {
	native.Register(fd, "set", "(I)J", set)
}

// private static native long set(int d);
// (I)J
func set(frame *rtda.Frame, gas uint64, contract interface{}, evm interface{}) {
	// todo
	frame.OperandStack().PushLong(0)
}
