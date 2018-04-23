package control

import "github.com/ethereum/go-ethereum/core/jvm/instructions/base"
import "github.com/ethereum/go-ethereum/core/jvm/rtda"

// Branch always
type GOTO struct{ base.BranchInstruction }

func (self *GOTO) Execute(frame *rtda.Frame) {
	base.Branch(frame, self.Offset)
}
