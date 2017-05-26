package framework_test

import (
	"testing"

	"github.com/ironzhang/matrix/framework"
	"github.com/ironzhang/matrix/framework/modules/pprof-module"
)

var _ = pprof_module.Module

func TestFramework(t *testing.T) {
	framework.Main()
}
