package users

type Status int

type ConstUserError UserError

//UserError - domain errors for user logic validations
type UserError struct {
	code       Status
	message    string
	innerError error
}

const (
	Unknow Status = iota
	NotFound
	AlreadyExistingItem
	InvalidData
)

const (
	USERNOTFOUND      string = "user not found"
	USERALREADYEXISTS string = "user already exists"
	INVALIDDATA       string = "invalid data"
)

func (e UserError) Error() string {
	return e.message
}

func UnknowError(e error) UserError {
	return UserError{code: Unknow, innerError: e, message: e.Error()}
}

func IsUserErrorType(arg1 ConstUserError, arg2 error) bool {
	if other, ok := arg2.(UserError); ok {
		return other.code == arg1.code && other.message == arg1.message
	}
	return false
}

var (
	ERRNOTFOUND      = ConstUserError{code: NotFound, message: USERNOTFOUND}
	ERRALREADYEXISTS = ConstUserError{code: AlreadyExistingItem, message: USERALREADYEXISTS}
	ERRINVALIDDATA   = ConstUserError{code: InvalidData, message: INVALIDDATA}
)
