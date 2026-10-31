package main

import (
	"bytes"
	"testing"
)

func TestColorWriterColorsOnlyLogLevel(t *testing.T) {
	var output bytes.Buffer
	writer := colorWriter{output: &output}

	message := "time=2026-10-30T23:04:23Z level=ERROR msg=\"request failed\" status=500\n"
	if _, err := writer.Write([]byte(message)); err != nil {
		t.Fatalf("Write() error = %v", err)
	}

	got := output.String()
	want := "\x1b[37mtime=2026-10-30T23:04:23Z level=\x1b[31mERROR\x1b[37m msg=\"request failed\" status=500\n\x1b[0m"
	if got != want {
		t.Fatalf("Write() = %q, want %q", got, want)
	}
}

func TestColorWriterKeepsUnknownLevelsWhite(t *testing.T) {
	var output bytes.Buffer
	writer := colorWriter{output: &output}

	message := "level=NOTICE msg=\"message\"\n"
	if _, err := writer.Write([]byte(message)); err != nil {
		t.Fatalf("Write() error = %v", err)
	}

	want := "\x1b[37m" + message + "\x1b[0m"
	if output.String() != want {
		t.Fatalf("Write() = %q, want %q", output.String(), want)
	}
}
