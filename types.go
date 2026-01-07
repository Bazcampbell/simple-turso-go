// pkg/types.go

package simpleturso

type Config struct {
	DbUrl string
	DbKey string
}

type Statement struct {
	Q      string        `json:"q"`
	Params []interface{} `json:"params"`
}

type TursoRequest struct {
	Statements []Statement `json:"statements"`
}

type TursoResponse []TursoResult

type TursoResult struct {
	Results struct {
		Columns []string        `json:"columns"`
		Rows    [][]interface{} `json:"rows"`
	} `json:"results"`
}

type LogLevel string

const (
	LogLevelInfo    LogLevel = "info"
	LogLevelSuccess LogLevel = "success"
	LogLevelWarning LogLevel = "warning"
	LogLevelError   LogLevel = "error"
)
