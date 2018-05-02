package jvm

//import "fmt"
//import "strings"
import "math"

import "github.com/ethereum/go-ethereum/common"

import "github.com/ethereum/go-ethereum/core/jvm/classpath"
import "github.com/ethereum/go-ethereum/core/jvm/instructions/base"
import "github.com/ethereum/go-ethereum/core/jvm/rtda"
import "github.com/ethereum/go-ethereum/core/jvm/rtda/heap"
import "github.com/ethereum/go-ethereum/core/jvm/interf"

type JVM struct {
	classLoader *heap.ClassLoader
	mainThread  *rtda.Thread
}

func newJVM() *JVM {
	cp := classpath.Parse("", "") //todo
	classLoader := heap.NewClassLoader(cp, true) //todo
	jvm:= &JVM{
		classLoader: classLoader,
		mainThread:  rtda.NewThread(),
	}
    jvm.initVM()
    return jvm
}

var jvm = newJVM()

func GetJVM() *JVM {
     return jvm
}

func (self *JVM) initVM() {
	vmClass := self.classLoader.LoadClass("sun/misc/VM")
	base.InitClass(self.mainThread, vmClass)
	interpret(self.mainThread, false, math.MaxUint64) //todo
}

func (self *JVM) Deploy(contractCode []byte, contractAddr common.Address, stateDB interf.StateDB, gas uint64) (uint64, error) {
     class := self.classLoader.LoadClassFromBytes(contractCode)
     obj := class.NewObject()
     
     method := class.GetConstructor("()V") 
     frame := self.mainThread.NewFrame(method)
     self.mainThread.PushFrame(frame)
     frame.LocalVars().SetRef(0, obj)
     return interpret(self.mainThread, false, gas)
}

func (self *JVM) ExecContract(contractCode []byte, input []byte, contractAddr common.Address, stateDB interf.StateDB, gas uint64) ([]byte, uint64, error) {
     class := self.classLoader.LoadClassFromBytes(contractCode)
     methodName := string(input) //todo
     method := class.GetInstanceMethod(methodName, "()V") //todo
     obj := class.NewObject()

     frame := self.mainThread.NewFrame(method)
     self.mainThread.PushFrame(frame)
     frame.LocalVars().SetRef(0, obj)
     gasLeft, err := interpret(self.mainThread, false, gas)
     return nil, gasLeft, err
}
