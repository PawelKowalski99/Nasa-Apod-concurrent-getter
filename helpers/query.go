package helpers

import (
	"fmt"
	"github.com/tidwall/gjson"
	"strings"
)

// GetValidJsonField lets query json picture due to its fields.
// If query got url == XXX AND copyright == XXX it checks if is right
func GetValidJsonField(json string, query map[string][]string, field string) string {
	queryCounter := 0
	for key, values := range query {

		connectedValues := strings.Join(values, " ")
		if strings.ReplaceAll(gjson.Get(json, key).Raw, `"`, "") == connectedValues {
			queryCounter++
		}
	}
	fmt.Println(queryCounter)
	if queryCounter == len(query) {
		return gjson.Get(json, field).String()
	}
	return ""
}
