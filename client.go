// pkg/client.go

package simpleturso

import (
	"fmt"
	"strings"
	"time"
)

func LogToTurso(level LogLevel, message, application, timezone string, keysAndValues ...interface{}) error {
	cfg, err := getConfig()
	if err != nil {
		return err
	}

	location, err := time.LoadLocation(timezone)
	if err != nil {
		return err
	}
	timestamp := time.Now().In(location).Format("02/01/2006 15:04:05")

	var messageParts []string
	messageParts = append(messageParts, message)

	for i := 0; i < len(keysAndValues); i += 2 {
		if i+1 >= len(keysAndValues) {
			// Odd number of arguments, skip the last one
			break
		}

		key := fmt.Sprintf("%v", keysAndValues[i])
		value := keysAndValues[i+1]

		messageParts = append(messageParts, fmt.Sprintf("%s=%v", key, value))
	}
	fullMessage := strings.Join(messageParts, ", ")

	req := TursoRequest{
		Statements: []Statement{
			{
				Q:      "INSERT INTO logs (timestamp, application, level, message) VALUES (?, ?, ?, ?)",
				Params: []interface{}{timestamp, application, level, fullMessage},
			},
		},
	}

	_, err = executeTursoRequestWithResponse(req, cfg)
	return err
}

// Select executes a SELECT query and scans results into type T
func Select[T any](query string, params []interface{}, scan func(map[string]interface{}) (T, error)) ([]T, error) {
	cfg, err := getConfig()
	if err != nil {
		return nil, err
	}

	if params == nil {
		params = []interface{}{}
	}

	req := TursoRequest{
		Statements: []Statement{
			{
				Q:      query,
				Params: params,
			},
		},
	}

	tursoResp, err := executeTursoRequestWithResponse(req, cfg)
	if err != nil {
		return nil, err
	}

	// 1. Check if the slice itself is empty
	if tursoResp == nil || len(*tursoResp) == 0 {
		return []T{}, nil
	}

	// 2. Access the first result in the array: (*tursoResp)[0]
	// Then drill into the .Results field
	data := (*tursoResp)[0].Results
	var items []T

	for _, row := range data.Rows {
		rowMap := make(map[string]interface{})
		for i, col := range data.Columns {
			if i < len(row) {
				rowMap[col] = row[i]
			}
		}

		item, err := scan(rowMap)
		if err != nil {
			return nil, fmt.Errorf("scan error: %w", err)
		}
		items = append(items, item)
	}

	return items, nil
}

func Execute(query string, params []interface{}) error {
	cfg, err := getConfig()
	if err != nil {
		return err
	}

	if params == nil {
		params = []interface{}{}
	}

	req := TursoRequest{
		Statements: []Statement{
			{
				Q:      query,
				Params: params,
			},
		},
	}

	_, err = executeTursoRequestWithResponse(req, cfg)
	return err
}
