package errx

const (
	TypeLogic   byte = 1 // логическая ошибка
	TypePerms   byte = 2 // ошибка прав доступа
	TypeInvalid byte = 3 // ошибка валидации
	TypeFatal   byte = 4 // критическая внешняя ошибка
)

type Base struct {
	message string
	kind    byte
}

func (err Base) Error() string {
	return err.message
}

func (err Base) Kind() byte {
	return err.kind
}

// Logic ошибка бизнес-логики (устранимая) приложения.
type Logic struct {
	Base
	code int // код внутренней классификации ошибок.
}

func (err Logic) Code() int {
	return err.code
}

func LogicNew(message string, code int) Logic {
	return Logic{Base{message, TypeLogic}, code}
}

func PermsNew(message string) Base {
	return Base{message, TypePerms}
}

func FatalNew(message string) Base {
	return Base{message, TypeFatal}
}

type Invalid struct {
	message string
	errors  ValidationErrors
}

func (err Invalid) Error() string {
	/*
		var result string
		messages := err.Errors()
		for _, message := range messages {
			if len(result) == 0 {
				result = message.String()
			} else {
				result += "; " + message.String()
			}
		}
	*/
	return err.message
}

func (err Invalid) Errors() []ValidationError {
	return err.errors
}

func (err Invalid) Fails() bool {
	return !err.errors.Empty()
}

func (err Invalid) Kind() byte {
	return TypeInvalid
}

func InvalidNew(message string) Invalid {
	return Invalid{message: message}
}
