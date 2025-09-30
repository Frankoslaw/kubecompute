package util

import (
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
)

func GetOsWorkDir() string {
	var base string
	if runtime.GOOS == "windows" {
		base = os.Getenv("ProgramData")
		if base == "" {
			slog.Warn("ProgramData environment variable not set, using temp dir instead")
			base = os.TempDir()
		}
	} else {
		base = "/var/lib"
	}

	return filepath.Join(base, "kubecompute")
}
