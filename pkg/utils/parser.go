package utils

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/charmbracelet/log"
)

// YamlDoc represents a YAML document structure
type YamlDoc struct {
	Data       map[string]string
	StringData map[string]string
}

// EnvVarObject represents a map of environment variables
type EnvVarObject map[string]string

// MergeDataFromManifests merges data from multiple YAML manifests
func MergeDataFromManifests(manifests []YamlDoc) EnvVarObject {
	envData := make(EnvVarObject)

	for _, manifest := range manifests {
		for k, v := range manifest.Data {
			envData[k] = v
		}
		for k, v := range manifest.StringData {
			envData[k] = v
		}
	}

	return envData
}

// GenerateEnvFile generates an environment file from the given EnvVarObject
func GenerateEnvFile(envObject EnvVarObject, filePath string, overwrite bool) {
	if _, err := os.Stat(filePath); err == nil {
		existingContent, err := os.ReadFile(filePath)
		if err != nil {
			log.Fatalf("error reading existing file: %v", err)
		}

		var hasDifferences bool

		// Read existing file
		scanner := bufio.NewScanner(strings.NewReader(string(existingContent)))
		for scanner.Scan() {
			line := scanner.Text()
			if len(line) != 0 && !strings.HasPrefix(line, "#") {
				parts := strings.SplitN(line, "=", 2)
				if len(parts) == 2 {
					key := strings.TrimSpace(parts[0])
					value := strings.TrimSpace(parts[1])
					// Remove quotes from value
					trimmedValue := strings.Trim(value, "'\"")
					if _, exists := envObject[key]; !exists {
						envObject[key] = trimmedValue
					} else if trimmedValue != envObject[key] {
						log.Warnf("%s has different values (existing: %q, new: %q)", key, trimmedValue, envObject[key])
						hasDifferences = true
					}
				}
			}
		}
		if hasDifferences && !overwrite {
			log.Fatalf("File already exists and --overwrite flag is not set")
		}

		if err := scanner.Err(); err != nil {
			log.Fatalf("Error scanning existing file: %v", err)
		}
	}

	var keys []string
	for k := range envObject {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var envContent strings.Builder
	for _, key := range keys {
		value := envObject[key]
		_, err := fmt.Fprintf(&envContent, "%s='%s'\n", key, value)
		if err != nil {
			log.Fatalf("error writing to string builder: %v", err)
		}
	}

	err := os.WriteFile(filePath, []byte(envContent.String()), 0644)
	if err != nil {
		log.Fatalf("error writing to file: %v", err)
	}
}
