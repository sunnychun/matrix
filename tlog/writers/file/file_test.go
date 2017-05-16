package file

import (
	"fmt"
	"os"
	"testing"
)

func TestFile(t *testing.T) {
	f, err := Open(Options{
		Dir:      os.TempDir(),
		Program:  "TestFile",
		MaxBytes: 500,
	})
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	for i := 0; i < 10; i++ {
		fmt.Fprintf(f, "%d: Hello, world\n", i)
	}
}
