package jvm

import (
    "math/big"

    "github.com/ethereum/go-ethereum/core/jvm/native"
    "github.com/ethereum/go-ethereum/core/jvm/rtda"
    "github.com/ethereum/go-ethereum/core/jvm/rtda/heap"
    "github.com/ethereum/go-ethereum/core/jvm/instructions/reserved"
    "github.com/ethereum/go-ethereum/core/jvm/interf"
)


func init() {
    native.Register("blockchain/Special", "gasLeft", "()J", gasLeft)
    native.Register("blockchain/Special", "msgSender", "()Lblockchain/types/Address;", msgSender)
    native.Register("blockchain/Special", "msgValue", "()J", msgValue)
    native.Register("blockchain/Special", "txOrigin", "()Lblockchain/types/Address;", txOrigin)
    native.Register("blockchain/Special", "gasPrice", "()J", gasPrice)
    native.Register("blockchain/Special", "thisAddr", "()Lblockchain/types/Address;", thisAddr)

    native.Register("blockchain/types/Address", "balance", "()J", balance)
    native.Register("blockchain/types/Address", "transfer", "(J)V", transfer)
    native.Register("blockchain/types/Address", "call", "(Ljava/lang/String;J[Ljava/lang/Object;)Ljava/lang/Object;", call)
}

func call(frame *rtda.Frame, gas uint64, contract interface{}, evm interface{}) {
     addr := interf.GetAddressFromObject(frame.LocalVars().GetThis())
     methodName := heap.GoString(frame.LocalVars().GetRef(1))
     value := big.NewInt(frame.LocalVars().GetLong(2)) //fixme
     // previous long arg consume 2 slots
     args := frame.LocalVars().GetRef(4).Refs()
     classLoader := frame.Thread().ClassLoader()
     gasLeft, err := evm.(interf.EVM).InternalCall(contract, addr, methodName, args, gas, value, func (bogusFrame *rtda.Frame, returnTypeDescriptor string) {
         var ret *heap.Object
         switch returnTypeDescriptor[0] {
         case 'V':
             ret = nil
         case 'Z':
             ret = heap.BoxBool(bogusFrame.OperandStack().PopInt(), classLoader)
         case 'B':
             ret = heap.BoxByte(bogusFrame.OperandStack().PopInt(), classLoader)
         case 'C':
             ret = heap.BoxChar(bogusFrame.OperandStack().PopInt(), classLoader)
         case 'S':
             ret = heap.BoxShort(bogusFrame.OperandStack().PopInt(), classLoader)
         case 'I':
             ret = heap.BoxInt(bogusFrame.OperandStack().PopInt(), classLoader)
         case 'F':
             ret = heap.BoxFloat(bogusFrame.OperandStack().PopFloat(), classLoader)
         case 'J':
             ret = heap.BoxLong(bogusFrame.OperandStack().PopLong(), classLoader)
         case 'D':
             ret = heap.BoxDouble(bogusFrame.OperandStack().PopDouble(), classLoader)
         case 'L', '[':
             ret = bogusFrame.OperandStack().PopRef()
         }
         frame.OperandStack().PushRef(ret)
     })
     //hack see invokenative.go
     reserved.SetGasAndError(gasLeft, err)
}


func balance(frame *rtda.Frame, gas uint64, contract interface{}, evm interface{}) {
     addrObj := frame.LocalVars().GetRef(0)
     addr := interf.GetAddressFromObject(addrObj)
     v := evm.(*EVM).StateDB.GetBalance(addr)

     vi64 := v.Int64()
     frame.OperandStack().PushLong(vi64)
}

func transfer(frame *rtda.Frame, gas uint64, contract interface{}, evm interface{}) {
}

func thisAddr(frame *rtda.Frame, gas uint64, contract interface{}, evm interface{}) {
     v := contract.(*Contract).Address()
     frame.OperandStack().PushRef(interf.AddressToObject(v, frame.Thread().ClassLoader()))
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
     frame.OperandStack().PushRef(interf.AddressToObject(v, frame.Thread().ClassLoader()))
}

func msgSender(frame *rtda.Frame, gas uint64, contract interface{}, evm interface{}) {
     v := contract.(*Contract).Caller()
     frame.OperandStack().PushRef(interf.AddressToObject(v, frame.Thread().ClassLoader()))
}

