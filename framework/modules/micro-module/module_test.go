package micro_module_test

import (
	"testing"

	"github.com/ironzhang/matrix/framework"

	_ "github.com/ironzhang/matrix/framework/modules/micro-module"
)

func TestModule(t *testing.T) {
	framework.Main()
}
