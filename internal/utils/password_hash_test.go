package utils

import "testing"

func TestHashAndCheckPassword(t *testing.T) {
	pw := "testpassword"
	hash, err := HashPassword(pw)
	if err != nil {
		t.Fatal(err)
	}
	if !CheckPassword(hash, pw) {
		t.Error("password should match hash")
	}
	if CheckPassword(hash, "wrong") {
		t.Error("wrong password should not match hash")
	}
}
