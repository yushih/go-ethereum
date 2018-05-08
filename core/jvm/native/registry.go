package native

import "github.com/ethereum/go-ethereum/core/jvm/rtda"

type NativeMethod func(frame *rtda.Frame, gas uint64, contract interface{})

var registry = map[string]NativeMethod{}

func emptyNativeMethod(frame *rtda.Frame, gas uint64, contract interface{}) {
	// do nothing
}

func Register(className, methodName, methodDescriptor string, method NativeMethod) {
	key := className + "~" + methodName + "~" + methodDescriptor
	registry[key] = method
}

func FindNativeMethod(className, methodName, methodDescriptor string) NativeMethod {
	key := className + "~" + methodName + "~" + methodDescriptor
	if method, ok := registry[key]; ok {
		return method
	}
	if methodDescriptor == "()V" {
		if methodName == "registerNatives" || methodName == "initIDs" {
			return emptyNativeMethod
		}
	}
	return nil
}
