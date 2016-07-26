package errs

var errors = map[Errno]string{}

var (
	ErrUnknown      = Errno(1)
	ErrInvalidParam = Errno(2)
)

func init() {
	errors[0] = "nil"
	errors[ErrUnknown] = "unknown"
	errors[ErrInvalidParam] = "invalid parameter"
}
