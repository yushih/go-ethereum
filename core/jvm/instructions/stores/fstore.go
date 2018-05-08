package stores

import "github.com/ethereum/go-ethereum/core/jvm/instructions/base"
import "github.com/ethereum/go-ethereum/core/jvm/rtda"

// Store float into local variable
type FSTORE struct{ base.Index8Instruction }

func (self *FSTORE) Execute(frame *rtda.Frame, gas uint64, contract interface{}) uint64 {
	_fstore(frame, uint(self.Index))
    return 100
}

type FSTORE_0 struct{ base.NoOperandsInstruction }

func (self *FSTORE_0) Execute(frame *rtda.Frame, gas uint64, contract interface{}) uint64 {
	_fstore(frame, 0)
    return 100
}

type FSTORE_1 struct{ base.NoOperandsInstruction }

func (self *FSTORE_1) Execute(frame *rtda.Frame, gas uint64, contract interface{}) uint64 {
	_fstore(frame, 1)
    return 100
}

type FSTORE_2 struct{ base.NoOperandsInstruction }

func (self *FSTORE_2) Execute(frame *rtda.Frame, gas uint64, contract interface{}) uint64 {
	_fstore(frame, 2)
    return 100
}

type FSTORE_3 struct{ base.NoOperandsInstruction }

func (self *FSTORE_3) Execute(frame *rtda.Frame, gas uint64, contract interface{}) uint64 {
	_fstore(frame, 3)
    return 100
}

func _fstore(frame *rtda.Frame, index uint) {
	val := frame.OperandStack().PopFloat()
	frame.LocalVars().SetFloat(index, val)
}
