package loads

import "github.com/ethereum/go-ethereum/core/jvm/instructions/base"
import "github.com/ethereum/go-ethereum/core/jvm/rtda"

// Load long from local variable
type LLOAD struct{ base.Index8Instruction }

func (self *LLOAD) Execute(frame *rtda.Frame, gas uint64) uint64 {
	_lload(frame, uint(self.Index))
    return 100
}

type LLOAD_0 struct{ base.NoOperandsInstruction }

func (self *LLOAD_0) Execute(frame *rtda.Frame, gas uint64) uint64 {
	_lload(frame, 0)
    return 100
}

type LLOAD_1 struct{ base.NoOperandsInstruction }

func (self *LLOAD_1) Execute(frame *rtda.Frame, gas uint64) uint64 {
	_lload(frame, 1)
    return 100
}

type LLOAD_2 struct{ base.NoOperandsInstruction }

func (self *LLOAD_2) Execute(frame *rtda.Frame, gas uint64) uint64 {
	_lload(frame, 2)
    return 100
}

type LLOAD_3 struct{ base.NoOperandsInstruction }

func (self *LLOAD_3) Execute(frame *rtda.Frame, gas uint64) uint64 {
	_lload(frame, 3)
    return 100
}

func _lload(frame *rtda.Frame, index uint) {
	val := frame.LocalVars().GetLong(index)
	frame.OperandStack().PushLong(val)
}
