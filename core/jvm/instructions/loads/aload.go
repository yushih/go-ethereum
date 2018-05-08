package loads

import "github.com/ethereum/go-ethereum/core/jvm/instructions/base"
import "github.com/ethereum/go-ethereum/core/jvm/rtda"

// Load reference from local variable
type ALOAD struct{ base.Index8Instruction }

func (self *ALOAD) Execute(frame *rtda.Frame, gas uint64, contract interface{}, evm interface{}) uint64 {
	_aload(frame, uint(self.Index))
    return 100
}

type ALOAD_0 struct{ base.NoOperandsInstruction }

func (self *ALOAD_0) Execute(frame *rtda.Frame, gas uint64, contract interface{}, evm interface{}) uint64 {
	_aload(frame, 0)
    return 100
}

type ALOAD_1 struct{ base.NoOperandsInstruction }

func (self *ALOAD_1) Execute(frame *rtda.Frame, gas uint64, contract interface{}, evm interface{}) uint64 {
	_aload(frame, 1)
    return 100
}

type ALOAD_2 struct{ base.NoOperandsInstruction }

func (self *ALOAD_2) Execute(frame *rtda.Frame, gas uint64, contract interface{}, evm interface{}) uint64 {
	_aload(frame, 2)
    return 100
}

type ALOAD_3 struct{ base.NoOperandsInstruction }

func (self *ALOAD_3) Execute(frame *rtda.Frame, gas uint64, contract interface{}, evm interface{}) uint64 {
	_aload(frame, 3)
    return 100
}

func _aload(frame *rtda.Frame, index uint) {
	ref := frame.LocalVars().GetRef(index)
	frame.OperandStack().PushRef(ref)
}
