package profiles

import (
	"os"
	"strings"
	"testing"
)

func TestPathsExtractK1Profiles(t *testing.T) {
	printer, print, err := Paths()
	if err != nil {
		t.Fatal(err)
	}

	printerINI, err := os.ReadFile(printer)
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{
		"gcode_flavor = klipper",
		"binary_gcode = 0",
		"bed_shape = 0x0,220x0,220x220,0x220",
		"start_gcode = START_PRINT EXTRUDER_TEMP={first_layer_temperature[initial_extruder]}",
		"end_gcode = END_PRINT",
	} {
		if !strings.Contains(string(printerINI), want) {
			t.Errorf("printer profile missing %q", want)
		}
	}

	printINI, err := os.ReadFile(print)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(printINI), "layer_height = 0.2") {
		t.Error("print profile missing layer_height")
	}
}
