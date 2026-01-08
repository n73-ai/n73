package utils_test

import (
	"ai-zustack/utils"
	"testing"
)

func TestPurgePullZoneCache(t *testing.T) {
	pullZoneID := int64(5117718) 
	if err := utils.PurgePullZoneCache(pullZoneID); err != nil {
		t.Errorf(err.Error())
	}
}

func TestAddRedirectEdgeRule(t *testing.T) {
	pullZoneID := int64(5118328)               
	bunnyHostname := "kool-pullzone.b-cdn.net" 
	customDomain := "hello.agustfricke.com"    

	if err := utils.AddRedirectEdgeRule(pullZoneID, bunnyHostname, customDomain); err != nil {
		t.Errorf(err.Error())
	}
}

func TestLoadFreeCertificate(t *testing.T) {
	customHostname := "hello.agustfricke.com" 
	if err := utils.LoadFreeCertificate(customHostname); err != nil {
		t.Errorf(err.Error())
	}
}

func TestEnableForceSSL(t *testing.T) {
	pullZoneID := int64(5118328)              
	customHostname := "hello.agustfricke.com" 
	if err := utils.EnableForceSSL(pullZoneID, customHostname); err != nil {
		t.Errorf(err.Error())
	}
}

func TestAddCustomHostname(t *testing.T) {
	pullZoneID := int64(5118328)              
	customHostname := "hello.agustfricke.com" 
	if err := utils.AddCustomHostname(pullZoneID, customHostname); err != nil {
		t.Errorf(err.Error())
	}
}

func TestCreatePullZone(t *testing.T) {
	err := utils.CreatePullZone()
	if err != nil {
		t.Errorf(err.Error())
	}
}

func TestUploadDirectory(t *testing.T) {
	err := utils.UploadDirectory()
	if err != nil {
		t.Errorf(err.Error())
	}
}

func TestCreateStorageZone(t *testing.T) {
  name := "foo"
	err := utils.CreateStorageZone(name)
	if err != nil {
		t.Errorf(err.Error())
	}
}
