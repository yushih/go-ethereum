package conversions

import "github.com/ethereum/go-ethereum/core/jvm/instructions/base"
import "github.com/ethereum/go-ethereum/core/jvm/rtda"

// Convert float to double
type F2D struct{ base.NoOperandsInstruction }

func (self *F2D) Execute(frame *rtda.Frame, gas uint64, contract interface{}) uint64 {
	stack := frame.OperandStack()
	f := stack.PopFloat()
	d := float64(f)
	stack.PushDouble(d)
    return 100
}

// Convert float to int
type F2I struct{ base.NoOperandsInstruction }

func (self *F2I) Execute(frame *rtda.Frame, gas uint64, contract interface{}) uint64 {
	stack := frame.OperandStack()
	f := stack.PopFloat()
	i := int32(f)
	stack.PushInt(i)
    return 100
}

// Convert float to long
type F2L struct{ base.NoOperandsInstruction }

func (self *F2L) Execute(frame *rtda.Frame, gas uint64, contract interface{}) uint64 {
	stack := frame.OperandStack()
	f := stack.PopFloat()
	l := int64(f)
	stack.PushLong(l)
    return 100
}
