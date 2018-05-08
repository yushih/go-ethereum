package loads

import "github.com/ethereum/go-ethereum/core/jvm/instructions/base"
import "github.com/ethereum/go-ethereum/core/jvm/rtda"

// Load double from local variable
type DLOAD struct{ base.Index8Instruction }

func (self *DLOAD) Execute(frame *rtda.Frame, gas uint64, contract interface{}) uint64 {
	_dload(frame, uint(self.Index))
    return 100
}

type DLOAD_0 struct{ base.NoOperandsInstruction }

func (self *DLOAD_0) Execute(frame *rtda.Frame, gas uint64, contract interface{}) uint64 {
	_dload(frame, 0)
    return 100
}

type DLOAD_1 struct{ base.NoOperandsInstruction }

func (self *DLOAD_1) Execute(frame *rtda.Frame, gas uint64, contract interface{}) uint64 {
	_dload(frame, 1)
    return 100
}

type DLOAD_2 struct{ base.NoOperandsInstruction }

func (self *DLOAD_2) Execute(frame *rtda.Frame, gas uint64, contract interface{}) uint64 {
	_dload(frame, 2)
    return 100
}

type DLOAD_3 struct{ base.NoOperandsInstruction }

func (self *DLOAD_3) Execute(frame *rtda.Frame, gas uint64, contract interface{}) uint64 {
	_dload(frame, 3)
    return 100
}

func _dload(frame *rtda.Frame, index uint) {
	val := frame.LocalVars().GetDouble(index)
	frame.OperandStack().PushDouble(val)
}
