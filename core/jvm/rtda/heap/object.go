package heap

type Object struct {
	class *Class
	data  interface{} // Slots for Object, []int32 for int[] ...
	extra interface{}
}

// create normal (non-array) object
func newObject(class *Class) *Object {
	return &Object{
		class: class,
		data:  newSlots(class.instanceSlotCount),
	}
}

// getters & setters
func (self *Object) Class() *Class {
	return self.class
}
func (self *Object) Data() interface{} {
	return self.data
}
func (self *Object) Fields() Slots {
	return self.data.(Slots)
}
func (self *Object) Extra() interface{} {
	return self.extra
}
func (self *Object) SetExtra(extra interface{}) {
	self.extra = extra
}

func (self *Object) IsInstanceOf(class *Class) bool {
	return class.IsAssignableFrom(self.class)
}

// reflection
func (self *Object) GetRefVar(name, descriptor string) *Object {
	field := self.class.getField(name, descriptor, false)
	slots := self.data.(Slots)
	return slots.GetRef(field.slotId)
}
func (self *Object) SetRefVar(name, descriptor string, ref *Object) {
	field := self.class.getField(name, descriptor, false)
	slots := self.data.(Slots)
	slots.SetRef(field.slotId, ref)
}
func (self *Object) SetIntVar(name, descriptor string, val int32) {
	field := self.class.getField(name, descriptor, false)
	slots := self.data.(Slots)
	slots.SetInt(field.slotId, val)
}
func (self *Object) GetIntVar(name, descriptor string) int32 {
	field := self.class.getField(name, descriptor, false)
	slots := self.data.(Slots)
	return slots.GetInt(field.slotId)
}
func (self *Object) GetFloatVar(name, descriptor string) float32 {
	field := self.class.getField(name, descriptor, false)
	slots := self.data.(Slots)
	return slots.GetFloat(field.slotId)
}
func (self *Object) SetFloatVar(name, descriptor string, val float32) {
	field := self.class.getField(name, descriptor, false)
	slots := self.data.(Slots)
	slots.SetFloat(field.slotId, val)
}
func (self *Object) GetDoubleVar(name, descriptor string) float64 {
	field := self.class.getField(name, descriptor, false)
	slots := self.data.(Slots)
	return slots.GetDouble(field.slotId)
}
func (self *Object) SetDoubleVar(name, descriptor string, val float64) {
	field := self.class.getField(name, descriptor, false)
	slots := self.data.(Slots)
	slots.SetDouble(field.slotId, val)
}
func (self *Object) GetLongVar(name, descriptor string) int64 {
	field := self.class.getField(name, descriptor, false)
	slots := self.data.(Slots)
	return slots.GetLong(field.slotId)
}
func (self *Object) SetLongVar(name, descriptor string, val int64) {
	field := self.class.getField(name, descriptor, false)
	slots := self.data.(Slots)
	slots.SetLong(field.slotId, val)
}
