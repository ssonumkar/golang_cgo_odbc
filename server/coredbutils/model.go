package coredbutils

import (
	"os/exec"
)

type DBDiagnostic struct {
	State       string `json:"State"`
	NativeError string `json:"NativeError"`
	Messge      string `json:"Message"`
}

type ActionError struct {
	APIName string         `json:"APIName"`
	Diag    []DBDiagnostic `json:"Diag,omitempty"`
}

// CPH ...
type CPH struct {
	InstanceName string
	CredFile     []byte
}

// UniqueProxy ...
type UniqueProxy struct {
	Cph      *CPH
	Port     string
	Cmd      *exec.Cmd
	FileName string
}

var proxyPool = make(map[string]UniqueProxy)
