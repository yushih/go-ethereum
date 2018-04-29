package jvm

//import "fmt"
//import "strings"
import "github.com/ethereum/go-ethereum/core/jvm/classpath"
import "github.com/ethereum/go-ethereum/core/jvm/instructions/base"
import "github.com/ethereum/go-ethereum/core/jvm/rtda"
import "github.com/ethereum/go-ethereum/core/jvm/rtda/heap"

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
	interpret(self.mainThread, false) //todo
}

func (self *JVM) Deploy(contractCode []byte, stateDB StateDB) {
     class := self.classLoader.LoadClassFromBytes(contractCode)
     obj := class.NewObject()
     
     method := class.GetConstructor("()V") 
     frame := self.mainThread.NewFrame(method)
     self.mainThread.PushFrame(frame)
     frame.LocalVars().SetRef(0, obj)
     interpret(self.mainThread, false)
     
}

func (self *JVM) ExecContract(contractCode []byte, input []byte) ([]byte, error) {
     class := self.classLoader.LoadClassFromBytes(contractCode)
     methodName := string(input) //todo
     method := class.GetStaticMethod(methodName, "()V") //todo
     
     frame := self.mainThread.NewFrame(method)
     self.mainThread.PushFrame(frame)
     interpret(self.mainThread, false)
     return nil, nil
}
