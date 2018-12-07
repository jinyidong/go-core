package autofac

import "reflect"

var typeRegistry = make(map[string]reflect.Type)

//注意使用方法，是否可以封装的更好
//RegisterType((*StructName)(nil))
func RegisterType(elem interface{}) {
	t := reflect.TypeOf(elem).Elem()
	typeRegistry[t.Name()] = t
}

//使用时有一个问题需要注意，通过NewStruct()之后，需要通过类型断言***.(struct)，才可完成string的反序列化工作
func NewStructByType(name string) (interface{}, bool) {
	elem, ok := typeRegistry[name]

	if !ok {
		return nil, false
	}
	return reflect.New(elem).Elem().Interface(), true
}
