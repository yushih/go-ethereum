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
     return &JVM{
         classLoader: jvm.classLoader,
         mainThread: rtda.NewThread(jvm.classLoader),
     }
}

func (self *JVM) initVM() {
	vmClass := self.classLoader.LoadClass("sun/misc/VM")
	base.InitClass(self.mainThread, vmClass)
	interpret(self.mainThread, false, nil, nil)
}

func (self *JVM) deploy(contractCode []byte, contractAddr common.Address, stateDB interf.StateDB, contract *Contract, evm *EVM) (uint64, error) {
     class := self.classLoader.LoadClassFromBytes(contractCode)
     obj := class.NewObject()
     
     method := class.GetConstructor("()V") 
     frame := self.mainThread.NewFrame(method)
     self.mainThread.PushFrame(frame)
     frame.LocalVars().SetRef(0, obj)
     gasLeft, err := interpret(self.mainThread, false, contract, evm)
     if err == nil {
         persistObjectGraph(obj, contractAddr, stateDB)
     }
     return gasLeft, err
}

func (self *JVM) execContract(input []byte, contractAddr common.Address, stateDB interf.StateDB, contract *Contract, evm *EVM) ([]byte, uint64, error) {
     class := self.classLoader.LoadClassFromBytes(contract.Code)
     var methodName string
     var suc bool
     if methodName, input, suc = readString(input); !suc {
         return nil, contract.Gas, ErrInput
     }
     method := class.GetPublicInstanceMethodByName(methodName);
     if method == nil {
         return nil, contract.Gas, ErrMethodNotFound
     }
     obj := class.NewObject()
     reincarnateObject(obj, contractAddr, stateDB)

     // this frame will accept the return value
     bogusFrame := rtda.NewBogusFrame()
     self.mainThread.PushFrame(bogusFrame)

     frame := self.mainThread.NewFrame(method)
     self.mainThread.PushFrame(frame)
     frame.LocalVars().SetRef(0, obj)

     for i, descriptor := range(method.ParameterTypeDescriptors()) {
         switch descriptor[0] {
         case 'I':
            if len(input) < 4 {
                return nil, contract.Gas, ErrInput
            }
            v := int32(binary.LittleEndian.Uint32(input[0:4]))
            input = input[4:]
            frame.LocalVars().SetInt(uint(i+1), v)
         case 'B', 'Z':
            if len(input) < 1 {
                return nil, contract.Gas, ErrInput
            }
            v := int32(input[0])
            input = input[1:]
            frame.LocalVars().SetInt(uint(i+1), v)
         case 'C', 'S':
            if len(input) < 2 {
                return nil, contract.Gas, ErrInput
            }
            v := int32(binary.LittleEndian.Uint16(input[0:2]))
            input = input[2:]
            frame.LocalVars().SetInt(uint(i+1), v)
         case 'F':
            if len(input) < 4 {
                return nil, contract.Gas, ErrInput
            }
            v := math.Float32frombits(binary.LittleEndian.Uint32(input[0:4]))
            input = input[4:]
            frame.LocalVars().SetFloat(uint(i+1), v)
         case 'D':
            if len(input) < 8 {
                return nil, contract.Gas, ErrInput
            }
            v := math.Float64frombits(binary.LittleEndian.Uint64(input[0:8]))
            input = input[8:]
            frame.LocalVars().SetDouble(uint(i+1), v)
         case 'J':
            if len(input) < 8 {
                return nil, contract.Gas, ErrInput
            }
            v := int64(binary.LittleEndian.Uint64(input[0:8]))
            input = input[8:]
            frame.LocalVars().SetLong(uint(i+1), v)
         case '[', 'L':
             //todo support arrays and object (by class name and ctor args)
             switch descriptor {
             case "Lblockchain/types/Address;":
                 if len(input) < common.AddressLength {
                     return nil, contract.Gas, ErrInput
                 }
                 addr := common.BytesToAddress(input[0:common.AddressLength])
                 frame.LocalVars().SetRef(uint(i+1), interf.AddressToObject(addr, self.classLoader))
                 input = input[common.AddressLength:]
             case "Ljava/lang/String;":
                 var s string
                 if s, input, suc = readString(input); !suc {
                     return nil, contract.Gas, ErrInput
                  }
                  frame.LocalVars().SetRef(uint(i+1), heap.JString(self.classLoader, s))
             }
         }// switch descriptor[0]
     }
     if len(input) !=0 {
         return nil, contract.Gas, ErrInput
     }

     gasLeft, err := interpret(self.mainThread, false, contract, evm)
     if err == nil {
         persistObjectGraph(obj, contractAddr, stateDB)
     } else {
         return nil, gasLeft, err
     }

     var ret []byte
     switch method.ReturnTypeDescriptor()[0] {
     case 'V':
         ret = nil
     case 'Z', 'B', 'C', 'S', 'I':
         v := bogusFrame.OperandStack().PopInt()
         ret = make([]byte, 4)
         binary.LittleEndian.PutUint32(ret, uint32(v))
     case 'F':
         v := bogusFrame.OperandStack().PopFloat()
         ret = make([]byte, 4)
         binary.LittleEndian.PutUint32(ret, math.Float32bits(v))
     case 'J':
         v := bogusFrame.OperandStack().PopLong()
         ret = make([]byte, 8)
         binary.LittleEndian.PutUint64(ret, uint64(v))
     case 'D':
         v := bogusFrame.OperandStack().PopDouble()
         ret = make([]byte, 8)
         binary.LittleEndian.PutUint64(ret, math.Float64bits(v))
     case 'L':
         o := bogusFrame.OperandStack().PopRef()
         if o!=nil && o.Class().Name()=="java/lang/String" {
             ret = []byte(heap.GoString(o))
         }
     case '[':
         //todo
     default:
     }
     return ret, gasLeft, err
}

func readString(bs []byte) (string, []byte, bool) {
     for index, runeValue := range string(bs) {
         if runeValue == 0 {
             return string(bs[0:index]), bs[index+1:], true
         }
     }
     return "", nil, false
}

func (self *JVM) internalExecContract(contractCode []byte, methodName string, args []*heap.Object, contractAddr common.Address, stateDB interf.StateDB, contract *Contract, evm *EVM, returnValueHandler func(*rtda.Frame, string))  (uint64, error) {

     class := self.classLoader.LoadClassFromBytes(contractCode)
     method := class.GetPublicInstanceMethodByName(methodName);
     obj := class.NewObject()
     reincarnateObject(obj, contractAddr, stateDB)

     // this frame will accept the return value
     bogusFrame := rtda.NewBogusFrame()
     self.mainThread.PushFrame(bogusFrame)

     frame := self.mainThread.NewFrame(method)
     self.mainThread.PushFrame(frame)
     frame.LocalVars().SetRef(0, obj)

     if len(method.ParameterTypeDescriptors()) != len(args) {
         return contract.Gas, ErrParameter
     }

     for _i, td := range(method.ParameterTypeDescriptors()) {
         i := uint(_i)
         switch td[0] {
         case 'I':
             o := args[i]
             if o==nil || o.Class().Name() != "java/lang/Integer" {
                 return contract.Gas, ErrParameter
             }
             v := heap.UnboxInt(o)
             frame.LocalVars().SetInt(i+1, v)
         case 'B':
             o := args[i]
             if o==nil || o.Class().Name() != "java/lang/Byte" {
                 return contract.Gas, ErrParameter
             }
             v := heap.UnboxByte(o)
             frame.LocalVars().SetInt(i+1, v)
         case 'C':
             o := args[i]
             if o==nil || o.Class().Name() != "java/lang/Character" {
                 return contract.Gas, ErrParameter
             }
             v := heap.UnboxChar(o)
             frame.LocalVars().SetInt(i+1, v)
         case 'S':
             o := args[i]
             if o==nil || o.Class().Name() != "java/lang/Short" {
                 return contract.Gas, ErrParameter
             }
             v := heap.UnboxShort(o)
             frame.LocalVars().SetInt(i+1, v)
         case 'Z':
             o := args[i]
             if o==nil || o.Class().Name() != "java/lang/Boolean" {
                 return contract.Gas, ErrParameter
             }
             v := heap.UnboxBool(o)
             frame.LocalVars().SetInt(i+1, v)
         case 'D':
             o := args[i]
             if o==nil || o.Class().Name() != "java/lang/Double" {
                 return contract.Gas, ErrParameter
             }
             v := heap.UnboxDouble(o)
             frame.LocalVars().SetDouble(i+1, v)
         case 'F':
             o := args[i]
             if o==nil || o.Class().Name() != "java/lang/Float" {
                 return contract.Gas, ErrParameter
             }
             v := heap.UnboxFloat(o)
             frame.LocalVars().SetFloat(i+1, v)
         case 'J':
             o := args[i]
             if o==nil || o.Class().Name() != "java/lang/Long" {
                 return contract.Gas, ErrParameter
             }
             v := heap.UnboxLong(o)
             frame.LocalVars().SetLong(i+1, v)
         case '[', 'L':
             frame.LocalVars().SetRef(i+1, args[i])
         }
     }

     gasLeft, err := interpret(self.mainThread, false, contract, evm)
     if err == nil {
         persistObjectGraph(obj, contractAddr, stateDB)
     } else {
         return gasLeft, err
     }

      returnValueHandler(bogusFrame, method.ReturnTypeDescriptor())
      return gasLeft, err
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
     } else {
         // Important! Otherwise may get a false null
         writeBytes(append(pathPrefix, NullSignal), []byte{}, contractAddr, stateDB)         
     }
     if path, ok := persisted[obj]; ok {
         writeBytes(append(pathPrefix, RefEntry), pathToBytes(path), contractAddr, stateDB)
         return
     } else {
         writeBytes(append(pathPrefix, RefEntry), []byte{}, contractAddr, stateDB)
     }
     persisted[obj] = pathPrefix

     if persistType {
         writeBytes(append(pathPrefix, TypeEntry), []byte(obj.Class().Name()), contractAddr, stateDB)
     }

     if obj.Class().IsArray() {
         write(append(pathPrefix, ArrayLengthEntry), int32ToHash(obj.ArrayLength()), contractAddr, stateDB)
         var i uint
         for i=0; i<uint(obj.ArrayLength()); i++ {
             path := append(pathPrefix, i+SlotIndexOffset) 

             switch obj.Data().(type) {
             case []int8:
                 write(path, int8ToHash(obj.Data().([]int8)[i]), contractAddr, stateDB)
             case []int16:
                 write(path, int16ToHash(obj.Data().([]int16)[i]), contractAddr, stateDB)
             case []int32:
                 write(path, int32ToHash(obj.Data().([]int32)[i]), contractAddr, stateDB)
             case []int64:
                 write(path, int64ToHash(obj.Data().([]int64)[i]), contractAddr, stateDB)
             case []uint16:
                 write(path, uint16ToHash(obj.Data().([]uint16)[i]), contractAddr, stateDB)
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

             //todo optimize for Ljava/lang/String; and Lblockchain/types/Address;

             switch descriptor[0] {
             case 'Z', 'B', 'C', 'S', 'I':
                 write(path, int32ToHash(slots.GetInt(slotId)), contractAddr, stateDB)
             case 'F':
                 write(path, floatToHash(slots.GetFloat(slotId)), contractAddr, stateDB)
             case 'J':
                 write(path, int64ToHash(slots.GetLong(slotId)), contractAddr, stateDB)
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

func uint64ToHash(i uint64) common.Hash {
     bs := make([]byte, 8)
     binary.LittleEndian.PutUint64(bs, i)
     return common.BytesToHash(bs)
}

func int8ToHash(i int8) common.Hash {
     return uint64ToHash(uint64(uint8(i))) //convert to uint8 first to avoid signed expansion
}

func int16ToHash(i int16) common.Hash {
     return uint64ToHash(uint64(uint16(i)))
}

func int32ToHash(i int32) common.Hash {
     return uint64ToHash(uint64(uint32(i)))
}

func int64ToHash(i int64) common.Hash {
     return uint64ToHash(uint64(i))
}

func uint16ToHash(i uint16) common.Hash {
     return uint64ToHash(uint64(i))
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