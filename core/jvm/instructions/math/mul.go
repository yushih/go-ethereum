package math

import "github.com/ethereum/go-ethereum/core/jvm/instructions/base"
import "github.com/ethereum/go-ethereum/core/jvm/rtda"

// Multiply double
type DMUL struct{ base.NoOperandsInstruction }

func (self *DMUL) Execute(frame *rtda.Frame, gas uint64, contract interface{}, evm interface{}) uint64 {
	stack := frame.OperandStack()
	v2 := stack.PopDouble()
	v1 := stack.PopDouble()
	result := v1 * v2
	stack.PushDouble(result)
    return 100
}

// Multiply float
type FMUL struct{ base.NoOperandsInstruction }

func (self *FMUL) Execute(frame *rtda.Frame, gas uint64, contract interface{}, evm interface{}) uint64 {
	stack := frame.OperandStack()
	v2 := stack.PopFloat()
	v1 := stack.PopFloat()
	result := v1 * v2
	stack.PushFloat(result)
    return 100
}

// Multiply int
type IMUL struct{ base.NoOperandsInstruction }

func (self *IMUL) Execute(frame *rtda.Frame, gas uint64, contract interface{}, evm interface{}) uint64 {
	stack := frame.OperandStack()
	v2 := stack.PopInt()
	v1 := stack.PopInt()
	result := v1 * v2
	stack.PushInt(result)
    return 100
}

// Multiply long
type LMUL struct{ base.NoOperandsInstruction }

func (self *LMUL) Execute(frame *rtda.Frame, gas uint64, contract interface{}, evm interface{}) uint64 {
	stack := frame.OperandStack()
	v2 := stack.PopLong()
	v1 := stack.PopLong()
	result := v1 * v2
	stack.PushLong(result)
    return 100
}
