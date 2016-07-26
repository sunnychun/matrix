package errs_test

import (
	"testing"

	"github.com/ironzhang/matrix/errs"
)

func TestErrno(t *testing.T) {
	var e error = errs.ErrInvalidParam
	t.Logf("%v\n", e)

	if e != errs.ErrInvalidParam {
		t.Fatalf("%v != %v", e, errs.ErrInvalidParam)
	} else {
		t.Logf("%v == %v", e, errs.ErrInvalidParam)
	}

	err := errs.GetErrno(e)
	t.Logf("%T %d %v", err, err, err)

	e = errs.GetErrno(nil)
	t.Logf("%T %d %v", e, e, e)
}
