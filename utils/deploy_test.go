package utils_test

import (
	"ai-zustack/utils"
	"fmt"
	"testing"
)

func TestPageExist(t *testing.T) {
	cloudflarePageName := "project-8cf17d28-c45d-40df-9a24-bcb0967bc363"
	exist, err := utils.PageExists(cloudflarePageName)
	if err != nil {
		t.Errorf(err.Error())
	}
	if exist {
		fmt.Println("page exist")
	} else {
		fmt.Println("page do not exist")
	}

}
