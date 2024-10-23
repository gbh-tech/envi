package manual

import (
	"fmt"
	"strings"

	parser "github.com/gbh-tech/envi/pkg/utils"
)

type ManualOptions struct {
	Value []string
	Path  []string
}

func GenerateEnvFile(options ManualOptions) error {
	envData := make(parser.EnvVarObject)
	for _, value := range options.Value {
		kp := strings.SplitN(value, "=", 2)
		envData[kp[0]] = kp[1]
	}

	for _, path := range options.Path {
		if err := parser.GenerateEnvFile(envData, path); err != nil {
			return fmt.Errorf("failed to generate env file at %s: %v", path, err)
		}
	}

	return nil
}
