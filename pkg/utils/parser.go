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
	StringData map[string]string `yaml:"stringData,omitempty"`
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

	// If the file exists, copy envObject to avoid modifying the original and merge with its contents
	if _, err := os.Stat(filePath); err == nil {
		envCopy = copyEnvObject(envObject)
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

// copyEnvObject creates a shallow copy of the original environment map.
// This prevents mutations to the original input when merging in existing file data.
func copyEnvObject(src EnvVarObject) EnvVarObject {
	dst := make(EnvVarObject, len(src))
	for k, v := range src {
		dst[k] = v
	}
	return dst
}

// parseAndMergeExistingEnv reads an existing .env file and merges values into the provided map.
// - Preserves values not present in the new manifest.
// - Warns on mismatched keys to allow manual resolution or require --overwrite.
func parseAndMergeExistingEnv(filePath string, env EnvVarObject) bool {
	content, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Error reading existing file: %v", err)
	}

	scanner := bufio.NewScanner(strings.NewReader(string(content)))
	hasDifferences := false

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		// Strip quotes and trim spaces from value to normalize formatting
		value := strings.TrimSpace(strings.Trim(parts[1], `"'`))

		if existing, exists := env[key]; !exists {
			// Preserve value from existing file if it's not defined in the manifest
			env[key] = value
		} else if existing != value {
			// Detect conflict in values to warn user
			log.Warnf("%s has different values (existing: %q, new: %q)", key, value, existing)
			hasDifferences = true
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("Error scanning existing file: %v", err)
	}

	return hasDifferences
}

// writeEnvFile writes the given environment map to a file in sorted key order.
// This ensures deterministic output and avoids unnecessary diffs in version control.
func writeEnvFile(filePath string, env EnvVarObject) {
	var keys []string
	for k := range env {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var builder strings.Builder
	for _, key := range keys {
		// Always quote values to avoid ambiguity and preserve consistency
		_, err := fmt.Fprintf(&builder, "%s='%s'\n", key, env[key])
		if err != nil {
			log.Fatalf("Error writing to buffer: %v", err)
		}
	}

	if err := os.WriteFile(filePath, []byte(builder.String()), 0644); err != nil {
		log.Fatalf("Error writing to file: %v", err)
	}
}
