package jsoncfg

import (
	"os"
	"testing"
	"time"
)

func TestDuration(t *testing.T) {
	type Config struct {
		A Duration
		B Duration
	}

	want := Config{
		A: Duration(60 * time.Second),
		B: Duration(90 * time.Minute),
	}

	if err := WriteToFile("duration.json", &want); err != nil {
		t.Fatalf("write to file: %v", err)
	}

	var got Config
	if err := LoadFromFile("duration.json", &got); err != nil {
		t.Fatalf("load from file: %v", err)
	}

	if got != want {
		t.Errorf("%v != %v", got, want)
	} else {
		t.Logf("%v == %v", got, want)
	}

	os.Remove("duration.json")
}
