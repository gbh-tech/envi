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
	var envCopy EnvVarObject

	if _, err := os.Stat(filePath); err == nil {
		// Defensive copy to avoid mutating the original map
		envCopy = copyEnvObject(envObject)
		// Parse the existing .env file line by line
		hasDifferences := parseAndMergeExistingEnv(filePath, envCopy)

		if hasDifferences && !overwrite {
			log.Fatalf("Environment file differs from manifest; use the --overwrite flag to replace it after saving any required values")
		}
	} else {
		// No existing file: just use original
		envCopy = envObject
	}

	writeEnvFile(filePath, envCopy)
}

func copyEnvObject(src EnvVarObject) EnvVarObject {
	dst := make(EnvVarObject, len(src))
	for k, v := range src {
		dst[k] = v
	}
	return dst
}

func parseAndMergeExistingEnv(filePath string, env EnvVarObject) bool {
	content, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Error reading existing file: %v", err)
	}

	scanner := bufio.NewScanner(strings.NewReader(string(content)))
	hasDifferences := false

	// Parse the existing .env file line by line
	for scanner.Scan() {
		line := scanner.Text()

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Split on first '=' only to handle values that may contain '='
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		// Strip surrounding single or double quotes from the value
		value := strings.TrimSpace(strings.Trim(parts[1], `"'`))

		// Handle two cases:
		// 1. Key doesn't exist in new manifest - preserve existing value
		// 2. Key exists but values differ - warn user and set conflict flag
		if existing, exists := env[key]; !exists {
			env[key] = value
		} else if existing != value {
			log.Warnf("%s has different values (existing: %q, new: %q)", key, value, existing)
			hasDifferences = true
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("Error scanning existing file: %v", err)
	}

	return hasDifferences
}

func writeEnvFile(filePath string, env EnvVarObject) {
	var keys []string
	for k := range env {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var builder strings.Builder
	for _, key := range keys {
		_, err := fmt.Fprintf(&builder, "%s='%s'\n", key, env[key])
		if err != nil {
			log.Fatalf("Error writing to buffer: %v", err)
		}
	}

	if err := os.WriteFile(filePath, []byte(builder.String()), 0644); err != nil {
		log.Fatalf("Error writing to file: %v", err)
	}
}
