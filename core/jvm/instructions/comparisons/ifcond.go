package comparisons

import "github.com/ethereum/go-ethereum/core/jvm/instructions/base"
import "github.com/ethereum/go-ethereum/core/jvm/rtda"

// Branch if int comparison with zero succeeds
type IFEQ struct{ base.BranchInstruction }

func (self *IFEQ) Execute(frame *rtda.Frame, gas uint64) uint64 {
	val := frame.OperandStack().PopInt()
	if val == 0 {
		base.Branch(frame, self.Offset)
	}
    return 100
}

type IFNE struct{ base.BranchInstruction }

func (self *IFNE) Execute(frame *rtda.Frame, gas uint64) uint64 {
	val := frame.OperandStack().PopInt()
	if val != 0 {
		base.Branch(frame, self.Offset)
	}
    return 100
}

type IFLT struct{ base.BranchInstruction }

func (self *IFLT) Execute(frame *rtda.Frame, gas uint64) uint64 {
	val := frame.OperandStack().PopInt()
	if val < 0 {
		base.Branch(frame, self.Offset)
	}
    return 100
}

type IFLE struct{ base.BranchInstruction }

func (self *IFLE) Execute(frame *rtda.Frame, gas uint64) uint64 {
	val := frame.OperandStack().PopInt()
	if val <= 0 {
		base.Branch(frame, self.Offset)
	}
    return 100
}

type IFGT struct{ base.BranchInstruction }

func (self *IFGT) Execute(frame *rtda.Frame, gas uint64) uint64 {
	val := frame.OperandStack().PopInt()
	if val > 0 {
		base.Branch(frame, self.Offset)
	}
    return 100
}

type IFGE struct{ base.BranchInstruction }

func (self *IFGE) Execute(frame *rtda.Frame, gas uint64) uint64 {
	val := frame.OperandStack().PopInt()
	if val >= 0 {
		base.Branch(frame, self.Offset)
	}
    return 100
}
