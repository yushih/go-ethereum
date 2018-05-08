package constants

import "github.com/ethereum/go-ethereum/core/jvm/instructions/base"
import "github.com/ethereum/go-ethereum/core/jvm/rtda"

// Push null
type ACONST_NULL struct{ base.NoOperandsInstruction }

func (self *ACONST_NULL) Execute(frame *rtda.Frame, gas uint64, contract interface{}, evm interface{}) uint64 {
	frame.OperandStack().PushRef(nil)
    return 100
}

// Push double
type DCONST_0 struct{ base.NoOperandsInstruction }

func (self *DCONST_0) Execute(frame *rtda.Frame, gas uint64, contract interface{}, evm interface{}) uint64 {
	frame.OperandStack().PushDouble(0.0)
    return 100
}

type DCONST_1 struct{ base.NoOperandsInstruction }

func (self *DCONST_1) Execute(frame *rtda.Frame, gas uint64, contract interface{}, evm interface{}) uint64 {
	frame.OperandStack().PushDouble(1.0)
    return 100
}

// Push float
type FCONST_0 struct{ base.NoOperandsInstruction }

func (self *FCONST_0) Execute(frame *rtda.Frame, gas uint64, contract interface{}, evm interface{}) uint64 {
	frame.OperandStack().PushFloat(0.0)
    return 100
}

type FCONST_1 struct{ base.NoOperandsInstruction }

func (self *FCONST_1) Execute(frame *rtda.Frame, gas uint64, contract interface{}, evm interface{}) uint64 {
	frame.OperandStack().PushFloat(1.0)
    return 100
}

type FCONST_2 struct{ base.NoOperandsInstruction }

func (self *FCONST_2) Execute(frame *rtda.Frame, gas uint64, contract interface{}, evm interface{}) uint64 {
	frame.OperandStack().PushFloat(2.0)
    return 100
}

// Push int constant
type ICONST_M1 struct{ base.NoOperandsInstruction }

func (self *ICONST_M1) Execute(frame *rtda.Frame, gas uint64, contract interface{}, evm interface{}) uint64 {
	frame.OperandStack().PushInt(-1)
    return 100
}

type ICONST_0 struct{ base.NoOperandsInstruction }

func (self *ICONST_0) Execute(frame *rtda.Frame, gas uint64, contract interface{}, evm interface{}) uint64 {
	frame.OperandStack().PushInt(0)
    return 100
}

type ICONST_1 struct{ base.NoOperandsInstruction }

func (self *ICONST_1) Execute(frame *rtda.Frame, gas uint64, contract interface{}, evm interface{}) uint64 {
	frame.OperandStack().PushInt(1)
    return 100
}

type ICONST_2 struct{ base.NoOperandsInstruction }

func (self *ICONST_2) Execute(frame *rtda.Frame, gas uint64, contract interface{}, evm interface{}) uint64 {
	frame.OperandStack().PushInt(2)
    return 100
}

type ICONST_3 struct{ base.NoOperandsInstruction }

func (self *ICONST_3) Execute(frame *rtda.Frame, gas uint64, contract interface{}, evm interface{}) uint64 {
	frame.OperandStack().PushInt(3)
    return 100
}

type ICONST_4 struct{ base.NoOperandsInstruction }

func (self *ICONST_4) Execute(frame *rtda.Frame, gas uint64, contract interface{}, evm interface{}) uint64 {
	frame.OperandStack().PushInt(4)
    return 100
}

type ICONST_5 struct{ base.NoOperandsInstruction }

func (self *ICONST_5) Execute(frame *rtda.Frame, gas uint64, contract interface{}, evm interface{}) uint64 {
	frame.OperandStack().PushInt(5)
    return 100
}

// Push long constant
type LCONST_0 struct{ base.NoOperandsInstruction }

func (self *LCONST_0) Execute(frame *rtda.Frame, gas uint64, contract interface{}, evm interface{}) uint64 {
	frame.OperandStack().PushLong(0)
    return 100
}

type LCONST_1 struct{ base.NoOperandsInstruction }

func (self *LCONST_1) Execute(frame *rtda.Frame, gas uint64, contract interface{}, evm interface{}) uint64 {
	frame.OperandStack().PushLong(1)
    return 100
}
