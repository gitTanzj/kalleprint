package quotes

import (
	"archive/zip"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gitTanzj/server/config"
)

// cubeSTL returns an ASCII STL of a size×size×size mm cube at the origin.
func cubeSTL(size float64) string {
	v := func(x, y, z float64) string { return fmt.Sprintf("vertex %g %g %g\n", x*size, y*size, z*size) }
	quads := [][4][3]float64{
		{{0, 0, 0}, {0, 1, 0}, {1, 1, 0}, {1, 0, 0}},
		{{0, 0, 1}, {1, 0, 1}, {1, 1, 1}, {0, 1, 1}},
		{{0, 0, 0}, {1, 0, 0}, {1, 0, 1}, {0, 0, 1}},
		{{0, 1, 0}, {0, 1, 1}, {1, 1, 1}, {1, 1, 0}},
		{{0, 0, 0}, {0, 0, 1}, {0, 1, 1}, {0, 1, 0}},
		{{1, 0, 0}, {1, 1, 0}, {1, 1, 1}, {1, 0, 1}},
	}
	var b strings.Builder
	b.WriteString("solid cube\n")
	for _, q := range quads {
		for _, tri := range [][3]int{{0, 1, 2}, {0, 2, 3}} {
			b.WriteString("facet normal 0 0 0\nouter loop\n")
			for _, i := range tri {
				b.WriteString(v(q[i][0], q[i][1], q[i][2]))
			}
			b.WriteString("endloop\nendfacet\n")
		}
	}
	b.WriteString("endsolid cube\n")
	return b.String()
}

// TestSliceFilamentWithPrusaSlicer runs the real slicer; it is skipped when
// PrusaSlicer is not installed (e.g. outside Docker).
func TestSliceFilamentWithPrusaSlicer(t *testing.T) {
	if _, err := exec.LookPath(config.Envs.PrusaSlicerPath); err != nil {
		t.Skipf("PrusaSlicer not available at %s", config.Envs.PrusaSlicerPath)
	}

	dir := t.TempDir()
	stl := filepath.Join(dir, "cube.stl")
	if err := os.WriteFile(stl, []byte(cubeSTL(20)), 0o644); err != nil {
		t.Fatal(err)
	}

	gcodePath := filepath.Join(dir, "cube.gcode")
	grams, err := SliceFilament(stl, ".stl", "", gcodePath)
	if err != nil {
		t.Fatal(err)
	}
	if grams < 1 || grams > 20 {
		t.Errorf("unexpected filament usage for a 20 mm cube: %.2f g", grams)
	}

	data, err := os.ReadFile(gcodePath)
	if err != nil {
		t.Fatal(err)
	}
	gcode := string(data)
	for _, want := range []string{
		"START_PRINT EXTRUDER_TEMP=215 BED_TEMP=60",
		"END_PRINT",
		"SET_PRESSURE_ADVANCE ADVANCE=0.04",
		"; gcode_flavor = klipper",
		"; bed_shape = 0x0,220x0,220x220,0x220",
		"; filament_max_volumetric_speed = 18",
		"; layer_height = 0.2",
	} {
		if !strings.Contains(gcode, want) {
			t.Errorf("gcode missing %q", want)
		}
	}
	if strings.Contains(gcode, "M190") || strings.Contains(gcode, "M109") {
		t.Error("PrusaSlicer emitted its own heating commands; START_PRINT should own heating")
	}

	t.Run("3mf embedded config is ignored", func(t *testing.T) {
		exported := filepath.Join(dir, "exported.3mf")
		out, err := exec.Command(config.Envs.PrusaSlicerPath, "--export-3mf", "--output", exported, stl).CombinedOutput()
		if err != nil {
			t.Fatalf("export 3mf: %v: %s", err, out)
		}
		// Embed a conflicting project config, as a 3MF saved from the GUI would carry.
		threeMF := filepath.Join(dir, "cube.3mf")
		if err := withEmbeddedConfig(exported, threeMF, "; layer_height = 0.3\n; temperature = 250\n"); err != nil {
			t.Fatal(err)
		}

		gcodePath := filepath.Join(dir, "cube3mf.gcode")
		if _, err := SliceFilament(threeMF, ".3mf", "", gcodePath); err != nil {
			t.Fatal(err)
		}
		data, _ := os.ReadFile(gcodePath)
		gcode := string(data)
		if !strings.Contains(gcode, "; layer_height = 0.2\n") || !strings.Contains(gcode, "; temperature = 210\n") {
			t.Error("embedded 3mf config overrode the K1 profiles")
		}
	})
}

func withEmbeddedConfig(src, dst, cfg string) error {
	zr, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer zr.Close()
	f, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer f.Close()
	zw := zip.NewWriter(f)
	for _, e := range zr.File {
		if err := zw.Copy(e); err != nil {
			return err
		}
	}
	w, err := zw.Create("Metadata/Slic3r_PE.config")
	if err != nil {
		return err
	}
	if _, err := w.Write([]byte(cfg)); err != nil {
		return err
	}
	return zw.Close()
}
