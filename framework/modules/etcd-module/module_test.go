package etcd_module_test

import (
	"testing"

	"github.com/ironzhang/matrix/framework"
	_ "github.com/ironzhang/matrix/framework/modules/etcd-module"
)

func TestModule(t *testing.T) {
	framework.Main()
}
