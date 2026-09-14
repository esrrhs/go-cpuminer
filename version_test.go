package main

import (
	"strings"
	"testing"
)

func TestVersionInfo(t *testing.T) {
	if Version == "" {
		t.Fatal("Version should not be empty")
	}

	info := GetVersionInfo()
	if !strings.Contains(info, "version "+Version) {
		t.Fatalf("expected version %s in info, got %s", Version, info)
	}

	BuildTime = "2026-09-14T12:00:00Z"
	GitCommit = "abcdef123456"
	GitBranch = "master"
	defer func() {
		BuildTime = ""
		GitCommit = ""
		GitBranch = ""
	}()

	infoWithBuild := GetVersionInfo()
	if !strings.Contains(infoWithBuild, "2026-09-14T12:00:00Z") {
		t.Fatalf("expected build time in info, got %s", infoWithBuild)
	}
	if !strings.Contains(infoWithBuild, "abcdef123456") {
		t.Fatalf("expected commit in info, got %s", infoWithBuild)
	}
	if !strings.Contains(infoWithBuild, "master") {
		t.Fatalf("expected branch in info, got %s", infoWithBuild)
	}
}
