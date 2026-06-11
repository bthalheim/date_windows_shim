package main

import (
	"bytes"
	"testing"
	"time"
)

func TestFormatDateSupportedSpecifiers(t *testing.T) {
	now := time.Date(2026, time.June, 11, 14, 18, 59, 0, time.FixedZone("PDT", -7*60*60))

	got, err := formatDate(now, "%F %T %Z %z %u %w %%")
	if err != nil {
		t.Fatalf("formatDate returned error: %v", err)
	}

	want := "2026-06-11 14:18:59 PDT -0700 4 4 %"
	if got != want {
		t.Fatalf("formatDate = %q, want %q", got, want)
	}
}

func TestRunDefaultOutput(t *testing.T) {
	now := time.Date(2026, time.June, 11, 14, 18, 59, 0, time.FixedZone("PDT", -7*60*60))
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	if err := run(&stdout, &stderr, nil, func() time.Time { return now }); err != nil {
		t.Fatalf("run returned error: %v", err)
	}

	want := "Thu Jun 11 14:18:59 PDT 2026\n"
	if stdout.String() != want {
		t.Fatalf("stdout = %q, want %q", stdout.String(), want)
	}

	if stderr.Len() != 0 {
		t.Fatalf("stderr = %q, want empty", stderr.String())
	}
}

func TestRunUTCFormat(t *testing.T) {
	now := time.Date(2026, time.June, 11, 14, 18, 59, 0, time.FixedZone("PDT", -7*60*60))
	var stdout bytes.Buffer

	if err := run(&stdout, &bytes.Buffer{}, []string{"-u", "+%F %T %Z"}, func() time.Time { return now }); err != nil {
		t.Fatalf("run returned error: %v", err)
	}

	want := "2026-06-11 21:18:59 UTC\n"
	if stdout.String() != want {
		t.Fatalf("stdout = %q, want %q", stdout.String(), want)
	}
}

func TestRunRejectsUnsupportedArgument(t *testing.T) {
	err := run(&bytes.Buffer{}, &bytes.Buffer{}, []string{"--set"}, time.Now)
	if err == nil {
		t.Fatal("expected error for unsupported argument")
	}
}

func TestFormatDateRejectsUnsupportedSpecifier(t *testing.T) {
	_, err := formatDate(time.Now(), "%Q")
	if err == nil {
		t.Fatal("expected error for unsupported format specifier")
	}
}
