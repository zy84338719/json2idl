package listener

type Target interface {
	// ExitJson is called when production json is exited.
	ExitJson(typeStr, valStr string) string
	// ExitObj is called when production obj is exited.
	PreExitObj(typeStr, valStr string) string
	ExitObj(typeStr, valStr string, isEnd bool, bIsMap bool) string
	PostExitObj(typeStr, valStr string) string
	// ExitPair is called when production pair is exited.
	ExitPair(index int, ketStr, typeStr, valStr, valType string) (string, string)
	// ExitArr is called when production arr is exited.
	ExitArr(typeStr, valStr string) string
	// ExitValue is called when production value is exited.
	ExitValue(typeStr, valStr string) string
}
