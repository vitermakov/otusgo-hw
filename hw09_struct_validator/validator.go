package hw09_struct_validator

import (
	"github.com/vitermakov/otusgo-hw/hw09_struct_validator/validator"
)

/*
	func main() {
		var u = User1{
			Role: "admin",
			Response: struct {
				Code int `validate:"in:200,400,404"`
				Body string
			}{
				Code: 200,
				Body: "",
			},
			Address: Address{
				Zipcode: "555934",
			},
			Email: "dd@yandex.ru",
			Age:   22,
		}
		err := Validate(u)
		fmt.Println(err)
	}
*/
func Validate(v interface{}) error {
	return validator.ValidateStruct(v)
}
