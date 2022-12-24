package errx

const (
	TypeNone    = iota
	TypeLogic   // логическая ошибка
	TypePerms   // ошибка прав доступа
	TypeInvalid // ошибка валидации
	TypeFatal   // критическая внешняя ошибка
)

type Base struct {
	err  error
	kind byte
}

func (err Base) Error() string {
	return err.err.Error()
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

func LogicNew(err error, code int) Logic {
	return Logic{Base{err, TypeLogic}, code}
}

func PermsNew(err error) Base {
	return Base{err, TypePerms}
}

func FatalNew(err error) Base {
	return Base{err, TypeFatal}
}

type Invalid struct {
	message string
	errors  ValidationErrors
}

func (err Invalid) Error() string {
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
