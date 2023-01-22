package resolve

import (
	"strings"
)

func validatePath(path string) bool {
	if path == "" {
		ErrorMessage("The path in the router cannot be empty.", 5)
		return true
	}
	if path[0] != '/' {
		ErrorMessage("The path in the router must start with \"/\"", 5)
		return true
	}

	if isParametric(path) {
		tokens := getTokens(path)
		for _, token := range tokens {
			if token[0] == ':' {
				if len(token) < 2 {
					ErrorMessage("The parameter parameter in the router path must have a name", 5)
					return true
				}
			}
		}
	}

	return false
}

func isParametric(path string) bool {
	return strings.Contains(path, ":")
}

func getPathTokens(path string) []pathToken {
	tokens := strings.Split(path, "/")
	if tokens[0] == "" {
		tokens = tokens[1:]
	}
	if tokens[len(tokens)-1] == "" {
		tokens = tokens[:len(tokens)-1]
	}

	pathTokens := make([]pathToken, len(tokens))
	for i, token := range tokens {
		parametricValue := ""
		if token[0] == ':' {
			parametricValue = token[1:]
		}

		pathTokens[i] = pathToken{
			IsParametric:    token[0] == ':' || token[0] == '*',
			Value:           token,
			ParametricValue: parametricValue,
		}
	}

	return pathTokens
}

func getTokens(path string) []string {
	tokens := strings.Split(path, "/")
	if tokens[0] == "" {
		tokens = tokens[1:]
	}
	if tokens[len(tokens)-1] == "" {
		tokens = tokens[:len(tokens)-1]
	}

	return tokens
}
