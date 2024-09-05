package logger

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func setupLogger() (*bytes.Buffer, *bytes.Buffer) {
	infoBuffer := new(bytes.Buffer)
	errorBuffer := new(bytes.Buffer)

	InfoLogger.SetOutput(infoBuffer)
	WarnLogger.SetOutput(infoBuffer)
	SuccessLogger.SetOutput(infoBuffer)
	ErrorLogger.SetOutput(errorBuffer)

	return infoBuffer, errorBuffer
}

func TestLogging(t *testing.T) {
	Init()
	infoBuffer, errorBuffer := setupLogger()

	Info("This is an info message")
	require.Contains(t, infoBuffer.String(), "This is an info message")

	Warn("This is a warning message")
	require.Contains(t, infoBuffer.String(), "This is a warning message")

	Success("This is a success message")
	require.Contains(t, infoBuffer.String(), "This is a success message")

	Error("This is an error message")
	require.Contains(t, errorBuffer.String(), "This is an error message")
}

func TestLoggingWithStruct(t *testing.T) {
	Init()
	infoBuffer, _ := setupLogger()

	type TestStruct struct {
		Field1 string
		Field2 int
	}
	Info(TestStruct{Field1: "value", Field2: 123})

	require.Contains(t, infoBuffer.String(), `{"Field1":"value","Field2":123}`)
}
