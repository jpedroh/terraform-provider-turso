package provider

import (
	"fmt"
	"strings"
)

func ExtractDbIdFromImportStateId(id string) (string, string, error) {
	parts := strings.SplitN(id, "/", 2)
	if len(parts) != 2 {
		return "", "", fmt.Errorf("%s is not a valid import state ID format. The expected format is: organizationSlug/databaseName", id)
	}
	return parts[0], parts[1], nil
}
