package fly_test

import (
	"ai-zustack/fly"
	"testing"
)

func TestCreateApp(t *testing.T) {
	err := fly.CreateApp("n73-test")
	if err != nil {
		t.Errorf(err.Error())
	}
}

func TestCreateMachine(t *testing.T) {
	err := fly.CreateMachine("/home/agust/work/ai/claude/fly.toml")
	if err != nil {
		t.Errorf(err.Error())
	}
}
