package math

import "github.com/ethereum/go-ethereum/core/jvm/instructions/base"
import "github.com/ethereum/go-ethereum/core/jvm/rtda"

// Boolean OR int
type IOR struct{ base.NoOperandsInstruction }

func (self *IOR) Execute(frame *rtda.Frame, gas uint64, contract interface{}, evm interface{}) uint64 {
	stack := frame.OperandStack()
	v2 := stack.PopInt()
	v1 := stack.PopInt()
	result := v1 | v2
	stack.PushInt(result)
    return 100
}

// Boolean OR long
type LOR struct{ base.NoOperandsInstruction }

func (self *LOR) Execute(frame *rtda.Frame, gas uint64, contract interface{}, evm interface{}) uint64 {
	stack := frame.OperandStack()
	v2 := stack.PopLong()
	v1 := stack.PopLong()
	result := v1 | v2
	stack.PushLong(result)
    return 100
}
