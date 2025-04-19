package sources

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type PropertiesSource struct {
	values map[string]string
}

func NewPropertiesFileSource(path string) (*PropertiesSource, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	values := make(map[string]string)
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip comments and empty lines
		if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, ";") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("malformed line: %s", line)
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		values[key] = value
	}

	err = scanner.Err()
	if err != nil {
		return nil, fmt.Errorf("failed to parse properties file %w", err)
	}

	return &PropertiesSource{values: values}, nil
}

func (s *PropertiesSource) Get(key string) (string, bool) {
	val, ok := s.values[key]
	return val, ok
}
