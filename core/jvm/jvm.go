package jvm

import "fmt"
//import "strings"
import "encoding/binary"
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
     gasLeft, err := interpret(self.mainThread, false, gas)
     if err == nil {
         persistObjectGraph(obj, contractAddr, stateDB)
     }
     return gasLeft, err
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
     if err == nil {
         persistObjectGraph(obj, contractAddr, stateDB)
     }
     return nil, gasLeft, err
}

//todo optimize array storage

func persistObjectGraph(rootObj *heap.Object, contractAddr common.Address, stateDB interf.StateDB) {
     persisted := make(map[*heap.Object][]uint)

     persist(rootObj, []uint{}, false, contractAddr, stateDB, persisted)     
}

func write (path []uint, hash common.Hash, contractAddr common.Address, stateDB interf.StateDB) {
     fmt.Printf("---writing %v = %v\n", path, hash)
     stateDB.SetState(contractAddr, pathToHash(path), hash)
 }

func writeBytes (path []uint, bs []byte, contractAddr common.Address, stateDB interf.StateDB) {
     var h common.Hash
     binary.LittleEndian.PutUint16(h[0:2], uint16(len(bs)))
     size := min(len(bs), common.HashLength-2)
     copy(h[2:], bs[0:size])
     write(path, h, contractAddr, stateDB)
     bs = bs[size:]
     var i uint = 0
     for len(bs)>0 {
         size = min(len(bs), common.HashLength)
         var h common.Hash
         copy(h[:], bs[0:size])
         write(append(path, i), h, contractAddr, stateDB)
         bs = bs[size:]
         i += 1 
     } 
 }

func persist (obj *heap.Object, pathPrefix []uint, persistType bool, contractAddr common.Address, stateDB interf.StateDB, persisted map[*heap.Object][]uint) {
     if path, ok := persisted[obj]; ok {
         writeBytes(append(pathPrefix, Ref), pathToBytes(path), contractAddr, stateDB)
         return
     }
     persisted[obj] = pathPrefix

     if persistType {
         writeBytes(append(pathPrefix, TypeOrLength), []byte(obj.Class().Name()), contractAddr, stateDB)
     }

     if obj.Class().IsArray() {
         write(append(pathPrefix, TypeOrLength), intToHash(int(obj.ArrayLength())), contractAddr, stateDB)
         for i:=0; i<int(obj.ArrayLength()); i++ {
             path := append(pathPrefix, uint(i+1)) // 0 is type or length

             switch obj.Data().(type) {
             case []int8:
                 write(path, intToHash(int(obj.Data().([]int8)[i])), contractAddr, stateDB)
             case []int16:
                 write(path, intToHash(int(obj.Data().([]int16)[i])), contractAddr, stateDB)
             case []int32:
                 write(path, intToHash(int(obj.Data().([]int32)[i])), contractAddr, stateDB)
             case []int64:
                 write(path, intToHash(int(obj.Data().([]int64)[i])), contractAddr, stateDB)
             case []uint16:
                 write(path, intToHash(int(obj.Data().([]uint16)[i])), contractAddr, stateDB)
             case []float32:
                 write(path, floatToHash(obj.Data().([]float32)[i]), contractAddr, stateDB)
             case []float64:
                 write(path, doubleToHash(obj.Data().([]float64)[i]), contractAddr, stateDB)
             case []*heap.Object:
                 o := obj.Data().([]*heap.Object)[i]
                 if o != nil {
                     persist(o, path, o.Class().Name()!=obj.Class().Name()[1:], contractAddr, stateDB, persisted)
                 }
             default:
                //todo
             }
         }
     } else {
         slots := obj.Fields()

         for _, field := range(obj.Class().Fields()) {
             if field.IsStatic() {
                 continue
             }
             descriptor := field.Descriptor()
             slotId := field.SlotId()
             
             fmt.Printf("---persisting field %v (slot %v) %v for %v\n", field.Name(), slotId, descriptor, obj.Class().Name())

             path := append(pathPrefix, slotId+1) // 0 is type of length

             //todo optimize for Ljava/lang/String;

             switch descriptor[0] {             
             case 'Z', 'B', 'C', 'S', 'I':
                 write(path, intToHash(int(slots.GetInt(slotId))), contractAddr, stateDB)
             case 'F':
                 write(path, floatToHash(slots.GetFloat(slotId)), contractAddr, stateDB)
             case 'J':
                 write(path, intToHash(int(slots.GetLong(slotId))), contractAddr, stateDB)
             case 'D':
                 write(path, doubleToHash(slots.GetDouble(slotId)), contractAddr, stateDB)
             case 'L':
                 o := slots.GetRef(slotId)
                 if o != nil {
                     persist(o, path, ("L"+o.Class().Name()+";")!=descriptor, contractAddr, stateDB, persisted)
                 }
             case '[':
                 o := slots.GetRef(slotId)
                 if o != nil {
                     persist(o, path, false, contractAddr, stateDB, persisted)
                 }
             default:
             // todo
             }
         } // for fields
     } // if array
} 

const (
    TypeOrLength = 0
    Ref = math.MaxUint16
)

func intToHash(i int) common.Hash {
     bs := make([]byte, 8)
     binary.LittleEndian.PutUint64(bs, uint64(i))
     return common.BytesToHash(bs)
}

func hashToInt(h common.Hash) int {
     return int(binary.LittleEndian.Uint64(h.Bytes()[:common.HashLength-8]))
}

func floatToHash(f float32) common.Hash {
     bs := make([]byte, 4)
     binary.LittleEndian.PutUint32(bs, math.Float32bits(f))
     return common.BytesToHash(bs)     
}

func doubleToHash(d float64) common.Hash {
     bs := make([]byte, 8)
     binary.LittleEndian.PutUint64(bs, math.Float64bits(d))
     return common.BytesToHash(bs)
}

func pathToBytes(path []uint) []byte {
     b := []byte{}
     for _, u := range(path) {
        bs := make([]byte, 16) //???big enough? 
        n := binary.PutUvarint(bs, uint64(u))
        b = append(b, bs[0:n]...)
     }
     return b
}

func pathToHash(path []uint) common.Hash {
     b := pathToBytes(path)

     if len(b) <= common.HashLength {
        var h common.Hash
        copy(h[:], b)
        return h
     } else {
        panic("todo")
     }
}

func min(x, y int) int {
    if x<y {
        return x
    } else {
        return y
    }
}
