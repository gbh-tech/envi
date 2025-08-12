package werf

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"

	"github.com/charmbracelet/log"
	parser "github.com/gbh-tech/envi/pkg/utils"

	"gopkg.in/yaml.v3"
)

func GenerateEnvFile(options Options) {
	environment := strings.TrimSpace(options.Environment)

	secretFile := fmt.Sprintf(".helm/secrets/%s.yaml", environment)
	valueFile := fmt.Sprintf(".helm/values/%s.yaml", environment)

	werfCommand := []string{
		"werf",
		"render",
		"--env",
		environment,
		"--values",
		valueFile,
	}

	if options.Dir != "" {
		werfCommand = append(
			werfCommand,
			"--dir",
			options.Dir,
		)
	}

	if options.Secrets {
		werfCommand = append(
			werfCommand,
			"--secret-values",
			secretFile,
		)
	}

	if options.Development {
		werfCommand = append(
			werfCommand,
			"--dev",
		)
	}

	if len(options.ValueFiles) > 0 {
		for _, file := range options.ValueFiles {
			werfCommand = append(
				werfCommand,
				"--values",
				file,
			)
		}
	}

	if options.Values != "" {
		extraVars := strings.TrimSpace(options.Values)
		werfCommand = append(
			werfCommand,
			"--set",
			extraVars,
		)
	}

	log.Infof("Werf command: %s", strings.Join(werfCommand, " "))
	cmd := exec.Command(werfCommand[0], werfCommand[1:]...)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		log.Printf("Command stdout:\n%s", stdout.String())
		log.Fatalf("Failed to execute Werf command: %s\nStderr: %s", err, stderr.String())
	}

	renderedManifests := stdout.Bytes()

	log.Infof("Obtaining env vars from rendered manifests...")

	var manifests []parser.YamlDoc
	decoder := yaml.NewDecoder(strings.NewReader(string(renderedManifests)))
	for {
		var doc parser.YamlDoc
		if err := decoder.Decode(&doc); err != nil {
			if err.Error() == "EOF" {
				break
			}
			log.Fatalf("Failed to decode YAML: %v", err)
		}
		manifests = append(manifests, doc)
	}

	envData := parser.MergeDataFromManifests(manifests)

	for _, path := range options.Path {
		parser.GenerateEnvFile(envData, path, options.Overwrite)
		log.Infof("File generated in %s using Werf!", path)
	}
}
