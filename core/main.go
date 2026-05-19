package main

/*
#include <stdlib.h>
*/
import "C"
import (
	"fmt"
	"unsafe"

	"com.argotunnel/core/capabilities"
	"com.argotunnel/core/engine"
)

var globalExecutor *engine.Executor

//export StartEngine
func StartEngine(dbPath *C.char, encryptionKey *C.char) *C.char {
	goDbPath := C.GoString(dbPath)
	goKey := C.GoString(encryptionKey)

	dbInstance, err := capabilities.NewEncryptedDB(goDbPath, goKey)
	if err != nil {
		return C.CString(fmt.Sprintf("ERROR: Failed to initialize DB: %v", err))
	}

	selector := capabilities.NewSelector(dbInstance)
	globalExecutor = engine.NewExecutor(selector, dbInstance)

	err = globalExecutor.Start()
	if err != nil {
		return C.CString(fmt.Sprintf("ERROR: Failed to start engine: %v", err))
	}

	return C.CString("SUCCESS: Engine started autonomously")
}

//export StopEngine
func StopEngine() *C.char {
	if globalExecutor == nil {
		return C.CString("ERROR: Engine not running")
	}

	err := globalExecutor.Stop()
	if err != nil {
		return C.CString(fmt.Sprintf("ERROR: Failed to stop engine: %v", err))
	}

	globalExecutor = nil
	return C.CString("SUCCESS: Engine stopped and memory scrubbed")
}

func main() {
	// Required for c-shared build mode
}
