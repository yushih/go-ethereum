package loads

import "github.com/ethereum/go-ethereum/core/jvm/instructions/base"
import "github.com/ethereum/go-ethereum/core/jvm/rtda"

// Load float from local variable
type FLOAD struct{ base.Index8Instruction }

func (self *FLOAD) Execute(frame *rtda.Frame, gas uint64) uint64 {
	_fload(frame, uint(self.Index))
    return 100
}

type FLOAD_0 struct{ base.NoOperandsInstruction }

func (self *FLOAD_0) Execute(frame *rtda.Frame, gas uint64) uint64 {
	_fload(frame, 0)
    return 100
}

type FLOAD_1 struct{ base.NoOperandsInstruction }

func (self *FLOAD_1) Execute(frame *rtda.Frame, gas uint64) uint64 {
	_fload(frame, 1)
    return 100
}

type FLOAD_2 struct{ base.NoOperandsInstruction }

func (self *FLOAD_2) Execute(frame *rtda.Frame, gas uint64) uint64 {
	_fload(frame, 2)
    return 100
}

type FLOAD_3 struct{ base.NoOperandsInstruction }

func (self *FLOAD_3) Execute(frame *rtda.Frame, gas uint64) uint64 {
	_fload(frame, 3)
    return 100
}

func _fload(frame *rtda.Frame, index uint) {
	val := frame.LocalVars().GetFloat(index)
	frame.OperandStack().PushFloat(val)
}
