package scanner

import (
	"bufio"
	"os"
	"strings"
)

func LoadResolvers(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var resolvers []string
	sc := bufio.NewScanner(f)

	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		resolvers = append(resolvers, line)
	}

	if err := sc.Err(); err != nil {
		return nil, err
	}
	return resolvers, nil
}
