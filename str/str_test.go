package str

import (
	"testing"
)

func TestFirstWord(t *testing.T) {
	if FirstWord("  SELECT ABC FROM DUAL") != "SELECT" {
		t.Error("FIRST WORD of [SELECT ABC FROM DUAL] ERROR")
	}
}

func TestStringContains(t *testing.T) {
	if !StringContains("all", "some", " ", "all") {
		t.Error("failed")
	}

	if StringContains("all", "some", " ", "") {
		t.Error("failed")
	}

	if !StringContains("a b", "a", " ", "") {
		t.Error("failed")
	}

	if StringContains("b", "a", " ", "") {
		t.Error("failed")
	}

	if StringContains("aa b", "a", " ", "") {
		t.Error("failed")
	}
}
