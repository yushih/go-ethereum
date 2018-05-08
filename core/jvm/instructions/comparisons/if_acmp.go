package comparisons

import "github.com/ethereum/go-ethereum/core/jvm/instructions/base"
import "github.com/ethereum/go-ethereum/core/jvm/rtda"

// Branch if reference comparison succeeds
type IF_ACMPEQ struct{ base.BranchInstruction }

func (self *IF_ACMPEQ) Execute(frame *rtda.Frame, gas uint64, contract interface{}, evm interface{}) uint64 {
	if _acmp(frame) {
		base.Branch(frame, self.Offset)
	}
    return 100
}

type IF_ACMPNE struct{ base.BranchInstruction }

func (self *IF_ACMPNE) Execute(frame *rtda.Frame, gas uint64, contract interface{}, evm interface{}) uint64 {
	if !_acmp(frame) {
		base.Branch(frame, self.Offset)
	}
    return 100
}

func _acmp(frame *rtda.Frame) bool {
	stack := frame.OperandStack()
	ref2 := stack.PopRef()
	ref1 := stack.PopRef()
	return ref1 == ref2 // todo
}
