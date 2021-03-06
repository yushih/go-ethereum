package atomic

import "github.com/ethereum/go-ethereum/core/jvm/native"
import "github.com/ethereum/go-ethereum/core/jvm/rtda"

func init() {
	native.Register("java/util/concurrent/atomic/AtomicLong", "VMSupportsCS8", "()Z", vmSupportsCS8)
}

func vmSupportsCS8(frame *rtda.Frame, gas uint64, contract interface{}, evm interface{}) {
	frame.OperandStack().PushBoolean(false)
}
