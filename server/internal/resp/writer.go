package resp

import "fmt"

func WriteSimpleString(value string) string {
	return "+" + value + "\r\n"
}

func WriteError(message string) string {
	return "-" + message + "\r\n"
}

func WriteInteger(value int) string {
	return fmt.Sprintf(":%d\r\n", value)
}

func WriteBulkString(value string) string {
	return fmt.Sprintf("$%d\r\n%s\r\n", len(value), value)
}

func WriteNullBulkString() string {
	return "$-1\r\n"
}

// WriteArrayHeader writes only the "*<count>" line. Use it when the items
// are added one by one after it.
func WriteArrayHeader(count int) string {
	return fmt.Sprintf("*%d\r\n", count)
}

func WriteArray(values []string) string {
	result := WriteArrayHeader(len(values))
	for _, value := range values {
		result += WriteBulkString(value)
	}
	return result
}

func WriteNullArray() string {
	return "*-1\r\n"
}
