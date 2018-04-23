package rtda

import "github.com/ethereum/go-ethereum/core/jvm/rtda/heap"

type Slot struct {
	num int32
	ref *heap.Object
}
