package command

import (
	"errors"
	"fmt"

	"github.com/benskia/Thresher/internal/config"
	"github.com/benskia/Thresher/internal/power"
)

const setDescription string = `Usage: Thresher set <name>
	Activates profile <name> by writing its values to power_supply files.
`

func commandSet(cfg *config.Config, args []string) error {
	if len(args) < 1 {
		return errors.New("set expects one arg: <name>")
	}

	name := args[0]

	profile, ok := cfg.Profiles[name]
	if !ok {
		return fmt.Errorf("profile %s not found", name)
	}

	fmt.Printf("Setting profile %s ...\n", profile.Name)
	if err := power.SaveThresholds(profile); err != nil {
		return fmt.Errorf("error saving thresholds: %w", err)
	}

	batteries, err := power.GetThresholds()
	if err != nil {
		return err
	}

	fmt.Println("\nCurrent Thresholds:")
	for _, battery := range batteries {
		fmt.Printf("\tName: %s\n", battery.Name)
		fmt.Printf("\tStart: %d\tEnd: %d\n\n", battery.Start, battery.End)
	}

	return nil
}
