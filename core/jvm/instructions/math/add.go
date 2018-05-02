package math

import "github.com/ethereum/go-ethereum/core/jvm/instructions/base"
import "github.com/ethereum/go-ethereum/core/jvm/rtda"

// Add double
type DADD struct{ base.NoOperandsInstruction }

func (self *DADD) Execute(frame *rtda.Frame, gas uint64) uint64 {
	stack := frame.OperandStack()
	v1 := stack.PopDouble()
	v2 := stack.PopDouble()
	result := v1 + v2
	stack.PushDouble(result)
    return 100
}

// Add float
type FADD struct{ base.NoOperandsInstruction }

func (self *FADD) Execute(frame *rtda.Frame, gas uint64) uint64 {
	stack := frame.OperandStack()
	v2 := stack.PopFloat()
	v1 := stack.PopFloat()
	result := v1 + v2
	stack.PushFloat(result)
    return 100
}

// Add int
type IADD struct{ base.NoOperandsInstruction }

func (self *IADD) Execute(frame *rtda.Frame, gas uint64) uint64 {
	stack := frame.OperandStack()
	v2 := stack.PopInt()
	v1 := stack.PopInt()
	result := v1 + v2
	stack.PushInt(result)
    return 100
}

// Add long
type LADD struct{ base.NoOperandsInstruction }

func (self *LADD) Execute(frame *rtda.Frame, gas uint64) uint64 {
	stack := frame.OperandStack()
	v2 := stack.PopLong()
	v1 := stack.PopLong()
	result := v1 + v2
	stack.PushLong(result)
    return 100
}
