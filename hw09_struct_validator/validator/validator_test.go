package validator_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/vitermakov/otusgo-hw/hw09_struct_validator/validator"
)

type UserRole string
type Status string

// Test the function on different structures and other types.
type (
	Address struct {
		Zipcode  string `validate:"len:6"`
		Value    string
		GeoPoint Point `validate:"nested"`
	}
	Point struct {
		Longtitude float32 `validate:"min:-180|max:180"`
		Latitude   float32 `validate:"min:-180|max:180"` //`validate:"len:15"`
	}
	Object struct {
		ID   string `json:"id" validate:"len:36"`
		Name string
	}
	User struct {
		// встроенный
		Object `validate:"nested"`
		Age    int      `validate:"min:18|max:50"`
		Email  string   `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole `validate:"in:admin,stuff"`
		Phones []string `validate:"len:11"`
		Status struct {
			Code Status `validate:"in:default,approved,banned"` // не должно проверяться, так как Status не `nested`
			Text string
		}
		Address   Address `validate:"nested"`
		isBlocked bool
	}
)

func (u *User) Block(value bool) {
	u.isBlocked = value
}
func (u User) IsBlocked() bool {
	return u.isBlocked
}

func TestValidateInit(t *testing.T) {
	err := validator.ValidateStruct(nil)
	require.True(t, errors.Is(err, validator.ErrInputStructIsNull))

	err = validator.ValidateStruct(32)
	require.True(t, errors.Is(err, validator.ErrInputNotStruct))

	// тестируем, что ошибки парсинга тегов прокидываются из ValidateStruct
	errorRetrieveCases := []interface{}{
		struct {
			Value int `validate:"min:x"`
		}{},
		struct {
			Value int `validate:"min:6,11"`
		}{},
		struct {
			Value string `validate:"min:2"`
		}{},
		struct {
			Value int `validate:"max:x"`
		}{},
		struct {
			Value int `validate:"max:6,11"`
		}{},
		struct {
			Value map[string]int `validate:"max:2"`
		}{},
		struct {
			Value string `validate:"len:x"`
		}{},
		struct {
			Value string `validate:"len:6,11"`
		}{},
		struct {
			Value float32 `validate:"len:6"`
		}{},
		struct {
			Value int `validate:"in:x,y"`
		}{},
		struct {
			Value bool `validate:"in:x,y"`
		}{},
		struct {
			Value string `validate:"regexp:(("`
		}{},
		struct {
			Value int `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		}{},
		struct {
			Value int `validate:"nested"`
		}{},
	}

	for i, ec := range errorRetrieveCases {
		t.Run(fmt.Sprintf("error case retrive %d", i+1), func(t *testing.T) {
			ec := ec
			t.Parallel()

			err := validator.ValidateStruct(ec)
			require.Error(t, err)
		})
	}
}

func TestValidateCheck(t *testing.T) {
	invalid := User{
		Object: Object{
			ID:   "444",
			Name: "Yury",
		},
		Role: "admin1",
		Status: struct {
			Code Status `validate:"in:default,approved,banned"`
			Text string
		}{
			Code: "approved1",
			Text: "",
		},
		Phones: []string{
			"790011122331",
			"7900111223",
		},
		Address: Address{
			Zipcode: "5559341",
			GeoPoint: Point{
				Longtitude: 200,
				Latitude:   -200,
			},
		},
		Email: "@yandex.ru",
		Age:   12,
	}
	checkErrorSet(
		t,
		validator.ValidateStruct(invalid),
		[]string{
			"User.Object.ID",
			"User.Role",
			"User.Phones.0",
			"User.Phones.1",
			"User.Address.Zipcode",
			"User.Address.GeoPoint.Longtitude",
			"User.Address.GeoPoint.Latitude",
			"User.Email",
			"User.Age",
		},
	)

	valid := User{
		Object: Object{
			ID:   "9a7ef00e-5991-11ec-a009-d00dde7fb0c3",
			Name: "Alex",
		},
		Role: "admin",
		Status: struct {
			Code Status `validate:"in:default,approved,banned"`
			Text string
		}{
			Code: "approved",
			Text: "",
		},
		Phones: []string{
			"79001112233",
			"79001112231",
		},
		Address: Address{
			Zipcode: "555934",
		},
		Email: "dd@yandex.ru",
		Age:   22,
	}

	err := validator.ValidateStruct(valid)

	require.NoError(t, err)
}

func checkErrorSet(t *testing.T, err error, errFields []string) {
	if errFields == nil {
		require.NoError(t, err, "no errors expected")
	} else {
		errors, ok := err.(validator.ValidationErrors)
		require.True(t, ok, "err is not ValidationErrors")

		result := make([]string, 0)
		for _, err := range errors {
			result = append(result, err.Field)
		}
		require.ElementsMatch(t, result, errFields)
	}
}