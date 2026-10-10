// Package profiles ships the static PrusaSlicer printer and print profiles for
// the Creality K1. They are embedded in the binary and written to disk on first
// use, because PrusaSlicer only accepts config files by path.
package profiles

import (
	"embed"
	"os"
	"path/filepath"
	"sync"

	"github.com/gitTanzj/server/config"
)

//go:embed k1/*.ini
var files embed.FS

const (
	printerFile = "k1/printer.ini"
	printFile   = "k1/print.ini"
)

var (
	once        sync.Once
	printerPath string
	printPath   string
	initErr     error
)

// Paths returns the printer and print profile paths, in that load order.
// PRUSA_PRINTER_PROFILE / PRUSA_PRINT_PROFILE override the embedded K1 files.
func Paths() (printer, print string, err error) {
	once.Do(func() {
		printerPath = config.Envs.PrusaPrinterProfile
		printPath = config.Envs.PrusaPrintProfile
		if printerPath != "" && printPath != "" {
			return
		}

		dir, err := os.MkdirTemp("", "kalleprint-profiles-*")
		if err != nil {
			initErr = err
			return
		}
		if printerPath == "" {
			if printerPath, initErr = extract(dir, printerFile); initErr != nil {
				return
			}
		}
		if printPath == "" {
			printPath, initErr = extract(dir, printFile)
		}
	})
	return printerPath, printPath, initErr
}

func extract(dir, name string) (string, error) {
	data, err := files.ReadFile(name)
	if err != nil {
		return "", err
	}
	path := filepath.Join(dir, filepath.Base(name))
	return path, os.WriteFile(path, data, 0o644)
}
