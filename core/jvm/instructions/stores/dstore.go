package stores

import "github.com/ethereum/go-ethereum/core/jvm/instructions/base"
import "github.com/ethereum/go-ethereum/core/jvm/rtda"

// Store double into local variable
type DSTORE struct{ base.Index8Instruction }

func (self *DSTORE) Execute(frame *rtda.Frame, gas uint64, contract interface{}, evm interface{}) uint64 {
	_dstore(frame, uint(self.Index))
    return 100
}

type DSTORE_0 struct{ base.NoOperandsInstruction }

func (self *DSTORE_0) Execute(frame *rtda.Frame, gas uint64, contract interface{}, evm interface{}) uint64 {
	_dstore(frame, 0)
    return 100
}

type DSTORE_1 struct{ base.NoOperandsInstruction }

func (self *DSTORE_1) Execute(frame *rtda.Frame, gas uint64, contract interface{}, evm interface{}) uint64 {
	_dstore(frame, 1)
    return 100
}

type DSTORE_2 struct{ base.NoOperandsInstruction }

func (self *DSTORE_2) Execute(frame *rtda.Frame, gas uint64, contract interface{}, evm interface{}) uint64 {
	_dstore(frame, 2)
    return 100
}

type DSTORE_3 struct{ base.NoOperandsInstruction }

func (self *DSTORE_3) Execute(frame *rtda.Frame, gas uint64, contract interface{}, evm interface{}) uint64 {
	_dstore(frame, 3)
    return 100
}

func _dstore(frame *rtda.Frame, index uint) {
	val := frame.OperandStack().PopDouble()
	frame.LocalVars().SetDouble(index, val)
}
