package loglight_test

import (
  "testing"
  "github.com/myles-mcdonnell/loglight"
)

func TestWhitelist{

  filter := loglight.NewPackageNameFilter([]string{"test"}, true)

  if !filter.Filter("test") {
    t.Fail()
  }
}