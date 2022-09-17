package printer

import (
	"os"

	"github.com/fatih/color"
)

var (
	ActiveItemColor = color.New(color.FgGreen, color.Bold)
)

func init() {
	EnableOrDisableColor(ActiveItemColor)
}

// useColors returns true if colors are force-enabled,
// false if colors are disabled, or nil for default behavior
// which is determined based on factors like if stdout is tty.
func useColors() *bool {
	tr, fa := true, false
	if os.Getenv("KUBE_FORCE_COLOR") != "" {
		return &tr
	} else if os.Getenv("NO_COLOR") != "" {
		return &fa
	}
	return nil
}

// EnableOrDisableColor determines if color should be force-enabled or force-disabled
// or left untouched based on environment configuration.
func EnableOrDisableColor(c *color.Color) {
	if v := useColors(); v != nil && *v {
		c.EnableColor()
	} else if v != nil && !*v {
		c.DisableColor()
	}
}
