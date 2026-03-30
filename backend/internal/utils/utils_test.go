package utils

import "testing"

func TestGetEnvVar(t *testing.T) {
	t.Run("use env value", func(t *testing.T) {
		t.Setenv("GET_ENV_VAR_TEST", " value ")
		got := GetEnvVar("GET_ENV_VAR_TEST", "default")
		if got != "value" {
			t.Fatalf("expected value, got %q", got)
		}
	})

	t.Run("fallback to default", func(t *testing.T) {
		t.Setenv("GET_ENV_VAR_TEST", "")
		got := GetEnvVar("GET_ENV_VAR_TEST", " default ")
		if got != "default" {
			t.Fatalf("expected default, got %q", got)
		}
	})
}
