package utils_test

import (
	"ai-zustack/utils"
	"testing"
)

func TestTryBuildProject(t *testing.T) {
  err := utils.TryBuildProject("/home/agust/work/ai/projects/147fdc5a-381c-43cd-9bd7-579b3660cc25/project")
  if err != nil {
    t.Errorf(err.Error())
  }
}
