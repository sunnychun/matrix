package server

import (
	"errors"
	"testing"
)

func errorEqual(e1, e2 error) bool {
	if e1 != nil && e2 != nil {
		return e1.Error() == e2.Error()
	}
	return e1 == e2
}

func TestParsePath(t *testing.T) {
	type test struct {
		path    string
		wantOP  string
		wantKey string
		wantErr error
	}

	tests := []test{
		{"/", "", "/", nil},
		{"//", "", "/", nil},
		{"///", "", "//", nil},
		{"/list", "list", "/", nil},
		{"/list/", "list", "/", nil},
		{"/list/dev", "list", "/dev/", nil},
		{"/list/dev/", "list", "/dev/", nil},

		{path: "", wantErr: errors.New(`invalid path: ""`)},
		{path: "list", wantErr: errors.New(`invalid path: "list"`)},
		{path: "list/", wantErr: errors.New(`invalid path: "list/"`)},
	}

	for _, tt := range tests {
		gotOP, gotKey, gotErr := parsePath(tt.path)
		if gotOP != tt.wantOP || gotKey != tt.wantKey || !errorEqual(gotErr, tt.wantErr) {
			t.Errorf("path=%s; got={%q, %q, %v}; want={%q, %q, %v}", tt.path, gotOP, gotKey, gotErr, tt.wantOP, tt.wantKey, tt.wantErr)
		} else {
			t.Logf("path=%s; got={%q, %q, %v}; want={%q, %q, %v}", tt.path, gotOP, gotKey, gotErr, tt.wantOP, tt.wantKey, tt.wantErr)
		}
	}
}
