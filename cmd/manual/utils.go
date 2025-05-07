package manual

import (
	"strings"

	"github.com/charmbracelet/log"

	"github.com/gbh-tech/envi/pkg/utils"
)

func generateEnvFileFromManualInput(options Options) {
	envData := make(utils.EnvVarObject)
	for _, value := range options.Values {
		keyPair := strings.SplitN(value, "=", 2)
		envData[keyPair[0]] = keyPair[1]
	}

	for _, path := range options.Paths {
		utils.GenerateEnvFile(envData, path, options.Overwrite)
		log.Infof("File generated in %s, from manual input!", path)
	}
}
