package werf

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"strings"

	parser "github.com/gbh-tech/envi/pkg/utils"

	"gopkg.in/yaml.v3"
)

type WerfOptions struct {
	Environment string
	Secrets     bool
	Values      string
	Path        []string
}

func GenerateEnvFile(options WerfOptions) error {
	environment := strings.TrimSpace(options.Environment)
	werfCommand := []string{
		"werf",
		"render",
		"--env",
		environment,
		"--values",
		fmt.Sprintf(".helm/values/%s.yaml", environment),
		"--dev",
	}

	if options.Secrets {
		werfCommand = append(werfCommand, "--secret-values", fmt.Sprintf(".helm/secrets/%s.yaml", environment))
	}

	if options.Values != "" {
		extraVars := strings.TrimSpace(options.Values)
		werfCommand = append(werfCommand, "--set", extraVars)
	}

	log.Println("Werf command: ", strings.Join(werfCommand, " "))
	cmd := exec.Command(werfCommand[0], werfCommand[1:]...)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		log.Printf("Command stdout:\n%s", stdout.String())
		log.Printf("Command stderr:\n%s", stderr.String())
		return fmt.Errorf("failed to execute Werf command: %w\nStderr: %s", err, stderr.String())
	}
	renderedManifests := stdout.Bytes()

	log.Println("Obtaining env vars from rendered manifests...")

	var manifests []parser.YamlDoc
	decoder := yaml.NewDecoder(strings.NewReader(string(renderedManifests)))
	for {
		var doc parser.YamlDoc
		if err := decoder.Decode(&doc); err != nil {
			if err.Error() == "EOF" {
				break
			}
			return fmt.Errorf("failed to decode YAML: %v", err)
		}
		manifests = append(manifests, doc)
	}

	envData := parser.MergeDataFromManifests(manifests)

	for _, path := range options.Path {
		if err := parser.GenerateEnvFile(envData, path); err != nil {
			return fmt.Errorf("failed to generate env file at %s: %v", path, err)
		}
	}

	return nil
}
