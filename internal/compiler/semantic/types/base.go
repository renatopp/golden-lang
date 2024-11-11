package types

var _type_id = uint64(0)

type baseType struct{ id uint64 }

func newBase() baseType {
	_type_id++
	return baseType{_type_id}
}
func (t baseType) Id() uint64 { return t.id }
