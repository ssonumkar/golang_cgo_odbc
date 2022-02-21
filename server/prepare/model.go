package prepare

import (
	"github.com/ssonumkar/test-odbc-snowflake/api"
)

type QueryMetadata struct {
	InputDT  *Input  `json:"input,omitempty"`
	OutputDT *Output `json:"output,omitempty"`
}
type Input struct {
	numInput  int        //`json:"numinput"`
	InputSets []InputSet `json:"inputset,omitempty"`
}
type InputSet struct {
	Parameters Parameter `json:"parameter"`
}

// Parameter ...
type Parameter struct {
	Name     string          `json:"name"`
	DataType string          `json:"datatype"`
	Nullable bool            `json:"nullable"`
	SqlType  api.SQLSMALLINT `json:"sqltype"`
}

// Output ...
type Output struct {
	numOutput  int         //`json:"numoutput"`
	OutputSets []OutputSet `json:"records,omitempty"`
}

// OutputSet ...
type OutputSet struct {
	Columns Column `json:"column"`
}

// Column ...
type Column struct {
	Name     string          `json:"name"`
	DataType string          `json:"datatype"`
	SqlType  api.SQLSMALLINT `json:"sqltype"`
}
