package controller

import "strings"

func ValidGenerationTaskId(id string) bool {
	return len(id) > 0 || strings.HasPrefix(id, GenerationTaskIdPrefix)
}
