package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"golang.org/x/crypto/bcrypt"
)

const testCost = bcrypt.MinCost

func mustRandom(t *testing.T, n int) string {
	t.Helper()
	s, err := randomString(n)
	if err != nil {
		t.Fatalf("randomString: %v", err)
	}
	return s
}

func Test_generateHashForPassword(t *testing.T) {
	lines, err := generateHashForPassword("superSecretPassword123", testCost)
	if err != nil || len(lines) != 1 {
		t.Fatalf("unexpected result: %v, %v", lines, err)
	}
	if err := checkHashAndPassword(lines[0]); err != nil {
		t.Error("hash and password do not match")
	}
}

func Test_generateHashForRandomPassword(t *testing.T) {
	lines, err := generateHashForPassword(mustRandom(t, 20), testCost)
	if err != nil {
		t.Fatal(err)
	}
	if err := checkHashAndPassword(lines[0]); err != nil {
		t.Error("hash and password do not match")
	}
}

func Test_passwordWithSpacesRoundTrips(t *testing.T) {
	lines, err := generateHashForPassword("correct horse battery", testCost)
	if err != nil {
		t.Fatal(err)
	}
	if err := checkHashAndPassword(lines[0]); err != nil {
		t.Errorf("password containing spaces failed to validate: %v", err)
	}
}

func Test_wrongPasswordFails(t *testing.T) {
	lines, _ := generateHashForPassword("right", testCost)
	wrong := strings.Replace(lines[0], "right", "wrong", 1)
	if err := checkHashAndPassword(wrong); err == nil {
		t.Error("expected mismatch for wrong password")
	}
	if _, err := matchPasswordAndHash(wrong); err == nil {
		t.Error("matchPasswordAndHash should report failure")
	}
}

func Test_malformedLineDoesNotPanic(t *testing.T) {
	for _, line := range []string{"", "nospace", "   "} {
		if _, err := matchPasswordAndHash(line); err == nil {
			t.Errorf("expected error for %q", line)
		}
	}
}

func Test_RandomPasswordLength(t *testing.T) {
	for _, n := range []int{0, 1, 16, 64} {
		if got := len(mustRandom(t, n)); got != n {
			t.Errorf("length %d: got %d", n, got)
		}
	}
}

func Test_RandomPasswordCharset(t *testing.T) {
	for _, c := range mustRandom(t, 500) {
		if !strings.ContainsRune(randomChars, c) {
			t.Fatalf("unexpected character %q", c)
		}
	}
}

func Test_writeLinesPermissions(t *testing.T) {
	path := filepath.Join(t.TempDir(), "out.txt")
	if err := writeLines([]string{"a b"}, path); err != nil {
		t.Fatal(err)
	}
	info, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	if perm := info.Mode().Perm(); perm != 0o600 {
		t.Errorf("expected mode 0600, got %o", perm)
	}
}

func Test_readLinesSkipsBlank(t *testing.T) {
	path := filepath.Join(t.TempDir(), "in.txt")
	if err := os.WriteFile(path, []byte("one\n\ntwo\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	lines, err := readLines(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(lines) != 2 || lines[0] != "one" || lines[1] != "two" {
		t.Errorf("unexpected lines: %q", lines)
	}
}

func Test_hashPasswordsFromFileSkipsLong(t *testing.T) {
	long := strings.Repeat("x", 73)
	lines, skipped, err := hashPasswordsFromFile([]string{"short", long, "other"}, testCost)
	if err != nil {
		t.Fatal(err)
	}
	if skipped != 1 || len(lines) != 2 {
		t.Errorf("skipped=%d lines=%d, want 1 and 2", skipped, len(lines))
	}
	if err := checkHashAndPassword(lines[0]); err != nil {
		t.Error("first line should validate")
	}
}

func Test_generateLongPasswordErrors(t *testing.T) {
	if _, err := generateHashForPassword(strings.Repeat("x", 73), testCost); err == nil {
		t.Error("expected error for password over 72 bytes")
	}
}
