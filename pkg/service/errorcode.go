package service

const (
	ErrLoginConfirmWait     = "login_confirm_wait"
	ErrLoginConfirmRejected = "login_confirm_rejected"
	ErrLoginConfirmRequired = "login_confirm_required"
	ErrMFARequired          = "mfa_required"
	ErrPasswordFailed       = "password_failed"
)

const (
	ErrSessionLoginConfirmWait =  "session_login_confirm_wait"
	ErrSessionLoginConfirmRejected = "session_login_confirm_rejected"
)

const successMsg = "ok"
