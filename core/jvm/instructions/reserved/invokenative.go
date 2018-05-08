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

func (self *INVOKE_NATIVE) Execute(frame *rtda.Frame, gas uint64, contract interface{}) uint64 {
	method := frame.Method()
	className := method.Class().Name()
	methodName := method.Name()
	methodDescriptor := method.Descriptor()

	nativeMethod := native.FindNativeMethod(className, methodName, methodDescriptor)
	if nativeMethod == nil {
		methodInfo := className + "." + methodName + methodDescriptor
		panic("java.lang.UnsatisfiedLinkError: " + methodInfo)
	}

  	nativeMethod(frame, gas, contract)

    return 100
}
