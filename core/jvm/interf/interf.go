package interf

import "encoding/binary"

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
    "github.com/ethereum/go-ethereum/core/jvm/rtda/heap"
    "github.com/ethereum/go-ethereum/core/jvm/rtda"
)

// StateDB is an EVM database for full state querying.
type StateDB interface {
	CreateAccount(common.Address)

	SubBalance(common.Address, *big.Int)
	AddBalance(common.Address, *big.Int)
	GetBalance(common.Address) *big.Int

	GetNonce(common.Address) uint64
	SetNonce(common.Address, uint64)

	GetCodeHash(common.Address) common.Hash
	GetCode(common.Address) []byte
	SetCode(common.Address, []byte)
	GetCodeSize(common.Address) int

	AddRefund(uint64)
	GetRefund() uint64

	GetState(common.Address, common.Hash) common.Hash
	SetState(common.Address, common.Hash, common.Hash)

	Suicide(common.Address) bool
	HasSuicided(common.Address) bool

	// Exist reports whether the given account exists in state.
	// Notably this should also return true for suicided accounts.
	Exist(common.Address) bool
	// Empty returns whether the given account is empty. Empty
	// is defined according to EIP161 (balance = nonce = code = 0).
	Empty(common.Address) bool

	RevertToSnapshot(int)
	Snapshot() int

	AddLog(*types.Log)
	AddPreimage(common.Hash, []byte)

	ForEachStorage(common.Address, func(common.Hash, common.Hash) bool)
}

func Uint32ToBytes(u uint32) []byte {
    bs := make([]byte, 4)
    binary.LittleEndian.PutUint32(bs, u)
    return bs
}

func BytesToUint32(bs []byte) uint32 {
    return binary.LittleEndian.Uint32(bs)
}

type EVM interface {
    InternalCall(_caller interface{}, addr common.Address, methodName string, args []*heap.Object, gas uint64, value *big.Int, returnValueHandler func(*rtda.Frame, string)) (leftOverGas uint64, err error)
}

func AddressToObject(addr common.Address, classLoader *heap.ClassLoader) *heap.Object {
     class := classLoader.LoadClass("blockchain/types/Address")
     obj := class.NewObject()

     arrClass := classLoader.LoadClass("[B")
     arr := arrClass.NewArray(common.AddressLength)
     
     for _, field := range(class.Fields()) {
         if !field.IsStatic() && field.Name()=="bytes" && field.Descriptor()=="[B" {
             obj.Fields().SetRef(field.SlotId(), arr)
             break;
         }
     }

     for i:=0; i<common.AddressLength; i++ {
         arr.Bytes()[i] = int8(addr[i])
     }
     return obj
}

func GetAddressFromObject(obj *heap.Object) common.Address {
     var arr *heap.Object;

     for _, field := range(obj.Class().Fields()) {
         if !field.IsStatic() && field.Name()=="bytes" && field.Descriptor()=="[B" {
             arr = obj.Fields().GetRef(field.SlotId())
         }
     }

     var addr common.Address
     for i:=0; i<common.AddressLength; i++ {
         addr[i] = byte(arr.Bytes()[i])
     }

     return addr
}
