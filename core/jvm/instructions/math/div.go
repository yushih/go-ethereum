package math

import "github.com/ethereum/go-ethereum/core/jvm/instructions/base"
import "github.com/ethereum/go-ethereum/core/jvm/rtda"

// Divide double
type DDIV struct{ base.NoOperandsInstruction }

func (self *DDIV) Execute(frame *rtda.Frame, gas uint64, contract interface{}, evm interface{}) uint64 {
	stack := frame.OperandStack()
	v2 := stack.PopDouble()
	v1 := stack.PopDouble()
	result := v1 / v2
	stack.PushDouble(result)
    return 100
}

// Divide float
type FDIV struct{ base.NoOperandsInstruction }

func (self *FDIV) Execute(frame *rtda.Frame, gas uint64, contract interface{}, evm interface{}) uint64 {
	stack := frame.OperandStack()
	v2 := stack.PopFloat()
	v1 := stack.PopFloat()
	result := v1 / v2
	stack.PushFloat(result)
    return 100
}

// Divide int
type IDIV struct{ base.NoOperandsInstruction }

func (self *IDIV) Execute(frame *rtda.Frame, gas uint64, contract interface{}, evm interface{}) uint64 {
	stack := frame.OperandStack()
	v2 := stack.PopInt()
	v1 := stack.PopInt()
	if v2 == 0 {
		panic("java.lang.ArithmeticException: / by zero")
	}

	result := v1 / v2
	stack.PushInt(result)
    return 100
}

// Divide long
type LDIV struct{ base.NoOperandsInstruction }

func (self *LDIV) Execute(frame *rtda.Frame, gas uint64, contract interface{}, evm interface{}) uint64 {
	stack := frame.OperandStack()
	v2 := stack.PopLong()
	v1 := stack.PopLong()
	if v2 == 0 {
		panic("java.lang.ArithmeticException: / by zero")
	}

	result := v1 / v2
	stack.PushLong(result)
    return 100
}
