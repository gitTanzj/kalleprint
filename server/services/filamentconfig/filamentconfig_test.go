package filamentconfig

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestDefaultsAreValid(t *testing.T) {
	for _, k := range Keys {
		if _, ok := Defaults[k]; !ok {
			t.Errorf("missing default for %s", k)
		}
	}
	if err := Validate(Defaults); err != nil {
		t.Fatalf("defaults should validate: %v", err)
	}
}

func TestValidateRejects(t *testing.T) {
	cases := map[string]map[string]string{
		"nozzle too hot":    {"temperature": "2100"},
		"bed too hot":       {"bed_temperature": "150"},
		"not a number":      {"filament_max_volumetric_speed": "fast"},
		"fractional fan":    {"max_fan_speed": "50.5"},
		"newline injection": {"temperature": "210\nstart_gcode = M112"},
		"bad filament type": {"filament_type": "PLA\nbed_shape = 0x0"},
		"min fan above max": {"min_fan_speed": "80", "max_fan_speed": "50"},
		"pressure advance":  {"pressure_advance": "1"},
		"not finite":        {"temperature": "NaN"},
	}
	for name, overrides := range cases {
		t.Run(name, func(t *testing.T) {
			if err := Validate(MergeParams(Defaults, overrides)); err == nil {
				t.Fatal("expected validation error")
			}
		})
	}
}

func TestBuildAndParseRoundTrip(t *testing.T) {
	params := MergeParams(Defaults, map[string]string{"pressure_advance": "0.035", "temperature": "205"})
	content := BuildContent(params)

	if strings.Contains(content, "pressure_advance") {
		t.Error("pressure_advance must not be written as a slicer key")
	}
	if !strings.Contains(content, `start_filament_gcode = "SET_PRESSURE_ADVANCE ADVANCE=0.035"`) {
		t.Errorf("missing pressure advance gcode in:\n%s", content)
	}

	path := filepath.Join(t.TempDir(), "f.ini")
	if err := WriteFile(path, params); err != nil {
		t.Fatal(err)
	}
	parsed, err := ParseIniFile(path)
	if err != nil {
		t.Fatal(err)
	}
	for _, k := range Keys {
		if parsed[k] != params[k] {
			t.Errorf("%s: got %q, want %q", k, parsed[k], params[k])
		}
	}
}

func TestWriteSlicerFile(t *testing.T) {
	t.Run("no stored ini uses defaults", func(t *testing.T) {
		path, err := WriteSlicerFile("")
		if err != nil {
			t.Fatal(err)
		}
		defer os.Remove(path)
		parsed, _ := ParseIniFile(path)
		if parsed["filament_max_volumetric_speed"] != Defaults["filament_max_volumetric_speed"] {
			t.Errorf("expected default volumetric speed, got %q", parsed["filament_max_volumetric_speed"])
		}
	})

	t.Run("legacy ini is completed and stripped of unknown keys", func(t *testing.T) {
		legacy := filepath.Join(t.TempDir(), "legacy.ini")
		os.WriteFile(legacy, []byte("temperature = 230\nbed_shape = 0x0,10x0,10x10\n"), 0o644)

		path, err := WriteSlicerFile(legacy)
		if err != nil {
			t.Fatal(err)
		}
		defer os.Remove(path)
		data, _ := os.ReadFile(path)
		content := string(data)
		if !strings.Contains(content, "temperature = 230\n") {
			t.Error("stored value lost")
		}
		if strings.Contains(content, "bed_shape") {
			t.Error("non-filament key leaked into slicer config")
		}
		if !strings.Contains(content, "min_fan_speed = 100\n") {
			t.Error("missing default for new key")
		}
	})

	t.Run("invalid stored ini is rejected", func(t *testing.T) {
		bad := filepath.Join(t.TempDir(), "bad.ini")
		os.WriteFile(bad, []byte("temperature = 900\n"), 0o644)
		if _, err := WriteSlicerFile(bad); err == nil {
			t.Fatal("expected error")
		}
	})
}
