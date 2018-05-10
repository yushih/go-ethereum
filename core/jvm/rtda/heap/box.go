package heap

//These methods depends on the fact that in each of the classes there is a 
//private field "value" that contains value of the primitive type. Should 
//call the public methods instead.

func BoxBool(v int32, classLoader *ClassLoader) *Object {
     o := classLoader.LoadClass("java/lang/Boolean").NewObject()
     o.SetIntVar("value", "Z", v)
     return o
}

func BoxByte(v int32, classLoader *ClassLoader) *Object {
     o := classLoader.LoadClass("java/lang/Byte").NewObject()
     o.SetIntVar("value", "B", v)
     return o
}

func BoxChar(v int32, classLoader *ClassLoader) *Object {
     o := classLoader.LoadClass("java/lang/Character").NewObject()
     o.SetIntVar("value", "C", v)
     return o
}

func BoxShort(v int32, classLoader *ClassLoader) *Object {
     o := classLoader.LoadClass("java/lang/Short").NewObject()
     o.SetIntVar("value", "S", v)
     return o
}

func BoxInt(v int32, classLoader *ClassLoader) *Object {
     o := classLoader.LoadClass("java/lang/Integer").NewObject()
     o.SetIntVar("value", "I", v)
     return o
}

func BoxFloat(v float32, classLoader *ClassLoader) *Object {
     o := classLoader.LoadClass("java/lang/Float").NewObject()
     o.SetFloatVar("value", "F", v)
     return o
}

func BoxDouble(v float64, classLoader *ClassLoader) *Object {
     o := classLoader.LoadClass("java/lang/Double").NewObject()
     o.SetDoubleVar("value", "D", v)
     return o
}

func BoxLong(v int64, classLoader *ClassLoader) *Object {
     o := classLoader.LoadClass("java/lang/Long").NewObject()
     o.SetLongVar("value", "J", v)
     return o
}

func UnboxInt(o *Object) int32 {
    return o.GetIntVar("value", "I")
}

func UnboxByte(o *Object) int32 {
    return o.GetIntVar("value", "Z")
}

func UnboxChar(o *Object) int32 {
    return o.GetIntVar("value", "C")
}

func UnboxShort(o *Object) int32 {
    return o.GetIntVar("value", "S")
}

func UnboxBool(o *Object) int32 {
    return o.GetIntVar("value", "Z")
}

func UnboxFloat(o *Object) float32 {
    return o.GetFloatVar("value", "F")
}

func UnboxDouble(o *Object) float64 {
    return o.GetDoubleVar("value", "D")
}

func UnboxLong(o *Object) int64 {
    return o.GetLongVar("value", "J")
}

