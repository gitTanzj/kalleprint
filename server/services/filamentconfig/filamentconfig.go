package filamentconfig

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/gitTanzj/server/types"
)

// Keys is the whitelist of settings a filament INI may contain. Printer and
// print settings live in the static profiles (see profiles/) and cannot be
// overridden from the admin panel.
//
// pressure_advance is not a PrusaSlicer key: it is written to the INI as
// start_filament_gcode = "SET_PRESSURE_ADVANCE ADVANCE=<n>".
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
	"filament_max_volumetric_speed",
	"cooling",
	"fan_always_on",
	"min_fan_speed",
	"max_fan_speed",
	"bridge_fan_speed",
	"disable_fan_first_layers",
	"full_fan_speed_layer",
	"slowdown_below_layer_time",
	"min_print_speed",
	"filament_retract_length",
	"filament_retract_speed",
	"pressure_advance",
}

// Defaults are tuned for PLA on a stock Creality K1 hotend.
var Defaults = map[string]string{
	"filament_type":                 "PLA",
	"temperature":                   "210",
	"bed_temperature":               "60",
	"first_layer_temperature":       "215",
	"first_layer_bed_temperature":   "60",
	"filament_diameter":             "1.75",
	"extrusion_multiplier":          "0.98",
	"filament_density":              "1.24",
	"filament_cost":                 "0",
	"filament_max_volumetric_speed": "18",
	"cooling":                       "1",
	"fan_always_on":                 "1",
	"min_fan_speed":                 "100",
	"max_fan_speed":                 "100",
	"bridge_fan_speed":              "100",
	"disable_fan_first_layers":      "1",
	"full_fan_speed_layer":          "3",
	"slowdown_below_layer_time":     "8",
	"min_print_speed":               "15",
	"filament_retract_length":       "0.8",
	"filament_retract_speed":        "40",
	"pressure_advance":              "0.04",
}

type rule struct {
	min, max float64
	integer  bool
}

// rules clamp admin input to values that are safe for the K1. A typo here would
// otherwise reach the printer as e.g. a 2100 °C nozzle target.
var rules = map[string]rule{
	"temperature":                   {170, 300, false},
	"first_layer_temperature":       {170, 300, false},
	"bed_temperature":               {0, 110, false},
	"first_layer_bed_temperature":   {0, 110, false},
	"filament_diameter":             {1, 3, false},
	"extrusion_multiplier":          {0.5, 1.5, false},
	"filament_density":              {0.5, 3, false},
	"filament_cost":                 {0, 100000, false},
	"filament_max_volumetric_speed": {1, 30, false},
	"cooling":                       {0, 1, true},
	"fan_always_on":                 {0, 1, true},
	"min_fan_speed":                 {0, 100, true},
	"max_fan_speed":                 {0, 100, true},
	"bridge_fan_speed":              {0, 100, true},
	"disable_fan_first_layers":      {0, 10, true},
	"full_fan_speed_layer":          {0, 20, true},
	"slowdown_below_layer_time":     {0, 60, false},
	"min_print_speed":               {1, 100, false},
	"filament_retract_length":       {0, 10, false},
	"filament_retract_speed":        {5, 100, false},
	"pressure_advance":              {0, 0.2, false},
}

var filamentTypePattern = regexp.MustCompile(`^[A-Za-z0-9+-]{1,16}$`)

var pressureAdvancePattern = regexp.MustCompile(`SET_PRESSURE_ADVANCE\s+ADVANCE=([0-9.]+)`)

// Validate checks every whitelisted value in params. Unknown keys are ignored
// since BuildContent never writes them.
func Validate(params map[string]string) error {
	for _, k := range Keys {
		v, ok := params[k]
		if !ok {
			continue
		}
		if k == "filament_type" {
			if !filamentTypePattern.MatchString(v) {
				return fmt.Errorf("filament_type must be 1-16 letters, digits, '+' or '-'")
			}
			continue
		}
		r := rules[k]
		n, err := strconv.ParseFloat(v, 64)
		if err != nil || math.IsNaN(n) || math.IsInf(n, 0) {
			return fmt.Errorf("%s must be a number", k)
		}
		if r.integer && n != math.Trunc(n) {
			return fmt.Errorf("%s must be an integer", k)
		}
		if n < r.min || n > r.max {
			return fmt.Errorf("%s must be between %g and %g", k, r.min, r.max)
		}
	}
	if minFan, maxFan := params["min_fan_speed"], params["max_fan_speed"]; minFan != "" && maxFan != "" {
		lo, _ := strconv.ParseFloat(minFan, 64)
		hi, _ := strconv.ParseFloat(maxFan, 64)
		if lo > hi {
			return fmt.Errorf("min_fan_speed must not exceed max_fan_speed")
		}
	}
	return nil
}

// ParseIniFile reads a PrusaSlicer-style ini file and returns its key/value pairs.
// start_filament_gcode is translated back into pressure_advance.
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
		if key == "start_filament_gcode" {
			if m := pressureAdvancePattern.FindStringSubmatch(val); m != nil {
				result["pressure_advance"] = m[1]
			}
			continue
		}
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
		if k == "pressure_advance" {
			fmt.Fprintf(&b, "start_filament_gcode = \"SET_PRESSURE_ADVANCE ADVANCE=%s\"\n", params[k])
			continue
		}
		fmt.Fprintf(&b, "%s = %s\n", k, params[k])
	}
	return b.String()
}

// WriteFile deletes any existing file at path, then writes new content.
func WriteFile(path string, params map[string]string) error {
	_ = os.Remove(path)
	return os.WriteFile(path, []byte(BuildContent(params)), 0o644)
}

// WriteSlicerFile writes the filament layer used for slicing to a new temp file
// and returns its path. The stored INI (if any) is re-merged with Defaults so
// filaments created before a key was added still get a complete, whitelisted,
// validated config. The caller removes the file.
func WriteSlicerFile(iniFilePath string) (string, error) {
	base := map[string]string{}
	if iniFilePath != "" {
		parsed, err := ParseIniFile(iniFilePath)
		if err != nil {
			return "", fmt.Errorf("read filament config: %w", err)
		}
		base = parsed
	}
	params := MergeParams(base, nil)
	if err := Validate(params); err != nil {
		return "", fmt.Errorf("invalid filament config %s: %w", iniFilePath, err)
	}

	f, err := os.CreateTemp("", "filament-*.ini")
	if err != nil {
		return "", err
	}
	defer f.Close()
	if _, err := f.WriteString(BuildContent(params)); err != nil {
		os.Remove(f.Name())
		return "", err
	}
	return f.Name(), nil
}

// EnrichFilament reads the filament's ini file and populates the slicer config fields.
// Silently skips if the path is empty or the file cannot be read.
func EnrichFilament(f *types.Filament) {
	if f.IniFilePath == "" {
		return
	}
	parsed, err := ParseIniFile(f.IniFilePath)
	if err != nil {
		return
	}
	params := MergeParams(parsed, nil)
	f.FilamentType = params["filament_type"]
	f.Temperature = params["temperature"]
	f.BedTemperature = params["bed_temperature"]
	f.FirstLayerTemperature = params["first_layer_temperature"]
	f.FirstLayerBedTemperature = params["first_layer_bed_temperature"]
	f.FilamentDiameter = params["filament_diameter"]
	f.ExtrusionMultiplier = params["extrusion_multiplier"]
	f.FilamentDensity = params["filament_density"]
	f.FilamentCost = params["filament_cost"]
	f.FilamentMaxVolumetricSpeed = params["filament_max_volumetric_speed"]
	f.Cooling = params["cooling"]
	f.FanAlwaysOn = params["fan_always_on"]
	f.MinFanSpeed = params["min_fan_speed"]
	f.MaxFanSpeed = params["max_fan_speed"]
	f.BridgeFanSpeed = params["bridge_fan_speed"]
	f.DisableFanFirstLayers = params["disable_fan_first_layers"]
	f.FullFanSpeedLayer = params["full_fan_speed_layer"]
	f.SlowdownBelowLayerTime = params["slowdown_below_layer_time"]
	f.MinPrintSpeed = params["min_print_speed"]
	f.FilamentRetractLength = params["filament_retract_length"]
	f.FilamentRetractSpeed = params["filament_retract_speed"]
	f.PressureAdvance = params["pressure_advance"]
}
