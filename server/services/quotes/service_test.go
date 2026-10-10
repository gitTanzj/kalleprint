package quotes

import (
	"archive/zip"
	"bytes"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestBuildSlicerArgsLoadOrder(t *testing.T) {
	got := buildSlicerArgs("model.stl", "out.gcode", "printer.ini", "print.ini", "filament.ini")
	want := []string{
		"--export-gcode", "--output", "out.gcode",
		"--load", "printer.ini", "--load", "print.ini", "--load", "filament.ini",
		"--center", "110,110", "model.stl",
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v\nwant %v", got, want)
	}
}

func TestStrip3MFConfig(t *testing.T) {
	path := filepath.Join(t.TempDir(), "in.3mf")
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)
	for _, name := range []string{
		"[Content_Types].xml",
		"3D/3dmodel.model",
		"Metadata/Slic3r_PE.config",
		"Metadata/Slic3r_PE_model.config",
		"Metadata/project_settings.config",
	} {
		w, _ := zw.Create(name)
		w.Write([]byte(name))
	}
	zw.Close()
	os.WriteFile(path, buf.Bytes(), 0o644)

	out, err := strip3MFConfig(path)
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(out)

	zr, err := zip.OpenReader(out)
	if err != nil {
		t.Fatal(err)
	}
	defer zr.Close()
	var names []string
	for _, f := range zr.File {
		names = append(names, f.Name)
	}
	want := []string{"[Content_Types].xml", "3D/3dmodel.model"}
	if !reflect.DeepEqual(names, want) {
		t.Fatalf("got %v, want %v", names, want)
	}
}

func TestParseFilamentGrams(t *testing.T) {
	cases := map[string]struct {
		gcode string
		want  float64
	}{
		"prefers grams":       {"; filament used [mm] = 1000\n; filament used [cm3] = 10\n; filament used [g] = 12.34\n", 12.34},
		"falls back to cm3":   {"; filament used [cm3] = 10\n", 12.4},
		"zero grams uses cm3": {"; filament used [cm3] = 10\n; filament used [g] = 0\n", 12.4},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "out.gcode")
			os.WriteFile(path, []byte("G1 X0\n"+c.gcode), 0o644)
			got, err := parseFilamentGrams(path)
			if err != nil {
				t.Fatal(err)
			}
			if diff := got - c.want; diff > 1e-9 || diff < -1e-9 {
				t.Fatalf("got %v, want %v", got, c.want)
			}
		})
	}

	path := filepath.Join(t.TempDir(), "empty.gcode")
	os.WriteFile(path, []byte("G1 X0\n"), 0o644)
	if _, err := parseFilamentGrams(path); err == nil {
		t.Fatal("expected error when usage is missing")
	}
}
