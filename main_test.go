package main

import (
	"os"
	"strings"
	"testing"
)

func assertErr(t *testing.T, got string, want string) {
	t.Helper()
	if got == "" {
		t.Fatalf("didn't get an error but expected %q", want)
	}
	if got != want {
		t.Fatalf("got the wrong error, got %q want %q", got, want)
	}
}

func TestProgram(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	t.Run("Unknown", func(t *testing.T) {
		defer os.Chdir(wd)
		_, err := run([]string{"envdir", "testdata/env", "thisprogramdoesnotexist"})
		assertErr(t, err, `exec: "thisprogramdoesnotexist": executable file not found in $PATH`)
	})

	t.Run("ExtraNewlinesStripped", func(t *testing.T) {
		defer os.Chdir(wd)
		execArgs, err := run([]string{"envdir", "testdata/env", "env"})
		if err != "" {
			t.Fatalf("expected no error message, got %q", err)
		}
		val, ok := getenv(execArgs.Env, "EXTRA_NEWLINES")
		if !ok {
			t.Fatal("did not find EXTRA_NEWLINES var in output")
		}
		want := "Varwithtabsandnewlines\n\nat the end"
		if val != want {
			t.Errorf("EXTRA_NEWLINES: want %q, got %q", want, val)
		}
	})

	t.Run("EmptyVarDeleted", func(t *testing.T) {
		defer os.Chdir(wd)
		os.Setenv("SHOULD_BE_DELETED", "not_deleted")
		defer os.Unsetenv("SHOULD_BE_DELETED")
		execArgs, err := run([]string{"envdir", "testdata/env", "env"})
		if err != "" {
			t.Fatalf("expected no error message, got %q", err)
		}
		_, ok := getenv(execArgs.Env, "SHOULD_BE_DELETED")
		if ok {
			t.Fatal("found SHOULD_BE_DELETED var in output, but should have been deleted")
		}
	})

	t.Run("TooLongTruncated", func(t *testing.T) {
		defer os.Chdir(wd)
		execArgs, err := run([]string{"envdir", "testdata/env", "env"})
		if err != "" {
			t.Fatalf("expected no error message, got %q", err)
		}
		val, ok := getenv(execArgs.Env, "1040_BYTE_ENV_VAR")
		if !ok {
			t.Fatal("did not find 1040_BYTE_ENV_VAR var in output")
		}
		want := strings.Repeat("a", 1024)
		if val != want {
			t.Errorf("1040_BYTE_ENV_VAR: want %q, got %q", want, val)
		}
	})
}
