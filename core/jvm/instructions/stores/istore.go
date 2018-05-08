package stores

import "github.com/ethereum/go-ethereum/core/jvm/instructions/base"
import "github.com/ethereum/go-ethereum/core/jvm/rtda"

// Store int into local variable
type ISTORE struct{ base.Index8Instruction }

func (self *ISTORE) Execute(frame *rtda.Frame, gas uint64, contract interface{}) uint64 {
	_istore(frame, uint(self.Index))
    return 100
}

type ISTORE_0 struct{ base.NoOperandsInstruction }

func (self *ISTORE_0) Execute(frame *rtda.Frame, gas uint64, contract interface{}) uint64 {
	_istore(frame, 0)
    return 100
}

type ISTORE_1 struct{ base.NoOperandsInstruction }

func (self *ISTORE_1) Execute(frame *rtda.Frame, gas uint64, contract interface{}) uint64 {
	_istore(frame, 1)
    return 100
}

type ISTORE_2 struct{ base.NoOperandsInstruction }

func (self *ISTORE_2) Execute(frame *rtda.Frame, gas uint64, contract interface{}) uint64 {
	_istore(frame, 2)
    return 100
}

type ISTORE_3 struct{ base.NoOperandsInstruction }

func (self *ISTORE_3) Execute(frame *rtda.Frame, gas uint64, contract interface{}) uint64 {
	_istore(frame, 3)
    return 100
}

func _istore(frame *rtda.Frame, index uint) {
	val := frame.OperandStack().PopInt()
	frame.LocalVars().SetInt(index, val)
}
