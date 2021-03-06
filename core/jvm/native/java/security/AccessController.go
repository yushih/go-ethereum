package security

import "github.com/ethereum/go-ethereum/core/jvm/instructions/base"
import "github.com/ethereum/go-ethereum/core/jvm/native"
import "github.com/ethereum/go-ethereum/core/jvm/rtda"

func init() {
	native.Register("java/security/AccessController", "doPrivileged", "(Ljava/security/PrivilegedAction;)Ljava/lang/Object;", doPrivileged)
	native.Register("java/security/AccessController", "doPrivileged", "(Ljava/security/PrivilegedExceptionAction;)Ljava/lang/Object;", doPrivileged)
	native.Register("java/security/AccessController", "getStackAccessControlContext", "()Ljava/security/AccessControlContext;", getStackAccessControlContext)
}

// @CallerSensitive
// public static native <T> T
//     doPrivileged(PrivilegedExceptionAction<T> action)
//     throws PrivilegedActionException;
// (Ljava/security/PrivilegedExceptionAction;)Ljava/lang/Object;

// @CallerSensitive
// public static native <T> T doPrivileged(PrivilegedAction<T> action);
// (Ljava/security/PrivilegedAction;)Ljava/lang/Object;
func doPrivileged(frame *rtda.Frame, gas uint64, contract interface{}, evm interface{}) {
	vars := frame.LocalVars()
	action := vars.GetRef(0)

	stack := frame.OperandStack()
	stack.PushRef(action)

	method := action.Class().GetInstanceMethod("run", "()Ljava/lang/Object;") // todo
	base.InvokeMethod(frame, method)
}

// private static native AccessControlContext getStackAccessControlContext();
// ()Ljava/security/AccessControlContext;
func getStackAccessControlContext(frame *rtda.Frame, gas uint64, contract interface{}, evm interface{}) {
	// todo
	frame.OperandStack().PushRef(nil)
}
