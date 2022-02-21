package coredbutils

import (
	"encoding/json"
)

func New(actionerr ActionError) error {
	return &ActionError{
		APIName: actionerr.APIName,
		Diag:    actionerr.Diag,
	}
}

func (acterr *ActionError) Error() string {
	out, err := json.Marshal(acterr)
	if err != nil {
		return "Error JSON Marshal Fialed"
	}
	return string(out)
}
