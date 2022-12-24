package jsontype

type Float32 float32
type Float64 float64

func (value *Float32) UnmarshalJSON(data []byte) (err error) {
	var val float32
	if err = UnmarshalJSON(data, &val); err == nil {
		*value = Float32(val)
	}
	return
}

func (value *Float64) UnmarshalJSON(data []byte) (err error) {
	var val float64
	if err = UnmarshalJSON(data, &val); err == nil {
		*value = Float64(val)
	}
	return
}
