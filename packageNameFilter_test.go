package loglight_test

import (
  "testing"
  "github.com/myles-mcdonnell/loglight"
)

func TestWhitelist (t *testing.T) {

  filter := loglight.NewPackageNameFilter([]string{"test"}, true)

  if !filter.Filter("test") {
    t.Fail()
  }
}

func TestBlacklist (t *testing.T) {

  filter := loglight.NewPackageNameFilter([]string{"test"}, false)

  if filter.Filter("test") {
    t.Fail()
  }
}