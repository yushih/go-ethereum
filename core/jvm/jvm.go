package jvm

import "fmt"
import "strings"
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
		mainThread:  rtda.NewThread(classLoader),
	}
    jvm.initVM()
    return jvm
}

var jvm = newJVM()

func getJVM() *JVM {
     return jvm
}

func (self *JVM) initVM() {
	vmClass := self.classLoader.LoadClass("sun/misc/VM")
	base.InitClass(self.mainThread, vmClass)
	interpret(self.mainThread, false, math.MaxUint64, nil) //todo
}

func (self *JVM) deploy(contractCode []byte, contractAddr common.Address, stateDB interf.StateDB, contract *Contract) (uint64, error) {
     gas := contract.Gas
     class := self.classLoader.LoadClassFromBytes(contractCode)
     obj := class.NewObject()
     
     method := class.GetConstructor("()V") 
     frame := self.mainThread.NewFrame(method)
     self.mainThread.PushFrame(frame)
     frame.LocalVars().SetRef(0, obj)
     gasLeft, err := interpret(self.mainThread, false, gas, contract)
     if err == nil {
         persistObjectGraph(obj, contractAddr, stateDB)
     }
     return gasLeft, err
}

func (self *JVM) execContract(contractCode []byte, input []byte, contractAddr common.Address, stateDB interf.StateDB, contract *Contract) ([]byte, uint64, error) {
     gas := contract.Gas
     class := self.classLoader.LoadClassFromBytes(contractCode)
     methodName := string(input) //todo
     method := class.GetInstanceMethod(methodName, "()V") //todo
     obj := class.NewObject()
     reincarnateObject(obj, contractAddr, stateDB)

     frame := self.mainThread.NewFrame(method)
     self.mainThread.PushFrame(frame)
     frame.LocalVars().SetRef(0, obj)
     gasLeft, err := interpret(self.mainThread, false, gas, contract)
     if err == nil {
         persistObjectGraph(obj, contractAddr, stateDB)
     }
     return nil, gasLeft, err
}

//todo optimize array storage
const (
    TypeEntry = 0 
    ArrayLengthEntry = 1
    NullSignal = 2
    RefEntry = 3
    SlotIndexOffset = RefEntry+1
)

func persistObjectGraph(rootObj *heap.Object, contractAddr common.Address, stateDB interf.StateDB) {
     persisted := make(map[*heap.Object][]uint)

     persist(rootObj, []uint{}, false, contractAddr, stateDB, persisted)     
}

func persist (obj *heap.Object, pathPrefix []uint, persistType bool, contractAddr common.Address, stateDB interf.StateDB, persisted map[*heap.Object][]uint) {
     if obj == nil {
         writeBytes(append(pathPrefix, NullSignal), []byte{1}, contractAddr, stateDB)
         return
     }
     if path, ok := persisted[obj]; ok {
         writeBytes(append(pathPrefix, RefEntry), pathToBytes(path), contractAddr, stateDB)
         return
     }
     persisted[obj] = pathPrefix

     if persistType {
         writeBytes(append(pathPrefix, TypeEntry), []byte(obj.Class().Name()), contractAddr, stateDB)
     }

     if obj.Class().IsArray() {
         write(append(pathPrefix, ArrayLengthEntry), intToHash(int(obj.ArrayLength())), contractAddr, stateDB)
         var i uint
         for i=0; i<uint(obj.ArrayLength()); i++ {
             path := append(pathPrefix, i+SlotIndexOffset) 

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
                 elem := obj.Data().([]*heap.Object)[i]
                 //todo: better persist type judegment
                 persist(elem, path, true, contractAddr, stateDB, persisted)
             default:
                //todo should not happen
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
             
             if (descriptor[0]=='L' || descriptor[0]=='[') && slots.GetRef(slotId)!=nil {
                 fmt.Printf("---%v is persisting field %v (slot %v) type %v with %v\n", obj.Class().Name(), field.Name(), slotId, descriptor, slots.GetRef(slotId).Class().Name())
             } else {
                 fmt.Printf("---%v is persisting field %v (slot %v) type %v\n", obj.Class().Name(), field.Name(), slotId, descriptor)
             }

             path := append(pathPrefix, slotId+SlotIndexOffset) 

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
                 var persistType bool
                 if o == nil {
                     persistType = false // actually not relevant
                 } else {
                     persistType = "L"+o.Class().Name()+";" != descriptor
                 }
                 persist(o, path, persistType, contractAddr, stateDB, persisted)
             case '[':
                 o := slots.GetRef(slotId)
                 var persistType bool
                 if o == nil {
                     persistType = false // actually not relevant
                 } else {
                     persistType = o.Class().Name() != descriptor
                 }
                 persist(o, path, persistType, contractAddr, stateDB, persisted)
             default:
             // todo
             }
         } // for fields
     } // if array
} 

func reincarnateObject(obj *heap.Object, contractAddr common.Address, stateDB interf.StateDB) {
     pool := make(map[string]*heap.Object) // actually (path []uint] => *heap.Object
     pool[""] = obj
     reincarnate(obj, []uint{}, contractAddr, stateDB, pool)
}

func reincarnate(obj *heap.Object, pathPrefix []uint, contractAddr common.Address, stateDB interf.StateDB, pool map[string]*heap.Object) {
     if obj.Class().IsArray() {
         var i uint
         for i=0; i<uint(obj.ArrayLength()); i++ {
             path := append(pathPrefix, i+SlotIndexOffset)
             switch obj.Data().(type) {
             case []int8:
                 v := hashToInt(read(path, contractAddr, stateDB))
                 obj.Bytes()[i] = int8(v)
             case []int16:
                 v := hashToInt(read(path, contractAddr, stateDB))
                 obj.Shorts()[i] = int16(v)
             case []int32:
                 v := hashToInt(read(path, contractAddr, stateDB))
                 obj.Ints()[i] = int32(v)
             case []int64:
                 v := hashToInt(read(path, contractAddr, stateDB))
                 obj.Longs()[i] = int64(v)
             case []uint16:
                 v := hashToInt(read(path, contractAddr, stateDB))
                 obj.Chars()[i] = uint16(v)
             case []float32:
                 v := hashToFloat(read(path, contractAddr, stateDB))
                 obj.Floats()[i] = v
             case []float64:
                 v := hashToDouble(read(path, contractAddr, stateDB))
                 obj.Doubles()[i] = v
             case []*heap.Object:
                 elem := loadObjectOrArray(path, obj.Class().Loader(), obj.Class().Name()[1:len(obj.Class().Name())-1], contractAddr, stateDB, pool)
                 obj.Refs()[i] = elem
             default:
                //todo should not happen
             }
         } // for array
     } else {
         slots := obj.Fields()

         for _, field := range(obj.Class().Fields()) {
             if field.IsStatic() {
                 continue
             }
             descriptor := field.Descriptor()
             slotId := field.SlotId()
             
             path := append(pathPrefix, slotId+SlotIndexOffset)

             fmt.Printf("---loading field %v (slot %v) %v from %v\n", field.Name(), slotId, descriptor, path)
             switch descriptor[0] {
             case 'Z', 'B', 'C', 'S', 'I':
                 bs := read(path, contractAddr, stateDB)
                 v := hashToInt(bs)
                 slots.SetInt(slotId, int32(v))
             case 'F':
                  bs := read(path, contractAddr, stateDB)
                  v := hashToFloat(bs)
                  slots.SetFloat(slotId, v)
              case 'J':
                  bs := read(path, contractAddr, stateDB)
                  v := hashToInt(bs)
                  slots.SetLong(slotId, int64(v))
              case 'D':
                  bs := read(path, contractAddr, stateDB)
                  v := hashToDouble(bs)
                  slots.SetDouble(slotId, v)
              case 'L':
                  slots.SetRef(slotId, loadObjectOrArray(path, obj.Class().Loader(), descriptor[1:len(descriptor)-1], contractAddr, stateDB, pool))
              case '[':
                  slots.SetRef(slotId, loadObjectOrArray(path, obj.Class().Loader(), descriptor, contractAddr, stateDB, pool))
              } // switch
          } // for fields
     } // if array
}

func loadObjectOrArray(path []uint, classLoader *heap.ClassLoader, descriptor string, contractAddr common.Address, stateDB interf.StateDB, pool map[string]*heap.Object) *heap.Object {
      isArray := descriptor[0]=='['
      var bs []byte
      if !allZeros(read(append(path, NullSignal), contractAddr, stateDB)) {
          return nil
      } else if bs = readBytes(append(path, RefEntry), contractAddr, stateDB); !allZeros(bs) {
          return pool[pathToString(bytesToPath(bs))]
      } else {
          bs = readBytes(append(path, TypeEntry), contractAddr, stateDB)
          var t string
          if allZeros(bs) {
              t = descriptor
          } else {
              t = string(bs)
          }
          class := classLoader.LoadClass(t)
          var o *heap.Object
          if isArray {
              count := uint(hashToInt(read(append(path, ArrayLengthEntry), contractAddr, stateDB)))
              o = class.NewArray(count)
          } else {
              o = class.NewObject()
          }
          pool[pathToString(path)] = o
          reincarnate(o, path, contractAddr, stateDB, pool)
          return o
     }
}

func read(path []uint, contractAddr common.Address, stateDB interf.StateDB) (ret []byte) {
     ret = stateDB.GetState(contractAddr, pathToHash(path)).Bytes()
     fmt.Printf("----read %v = %v\n", path, ret)
     return
}

func write(path []uint, hash common.Hash, contractAddr common.Address, stateDB interf.StateDB) {
     fmt.Printf("---writing %v = %v\n", path, hash)
     stateDB.SetState(contractAddr, pathToHash(path), hash)
}

func readBytes(path []uint, contractAddr common.Address, stateDB interf.StateDB) (ret []byte) {
     bs := read(path, contractAddr, stateDB)
     left := int(binary.LittleEndian.Uint16(bs[0:2]))
     bs = bs[2:]
     var i uint = 0
     for left>0 {
         l := min(left, len(bs))
         ret = append(ret, bs[0:l]...)
         left -= l
         bs = read(append(path, i), contractAddr, stateDB)
         i += 1
     } 
     return
}

func writeBytes(path []uint, bs []byte, contractAddr common.Address, stateDB interf.StateDB) {
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

func intToHash(i int) common.Hash {
     bs := make([]byte, 8)
     binary.LittleEndian.PutUint64(bs, uint64(i))
     return common.BytesToHash(bs)
}

func hashToInt(h []byte) int {
     return int(binary.LittleEndian.Uint64(h[common.HashLength-8:]))
}

func floatToHash(f float32) common.Hash {
     bs := make([]byte, 4)
     binary.LittleEndian.PutUint32(bs, math.Float32bits(f))
     return common.BytesToHash(bs)     
}

func hashToFloat(h []byte) float32 {
     return math.Float32frombits(binary.LittleEndian.Uint32(h[common.HashLength-4:]))
}

func hashToDouble(h []byte) float64 {
     return math.Float64frombits(binary.LittleEndian.Uint64(h[common.HashLength-8:]))
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

func bytesToPath(bs []byte) []uint {
     path := []uint{}
     for len(bs) > 0 {
         //todo: handle exception
         u, n := binary.Uvarint(bs)
         path = append(path, uint(u))
         bs = bs[n:]
     }
     return path
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

func allZeros(bs []byte) bool {
     for _, b := range(bs) {
         if b != 0 {
             return false
         }
     }
     return true
}

func pathToString(path []uint) string {
     ss := make([]string, len(path))
     for i, u := range(path) {
         ss[i] = string(u)
     }
     return strings.Join(ss, ",")
}