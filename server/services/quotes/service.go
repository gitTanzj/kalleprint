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

var profitMultiplier = 1.2

func CalculateQuote(amount float64, filamentCost float64) float64 {
	price := amount * filamentCost * profitMultiplier
	return math.Round(price*100) / 100
}

func SliceFilament(inputPath, ext string) (float64, error) {
	tmpGcode, err := os.CreateTemp("", "quote-*.gcode")
	if err != nil {
		return 0, err
	}
	tmpGcode.Close()
	defer os.Remove(tmpGcode.Name())

	args := []string{"--slice", "--output", tmpGcode.Name()}
	if !strings.EqualFold(ext, ".3mf") {
		args = append(args, "--load", config.Envs.PrusaSlicerConfig)
	}
	args = append(args, inputPath)

	if out, err := exec.Command(config.Envs.PrusaSlicerPath, args...).CombinedOutput(); err != nil {
		return 0, fmt.Errorf("slicer failed: %w: %s", err, out)
	}

	return parseFilamentGrams(tmpGcode.Name())
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
