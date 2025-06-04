package utils

import "testing"

func TestGenerateAndValidateJWT(t *testing.T) {
	token, err := GenerateJWT(1, "test@example.com", "admin")
	if err != nil {
		t.Fatal(err)
	}
	valid, err := IsValidJWT(token)
	if err != nil || !valid {
		t.Error("token should be valid")
	}
}
