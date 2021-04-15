package jsobj

//Object JavaScript Object
type Object interface {
	Interface() interface{}
}

type object struct {
	value interface{}
}

//Interface Get inner interface data from object instance
func (obj *object) Interface() interface{} {
	return obj.value
}
