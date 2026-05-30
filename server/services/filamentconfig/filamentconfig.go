package filamentconfig

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/gitTanzj/server/types"
)

var Keys = []string{
	"filament_type",
	"temperature",
	"bed_temperature",
	"first_layer_temperature",
	"first_layer_bed_temperature",
	"filament_diameter",
	"extrusion_multiplier",
	"filament_density",
	"filament_cost",
}

var Defaults = map[string]string{
	"filament_type":               "PLA",
	"temperature":                 "220",
	"bed_temperature":             "60",
	"first_layer_temperature":     "225",
	"first_layer_bed_temperature": "65",
	"filament_diameter":           "1.75",
	"extrusion_multiplier":        "1",
	"filament_density":            "1.24",
	"filament_cost":               "0",
}

// ParseIniFile reads a PrusaSlicer-style ini file and returns its key/value pairs.
// Only keys listed in Keys are returned; unknown lines are ignored.
func ParseIniFile(path string) (map[string]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	result := make(map[string]string, len(Keys))
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, ";") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		val := strings.TrimSpace(parts[1])
		result[key] = val
	}
	return result, scanner.Err()
}

// MergeParams returns a complete params map by starting from base and overriding
// with any non-empty values from overrides. Missing keys fall back to Defaults.
func MergeParams(base, overrides map[string]string) map[string]string {
	merged := make(map[string]string, len(Keys))
	for _, k := range Keys {
		if v, ok := overrides[k]; ok && v != "" {
			merged[k] = v
		} else if v, ok := base[k]; ok && v != "" {
			merged[k] = v
		} else {
			merged[k] = Defaults[k]
		}
	}
	return merged
}

// BuildContent serialises params into ini format in canonical key order.
func BuildContent(params map[string]string) string {
	var b strings.Builder
	for _, k := range Keys {
		fmt.Fprintf(&b, "%s = %s\n", k, params[k])
	}
	return b.String()
}

// WriteFile deletes any existing file at path, then writes new content.
func WriteFile(path string, params map[string]string) error {
	_ = os.Remove(path)
	return os.WriteFile(path, []byte(BuildContent(params)), 0o644)
}

// EnrichFilament reads the filament's ini file and populates the slicer config fields.
// Silently skips if the path is empty or the file cannot be read.
func EnrichFilament(f *types.Filament) {
	if f.IniFilePath == "" {
		return
	}
	params, err := ParseIniFile(f.IniFilePath)
	if err != nil {
		return
	}
	f.FilamentType = params["filament_type"]
	f.Temperature = params["temperature"]
	f.BedTemperature = params["bed_temperature"]
	f.FirstLayerTemperature = params["first_layer_temperature"]
	f.FirstLayerBedTemperature = params["first_layer_bed_temperature"]
	f.FilamentDiameter = params["filament_diameter"]
	f.ExtrusionMultiplier = params["extrusion_multiplier"]
	f.FilamentDensity = params["filament_density"]
	f.FilamentCost = params["filament_cost"]
}
