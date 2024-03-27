package logs

import "fmt"

const (
	JsonLogFormat    = "json"
	ConsoleLogFormat = "console"
)

var DefaultLogFormat = JsonLogFormat

// ConvertToZapFormat converts and validated log format string.
func ConvertToZapFormat(format string) (string, error) {
	switch format {
	case ConsoleLogFormat:
		return ConsoleLogFormat, nil
	case JsonLogFormat:
		return JsonLogFormat, nil
	case "":
		return DefaultLogFormat, nil
	default:
		return "", fmt.Errorf("unknown log format: %s, supported values json, console", format)
	}
}
