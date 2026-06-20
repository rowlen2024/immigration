package service

import (
	"testing"
	"time"

	"mygo-immigration/backend/internal/repository"
)

func TestPublicVersionRegistryResolve(t *testing.T) {
	reg := NewPublicVersionRegistry()
	base := time.Date(2026, 6, 20, 12, 0, 0, 0, time.UTC)

	reg.Register("public:project:", func(key string) (repository.PublicVersion, error) {
		return repository.PublicVersion{UpdatedAt: base, Count: int64(len(key))}, nil
	})
	reg.Register("public:project:special", func(string) (repository.PublicVersion, error) {
		return repository.PublicVersion{UpdatedAt: base.Add(time.Hour), Count: 1}, nil
	})

	got, err := reg.Resolve([]string{
		"public:project:eb5",
		"public:project:special",
		"public:unknown",
	})
	if err != nil {
		t.Fatalf("Resolve returned error: %v", err)
	}

	if _, ok := got["public:unknown"]; ok {
		t.Fatal("unknown key should not be returned")
	}
	if got["public:project:eb5"] != "2026-06-20T12:00:00Z:18" {
		t.Fatalf("unexpected project version: %s", got["public:project:eb5"])
	}
	if got["public:project:special"] != "2026-06-20T13:00:00Z:1" {
		t.Fatalf("longest prefix resolver was not used: %s", got["public:project:special"])
	}
}
