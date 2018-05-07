package jvm

import (
    "math/big"

	"github.com/ethereum/go-ethereum/common"
)

//opcodes.go

// OpCode is an EVM opcode
type OpCode byte

func (op OpCode) IsPush() bool {
	return false
}

func (op OpCode) IsStaticJump() bool {
    return false
}

func (o OpCode) String() string {
    return "dummy"
}

//memory.go
type Memory struct {
	store       []byte
	lastGasCost uint64
}
// Len returns the length of the backing slice
func (m *Memory) Len() int {
	return len(m.store)
}

// Get returns offset + size as a new slice
func (self *Memory) Get(offset, size int64) (cpy []byte) {
	if size == 0 {
		return nil
	}

	if len(self.store) > int(offset) {
		cpy = make([]byte, size)
		copy(cpy, self.store[offset:offset+size])

		return
	}

	return
}

// GetPtr returns the offset + size
func (self *Memory) GetPtr(offset, size int64) []byte {
	if size == 0 {
		return nil
	}

	if len(self.store) > int(offset) {
		return self.store[offset : offset+size]
	}

	return nil
}

//stack.go

type Stack struct {
	data []*big.Int
}

func (st *Stack) Data() []*big.Int {
	return st.data
}

//logger.go
// StructLog is emitted to the EVM each cycle and lists information about the current internal state
// prior to the execution of the statement.
type StructLog struct {
	Pc         uint64                      `json:"pc"`
	Op         OpCode                      `json:"op"`
	Gas        uint64                      `json:"gas"`
	GasCost    uint64                      `json:"gasCost"`
	Memory     []byte                      `json:"memory"`
	MemorySize int                         `json:"memSize"`
	Stack      []*big.Int                  `json:"stack"`
	Storage    map[common.Hash]common.Hash `json:"-"`
	Depth      int                         `json:"depth"`
	Err        error                       `json:"-"`
}

// LogConfig are the configuration options for structured logger the EVM
type LogConfig struct {
	DisableMemory  bool // disable memory capture
	DisableStack   bool // disable stack capture
	DisableStorage bool // disable storage capture
	Limit          int  // maximum length of output, but zero means unlimited
}

