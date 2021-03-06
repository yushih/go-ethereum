package stores

import "github.com/ethereum/go-ethereum/core/jvm/instructions/base"
import "github.com/ethereum/go-ethereum/core/jvm/rtda"

// Store long into local variable
type LSTORE struct{ base.Index8Instruction }

func (self *LSTORE) Execute(frame *rtda.Frame, gas uint64, contract interface{}, evm interface{}) uint64 {
	_lstore(frame, uint(self.Index))
    return 100
}

type LSTORE_0 struct{ base.NoOperandsInstruction }

func (self *LSTORE_0) Execute(frame *rtda.Frame, gas uint64, contract interface{}, evm interface{}) uint64 {
	_lstore(frame, 0)
    return 100
}

type LSTORE_1 struct{ base.NoOperandsInstruction }

func (self *LSTORE_1) Execute(frame *rtda.Frame, gas uint64, contract interface{}, evm interface{}) uint64 {
	_lstore(frame, 1)
    return 100
}

type LSTORE_2 struct{ base.NoOperandsInstruction }

func (self *LSTORE_2) Execute(frame *rtda.Frame, gas uint64, contract interface{}, evm interface{}) uint64 {
	_lstore(frame, 2)
    return 100
}

type LSTORE_3 struct{ base.NoOperandsInstruction }

func (self *LSTORE_3) Execute(frame *rtda.Frame, gas uint64, contract interface{}, evm interface{}) uint64 {
	_lstore(frame, 3)
    return 100
}

func _lstore(frame *rtda.Frame, index uint) {
	val := frame.OperandStack().PopLong()
	frame.LocalVars().SetLong(index, val)
}
