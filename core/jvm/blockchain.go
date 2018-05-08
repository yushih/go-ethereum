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
    native.Register("blockchain/Special", "msgValue", "()J", msgValue)
    native.Register("blockchain/Special", "txOrigin", "()Lblockchain/types/Address;", txOrigin)
    native.Register("blockchain/Special", "gasPrice", "()J", gasPrice)
    native.Register("blockcahin/Special", "thisAddr", "()Lblockchain/types/Address;", thisAddr)
}

func thisAddr(frame *rtda.Frame, gas uint64, contract interface{}, evm interface{}) {
     v := contract.(*Contract).Address()
     frame.OperandStack().PushRef(addressToObject(v, frame.Thread().ClassLoader()))
}

func msgValue(frame *rtda.Frame, gas uint64, contract interface{}, evm interface{}) {
     //fixme
     v := contract.(*Contract).Value()
     vi64 := v.Int64()
     frame.OperandStack().PushLong(vi64)
}

func gasPrice(frame *rtda.Frame, gas uint64, contract interface{}, evm interface{}) {
     //fixme
     v := evm.(*EVM).GasPrice
     vi64 := v.Int64()
     frame.OperandStack().PushLong(vi64)
}

func gasLeft(frame *rtda.Frame, gas uint64, contract interface{}, evm interface{}) {
     frame.OperandStack().PushLong(int64(gas))
}

func txOrigin(frame *rtda.Frame, gas uint64, contract interface{}, evm interface{}) {
     v := evm.(*EVM).Origin
     frame.OperandStack().PushRef(addressToObject(v, frame.Thread().ClassLoader()))
}

func msgSender(frame *rtda.Frame, gas uint64, contract interface{}, evm interface{}) {
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