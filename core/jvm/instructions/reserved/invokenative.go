package reserved

import "github.com/ethereum/go-ethereum/core/jvm/instructions/base"
import "github.com/ethereum/go-ethereum/core/jvm/rtda"
import "github.com/ethereum/go-ethereum/core/jvm/native"
import _ "github.com/ethereum/go-ethereum/core/jvm/native/java/io"
import _ "github.com/ethereum/go-ethereum/core/jvm/native/java/lang"
import _ "github.com/ethereum/go-ethereum/core/jvm/native/java/security"
import _ "github.com/ethereum/go-ethereum/core/jvm/native/java/util/concurrent/atomic"
import _ "github.com/ethereum/go-ethereum/core/jvm/native/sun/io"
import _ "github.com/ethereum/go-ethereum/core/jvm/native/sun/misc"
import _ "github.com/ethereum/go-ethereum/core/jvm/native/sun/reflect"

// Invoke native method
type INVOKE_NATIVE struct{ base.NoOperandsInstruction }

func (self *INVOKE_NATIVE) Execute(frame *rtda.Frame, gas uint64, contract interface{}, evm interface{}) uint64 {
	method := frame.Method()
	className := method.Class().Name()
	methodName := method.Name()
	methodDescriptor := method.Descriptor()

	nativeMethod := native.FindNativeMethod(className, methodName, methodDescriptor)
	if nativeMethod == nil {
		methodInfo := className + "." + methodName + methodDescriptor
		panic("java.lang.UnsatisfiedLinkError: " + methodInfo)
	}

    gasLeft = nil
    err = nil

    nativeMethod(frame, gas, contract, evm)

    if gasLeft != nil {
        return gas - *gasLeft
    } else {
        return 100
    }
}

//normally we should have nativeMethod return these but currently there is 
//only one (blockchain.go/call) actually needs to so here is some hacking
var gasLeft *uint64
var err *error

func SetGasAndError(g uint64, e error) {
     gasLeft = &g
     err = &e
}