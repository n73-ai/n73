package utils_test

import (
	"ai-zustack/utils"
	"testing"
)

func TestUnzip(t *testing.T) {
  err := utils.Unzip("/home/agust/some.zip", "/home/agust/personal")
  if err != nil {
    t.Errorf(err.Error())
  }
}
