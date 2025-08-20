package utils_test

import (
	"ai-zustack/utils"
	"testing"
)

func TestDeleteGhRepo(t *testing.T) {
  t.Run("delete remote github repository", func(t *testing.T) {
    projectID := "7504268a-69df-413e-baac-1d13158fd3cb"
    err := utils.DeleteGhRepo(projectID)
    if err != nil {
      t.Errorf("DeleteGhRepo() faild because of: %v", err.Error())
    }
  })
}

func TestDeleteCfPage(t *testing.T) {
  t.Run("delete remote cloudflare page", func(t *testing.T) {
    projectID := "7504268a-69df-413e-baac-1d13158fd3cb"
    err := utils.DeleteCfPage(projectID)
    if err != nil {
      t.Errorf("DeleteGhRepo() faild because of: %v", err.Error())
    }
  })
}
