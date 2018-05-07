package interf

import (
    "github.com/ethereum/go-ethereum/core/jvm/native"
    "github.com/ethereum/go-ethereum/core/jvm/rtda"
)

func init() {
    native.Register("blockchain/Special", "gasLeft", "()J", gasLeft)
}

func gasLeft(frame *rtda.Frame) {
     // see invokenative.go
}