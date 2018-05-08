package jvm

import (
    "github.com/ethereum/go-ethereum/core/jvm/native"
    "github.com/ethereum/go-ethereum/core/jvm/rtda"
    "github.com/ethereum/go-ethereum/core/jvm/rtda/heap"
	"github.com/ethereum/go-ethereum/common"
)

func init() {
    native.Register("blockchain/Special", "gasLeft", "()J", gasLeft)
    native.Register("blockchain/Special", "msgSender", "()Lblockchain/types/Address;", msgSender)
}

func gasLeft(frame *rtda.Frame, gas uint64, contract interface{}) {
     frame.OperandStack().PushLong(int64(gas))
}

func msgSender(frame *rtda.Frame, gas uint64, contract interface{}) {
     addrObj := addressToObject(contract.(*Contract).Caller(), frame.Thread().ClassLoader())
     frame.OperandStack().PushRef(addrObj)
}

func addressToObject(addr common.Address, classLoader *heap.ClassLoader) *heap.Object {
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