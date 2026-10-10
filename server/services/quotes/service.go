package quotes

import (
	"archive/zip"
	"bufio"
	"fmt"
	"io"
	"math"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/gitTanzj/server/config"
	"github.com/gitTanzj/server/profiles"
	"github.com/gitTanzj/server/services/filamentconfig"
)

const PROFIT_MULTIPLIER = 1.8

// Bed centre of the Creality K1 (220×220).
const bedCenter = "110,110"

func CalculateQuote(amount float64, filamentCost float64) float64 {
	price := amount * filamentCost * PROFIT_MULTIPLIER
	price_rounded := math.Round(price*100) / 100
	return math.Max(price_rounded, 1)
}

// SliceFilament slices the input file and returns grams of filament used.
// If gcodeOut is non-empty the gcode is written there and kept; otherwise a
// temporary file is used and removed automatically.
//
// Config is layered printer → print → filament; later --load files win.
func SliceFilament(inputPath, ext, iniFilePath, gcodeOut string) (float64, error) {
	cleanup := gcodeOut == ""
	if cleanup {
		tmpGcode, err := os.CreateTemp("", "quote-*.gcode")
		if err != nil {
			return 0, err
		}
		tmpGcode.Close()
		gcodeOut = tmpGcode.Name()
		defer os.Remove(gcodeOut)
	}

	printerProfile, printProfile, err := profiles.Paths()
	if err != nil {
		return 0, fmt.Errorf("load printer profiles: %w", err)
	}

	filamentProfile, err := filamentconfig.WriteSlicerFile(iniFilePath)
	if err != nil {
		return 0, err
	}
	defer os.Remove(filamentProfile)

	if strings.EqualFold(ext, ".3mf") {
		stripped, err := strip3MFConfig(inputPath)
		if err != nil {
			return 0, fmt.Errorf("read 3mf: %w", err)
		}
		defer os.Remove(stripped)
		inputPath = stripped
	}

	args := buildSlicerArgs(inputPath, gcodeOut, printerProfile, printProfile, filamentProfile)
	if out, err := exec.Command(config.Envs.PrusaSlicerPath, args...).CombinedOutput(); err != nil {
		return 0, fmt.Errorf("slicer failed: %w: %s", err, out)
	}

	return parseFilamentGrams(gcodeOut)
}

func buildSlicerArgs(inputPath, gcodeOut string, configs ...string) []string {
	args := []string{"--export-gcode", "--output", gcodeOut}
	for _, c := range configs {
		args = append(args, "--load", c)
	}
	return append(args, "--center", bedCenter, inputPath)
}

// slicerConfigEntries are the 3MF entries where PrusaSlicer, OrcaSlicer and
// Bambu Studio store print settings and per-object overrides. Removing them
// leaves plain geometry, so our profiles always decide how a model is sliced.
var slicerConfigEntries = map[string]bool{
	"metadata/slic3r_pe.config":        true,
	"metadata/slic3r_pe_model.config":  true,
	"metadata/prusaslicer.config":      true,
	"metadata/project_settings.config": true,
	"metadata/model_settings.config":   true,
}

// strip3MFConfig copies the 3MF at path to a temp file without embedded slicer
// settings and returns the new path. The caller removes the file.
func strip3MFConfig(path string) (string, error) {
	src, err := zip.OpenReader(path)
	if err != nil {
		return "", err
	}
	defer src.Close()

	out, err := os.CreateTemp("", "stripped-*.3mf")
	if err != nil {
		return "", err
	}
	defer out.Close()

	if err := copyZipWithout(&src.Reader, out, slicerConfigEntries); err != nil {
		os.Remove(out.Name())
		return "", err
	}
	return out.Name(), nil
}

func copyZipWithout(src *zip.Reader, dst io.Writer, skip map[string]bool) error {
	zw := zip.NewWriter(dst)
	for _, f := range src.File {
		if skip[strings.ToLower(f.Name)] {
			continue
		}
		if err := zw.Copy(f); err != nil {
			return err
		}
	}
	return zw.Close()
}

// parseFilamentGrams reads PrusaSlicer's "filament used [g]" summary, which is
// computed from the filament_density in the loaded config. Older gcode without
// it falls back to cm³ × default density.
func parseFilamentGrams(gcodePath string) (float64, error) {
	f, err := os.Open(gcodePath)
	if err != nil {
		return 0, err
	}
	defer f.Close()

	defaultDensity, _ := strconv.ParseFloat(filamentconfig.Defaults["filament_density"], 64)

	var grams, cm3 float64
	var haveGrams, haveCm3 bool

	scanner := bufio.NewScanner(f)
	scanner.Buffer(make([]byte, 0, 64*1024), 1<<20)
	for scanner.Scan() {
		line := scanner.Text()
		if v, ok := summaryValue(line, "; filament used [g] ="); ok {
			grams, haveGrams = v, true
		} else if v, ok := summaryValue(line, "; filament used [cm3] ="); ok {
			cm3, haveCm3 = v, true
		}
	}
	if err := scanner.Err(); err != nil {
		return 0, err
	}

	switch {
	case haveGrams && grams > 0:
		return grams, nil
	case haveCm3:
		return cm3 * defaultDensity, nil
	}
	return 0, fmt.Errorf("filament usage not found in gcode output")
}

func summaryValue(line, prefix string) (float64, bool) {
	if !strings.HasPrefix(line, prefix) {
		return 0, false
	}
	v, err := strconv.ParseFloat(strings.TrimSpace(strings.TrimPrefix(line, prefix)), 64)
	return v, err == nil
}
