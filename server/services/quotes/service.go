package quotes

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/gitTanzj/server/config"
)

const PROFIT_MULTIPLIER = 1.8

func CalculateQuote(amount float64, filamentCost float64) float64 {
	price := amount * filamentCost * PROFIT_MULTIPLIER
	price_rounded := math.Round(price*100) / 100
	return math.Max(price_rounded, 1)
}

// SliceFilament slices the input file and returns grams of filament used.
// If gcodeOut is non-empty the gcode is written there and kept; otherwise a
// temporary file is used and removed automatically.
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

	args := []string{"--slice", "--output", gcodeOut}
	if !strings.EqualFold(ext, ".3mf") {
		cfg := iniFilePath
		if cfg == "" {
			cfg = config.Envs.PrusaSlicerConfig
		}
		args = append(args, "--load", cfg)
	}
	args = append(args, inputPath)

	if out, err := exec.Command(config.Envs.PrusaSlicerPath, args...).CombinedOutput(); err != nil {
		return 0, fmt.Errorf("slicer failed: %w: %s", err, out)
	}

	return parseFilamentGrams(gcodeOut)
}

func parseFilamentGrams(gcodePath string) (float64, error) {
	f, err := os.Open(gcodePath)
	if err != nil {
		return 0, err
	}
	defer f.Close()

	const defaultDensity = 1.25

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "; filament used [cm3] =") {
			parts := strings.SplitN(line, "=", 2)
			if len(parts) == 2 {
				cm3, err := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
				if err != nil {
					return 0, err
				}
				return cm3 * defaultDensity, nil
			}
		}
	}
	return 0, fmt.Errorf("filament usage not found in gcode output")
}
