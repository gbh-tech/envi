package manual

import (
	"github.com/charmbracelet/log"
	"strings"

	"github.com/gbh-tech/envi/pkg/utils"
)

func generateEnvFileFromManualInput(options Options) {
	envData := make(utils.EnvVarObject)
	for _, value := range options.Values {
		keyPair := strings.SplitN(value, "=", 2)
		envData[keyPair[0]] = keyPair[1]
	}

	for _, path := range options.Paths {
		if err := utils.GenerateEnvFile(envData, path); err != nil {
			log.Fatalf("Failed to generate env file at %s: %v", path, err)
		}
	}
}
