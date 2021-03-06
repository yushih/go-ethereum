package conversions

import "github.com/ethereum/go-ethereum/core/jvm/instructions/base"
import "github.com/ethereum/go-ethereum/core/jvm/rtda"

// Convert long to double
type L2D struct{ base.NoOperandsInstruction }

func (self *L2D) Execute(frame *rtda.Frame, gas uint64, contract interface{}, evm interface{}) uint64 {
	stack := frame.OperandStack()
	l := stack.PopLong()
	d := float64(l)
	stack.PushDouble(d)
    return 100
}

// Convert long to float
type L2F struct{ base.NoOperandsInstruction }

func (self *L2F) Execute(frame *rtda.Frame, gas uint64, contract interface{}, evm interface{}) uint64 {
	stack := frame.OperandStack()
	l := stack.PopLong()
	f := float32(l)
	stack.PushFloat(f)
    return 100
}

// Convert long to int
type L2I struct{ base.NoOperandsInstruction }

func (self *L2I) Execute(frame *rtda.Frame, gas uint64, contract interface{}, evm interface{}) uint64 {
	stack := frame.OperandStack()
	l := stack.PopLong()
	i := int32(l)
	stack.PushInt(i)
    return 100
}
