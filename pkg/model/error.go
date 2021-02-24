package model

var _ error = (*CustomError)(nil)

type CustomError string

func (c CustomError) Error() string {
	return string(c)
}

const (
	ErrConfirmReject CustomError = CustomError("reject Login Confirm")

	ErrConfirmCancel CustomError = CustomError("cancel Login Confirm")

	ErrConfirmRequestFailure CustomError = CustomError("request Failed")
)
