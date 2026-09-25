package auth

import (
	"testing"
)

func TestGenerateAndParseToken(
	t *testing.T,
) {
	const secret = "unit-test-secret"

	token, err :=
		GenerateToken(
			123,
			"royandi@example.com",
			secret,
		)

	if err != nil {
		t.Fatalf(
			"GenerateToken gagal: %v",
			err,
		)
	}

	if token == "" {
		t.Fatal(
			"token tidak boleh kosong",
		)
	}

	claims, err :=
		ParseToken(
			token,
			secret,
		)

	if err != nil {
		t.Fatalf(
			"ParseToken gagal: %v",
			err,
		)
	}

	if claims.UserID !=
		123 {
		t.Fatalf(
			"UserID = %d, want 123",
			claims.UserID,
		)
	}

	if claims.Email !=
		"royandi@example.com" {
		t.Fatalf(
			"Email = %q, want royandi@example.com",
			claims.Email,
		)
	}
}

func TestParseToken_WrongSecret(
	t *testing.T,
) {
	token, err :=
		GenerateToken(
			1,
			"user@example.com",
			"correct-secret",
		)

	if err != nil {
		t.Fatalf(
			"GenerateToken gagal: %v",
			err,
		)
	}

	_, err =
		ParseToken(
			token,
			"wrong-secret",
		)

	if err == nil {
		t.Fatal(
			"token dengan secret berbeda seharusnya ditolak",
		)
	}
}
