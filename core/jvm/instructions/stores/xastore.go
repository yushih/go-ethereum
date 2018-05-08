package stores

import "github.com/ethereum/go-ethereum/core/jvm/instructions/base"
import "github.com/ethereum/go-ethereum/core/jvm/rtda"
import "github.com/ethereum/go-ethereum/core/jvm/rtda/heap"

// Store into reference array
type AASTORE struct{ base.NoOperandsInstruction }

func (self *AASTORE) Execute(frame *rtda.Frame, gas uint64, contract interface{}) uint64 {
	stack := frame.OperandStack()
	ref := stack.PopRef()
	index := stack.PopInt()
	arrRef := stack.PopRef()

	checkNotNil(arrRef)
	refs := arrRef.Refs()
	checkIndex(len(refs), index)
	refs[index] = ref
    return 100
}

// Store into byte or boolean array
type BASTORE struct{ base.NoOperandsInstruction }

func (self *BASTORE) Execute(frame *rtda.Frame, gas uint64, contract interface{}) uint64 {
	stack := frame.OperandStack()
	val := stack.PopInt()
	index := stack.PopInt()
	arrRef := stack.PopRef()

	checkNotNil(arrRef)
	bytes := arrRef.Bytes()
	checkIndex(len(bytes), index)
	bytes[index] = int8(val)
    return 100
}

// Store into char array
type CASTORE struct{ base.NoOperandsInstruction }

func (self *CASTORE) Execute(frame *rtda.Frame, gas uint64, contract interface{}) uint64 {
	stack := frame.OperandStack()
	val := stack.PopInt()
	index := stack.PopInt()
	arrRef := stack.PopRef()

	checkNotNil(arrRef)
	chars := arrRef.Chars()
	checkIndex(len(chars), index)
	chars[index] = uint16(val)
    return 100
}

// Store into double array
type DASTORE struct{ base.NoOperandsInstruction }

func (self *DASTORE) Execute(frame *rtda.Frame, gas uint64, contract interface{}) uint64 {
	stack := frame.OperandStack()
	val := stack.PopDouble()
	index := stack.PopInt()
	arrRef := stack.PopRef()

	checkNotNil(arrRef)
	doubles := arrRef.Doubles()
	checkIndex(len(doubles), index)
	doubles[index] = float64(val)
    return 100
}

// Store into float array
type FASTORE struct{ base.NoOperandsInstruction }

func (self *FASTORE) Execute(frame *rtda.Frame, gas uint64, contract interface{}) uint64 {
	stack := frame.OperandStack()
	val := stack.PopFloat()
	index := stack.PopInt()
	arrRef := stack.PopRef()

	checkNotNil(arrRef)
	floats := arrRef.Floats()
	checkIndex(len(floats), index)
	floats[index] = float32(val)
    return 100
}

// Store into int array
type IASTORE struct{ base.NoOperandsInstruction }

func (self *IASTORE) Execute(frame *rtda.Frame, gas uint64, contract interface{}) uint64 {
	stack := frame.OperandStack()
	val := stack.PopInt()
	index := stack.PopInt()
	arrRef := stack.PopRef()

	checkNotNil(arrRef)
	ints := arrRef.Ints()
	checkIndex(len(ints), index)
	ints[index] = int32(val)
    return 100
}

// Store into long array
type LASTORE struct{ base.NoOperandsInstruction }

func (self *LASTORE) Execute(frame *rtda.Frame, gas uint64, contract interface{}) uint64 {
	stack := frame.OperandStack()
	val := stack.PopLong()
	index := stack.PopInt()
	arrRef := stack.PopRef()

	checkNotNil(arrRef)
	longs := arrRef.Longs()
	checkIndex(len(longs), index)
	longs[index] = int64(val)
    return 100
}

// Store into short array
type SASTORE struct{ base.NoOperandsInstruction }

func (self *SASTORE) Execute(frame *rtda.Frame, gas uint64, contract interface{}) uint64 {
	stack := frame.OperandStack()
	val := stack.PopInt()
	index := stack.PopInt()
	arrRef := stack.PopRef()

	checkNotNil(arrRef)
	shorts := arrRef.Shorts()
	checkIndex(len(shorts), index)
	shorts[index] = int16(val)
    return 100
}

func checkNotNil(ref *heap.Object) {
	if ref == nil {
		panic("java.lang.NullPointerException")
	}
}
func checkIndex(arrLen int, index int32) {
	if index < 0 || index >= int32(arrLen) {
		panic("ArrayIndexOutOfBoundsException")
	}
}
