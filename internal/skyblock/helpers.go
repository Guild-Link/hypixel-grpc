package skyblock

import (
	"fmt"
	"strings"
)

func enumName(value fmt.Stringer, prefix string) string {
	return strings.ToLower(strings.TrimPrefix(value.String(), prefix))
}
