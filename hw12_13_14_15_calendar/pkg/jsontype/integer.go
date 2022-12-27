package jsontype

type (
	Int   int
	Int32 int32
	Int64 int64
)

type (
	Uint   uint
	Uint32 uint32
	Uint64 uint64
)

func (value *Int) UnmarshalJSON(data []byte) (err error) {
	var val int
	if err = UnmarshalJSON(data, &val); err == nil {
		*value = Int(val)
	}
	return
}

func (value *Int32) UnmarshalJSON(data []byte) (err error) {
	var val int32
	if err = UnmarshalJSON(data, &val); err == nil {
		*value = Int32(val)
	}
	return
}

func (value *Int64) UnmarshalJSON(data []byte) (err error) {
	var val int64
	if err = UnmarshalJSON(data, &val); err == nil {
		*value = Int64(val)
	}
	return
}

func (value *Uint) UnmarshalJSON(data []byte) (err error) {
	var val uint
	if err = UnmarshalJSON(data, &val); err == nil {
		*value = Uint(val)
	}
	return
}

func (value *Uint32) UnmarshalJSON(data []byte) (err error) {
	var val uint32
	if err = UnmarshalJSON(data, &val); err == nil {
		*value = Uint32(val)
	}
	return
}

func (value *Uint64) UnmarshalJSON(data []byte) (err error) {
	var val uint64
	if err = UnmarshalJSON(data, &val); err == nil {
		*value = Uint64(val)
	}
	return
}
