package loads

import "github.com/ethereum/go-ethereum/core/jvm/instructions/base"
import "github.com/ethereum/go-ethereum/core/jvm/rtda"

// Load int from local variable
type ILOAD struct{ base.Index8Instruction }

func (self *ILOAD) Execute(frame *rtda.Frame, gas uint64) uint64 {
	_iload(frame, uint(self.Index))
    return 100
}

type ILOAD_0 struct{ base.NoOperandsInstruction }

func (self *ILOAD_0) Execute(frame *rtda.Frame, gas uint64) uint64 {
	_iload(frame, 0)
    return 100
}

type ILOAD_1 struct{ base.NoOperandsInstruction }

func (self *ILOAD_1) Execute(frame *rtda.Frame, gas uint64) uint64 {
	_iload(frame, 1)
    return 100
}

type ILOAD_2 struct{ base.NoOperandsInstruction }

func (self *ILOAD_2) Execute(frame *rtda.Frame, gas uint64) uint64 {
	_iload(frame, 2)
    return 100
}

type ILOAD_3 struct{ base.NoOperandsInstruction }

func (self *ILOAD_3) Execute(frame *rtda.Frame, gas uint64) uint64 {
	_iload(frame, 3)
    return 100
}

func _iload(frame *rtda.Frame, index uint) {
	val := frame.LocalVars().GetInt(index)
	frame.OperandStack().PushInt(val)
}
