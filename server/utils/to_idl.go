package listener

import (
	"regexp"
)

type IdlTarget struct {
}

// func (i*IdlTarget)ExitJson is called when production json is exited.
func (i *IdlTarget) ExitJson(typeStr, valStr string) string {
	/*if typeStr=="struct"{
	      return "{"+valStr+"}"
		}else if typeStr=="array"{
			return "["+valStr+"]"
		}
		return typeStr+":"+valStr
	*/
	if stripNewLine(valStr)[0] == '[' {
		return "struct GeneratedStruct " + stripNewLine(stripArr(valStr)) + "\n" + "struct GeneratedArr {\n 1:  list<GeneratedStruct> generatedStruct \n}"
	}
	return "struct GeneratedStruct " + stripNewLine(valStr)
}

// func (i*IdlTarget)ExitObj is called when production obj is exited.
func (i *IdlTarget) PreExitObj(typeStr, valStr string) string {
	//return "{"+valStr+"}"
	return "{\n"
}
func (i *IdlTarget) ExitObj(typeStr, valStr string, isEnd bool, bIsMap bool) string {
	//return "{"+valStr+"}"
	//fmt.Println(typeStr,valStr,isEnd)
	if isEnd {
		return valStr + "\n"
	}
	return valStr + "\n"
}
func (i *IdlTarget) PostExitObj(typeStr, valStr string) string {
	//return "{"+valStr+"}"
	return "}"
}

// func (i*IdlTarget)ExitPair is called when production pair is exited.
func (i *IdlTarget) ExitPair(index int, keyStr, typeStr, valStr, valType string) (string, string) {
	if typeStr == "struct" {
		if IsNumber(keyStr) {
			return typeStr + " " + " " + camel(stripQuotes(keyStr)) + " " + stripNewLine(valStr) + "\n", camel(stripQuotes(keyStr)) + " " + stripQuotes(keyStr)
			//fmt.Println("number key",valStr)
			//return  "","map<string,string>"
		}
		return typeStr + " " + " " + camel(stripQuotes(keyStr)) + " " + stripNewLine(valStr) + "\n", camel(stripQuotes(keyStr)) + " " + stripQuotes(keyStr) //camel(stripQuotes(keyStr))+":"+keyStr
	}

	if typeStr == "array" {
		if valType == "struct" {
			return "struct" + " " + camel(stripQuotes(keyStr)) + stripNewLine(stripArr(valStr)) + "\n", "list<" + camel(stripQuotes(keyStr)) + "> " + stripQuotes(keyStr)
		}
		if valType == "" {
			return "struct" + " " + camel(stripQuotes(keyStr)) + "{\n}\n", "list<" + camel(stripQuotes(keyStr)) + "> " + stripQuotes(keyStr)
		}
		if valType == "array" {
			return "", "list< " + i.typeMap(valType, valStr) + "> " + stripQuotes(keyStr) + " //eg:" + removeEndline(valStr)
		}
		return "", "list<" + i.typeMap(valType, valStr) + "> " + stripQuotes(keyStr) + " //eg:" + removeEndline(valStr)
	}

	if typeStr == "null" {
		return "struct" + " " + camel(stripQuotes(keyStr)) + "{\n}\n", camel(stripQuotes(keyStr)) + " " + stripQuotes(keyStr) + " //eg:" + valStr //keyStr+":"+valStr
	}
	return "", i.typeMap(valType, valStr) + " " + stripQuotes(keyStr) + " //eg:" + valStr //keyStr+":"+valStr
}

// func (i*IdlTarget)ExitArr is called when production arr is exited.
func (i *IdlTarget) ExitArr(typeStr, valStr string) string {
	return "[\n" + valStr + "\n]"
}

// func (i*IdlTarget)ExitValue is called when production value is exited.
func (i *IdlTarget) ExitValue(typeStr, valStr string) string {
	return valStr
}

func (i *IdlTarget) typeMap(valType, val string) string {
	if valType == "array" {
		return "list<" + valType + ">"
	}

	if valType == "true" || valType == "false" {
		return "bool"
	}

	if valType == "float64" {
		f, _ := regexp.Compile("^-?([1-9]\\d*\\.\\d*|0\\.\\d*[1-9]\\d*|0?\\.0+)$")
		fe, _ := regexp.Compile("^-?((\\d+\\.?\\d*)|(\\.\\d+))[Ee][+-]?\\d+$")
		if f.Match([]byte(val)) || fe.Match([]byte(val)) {
			return "double"
		}
		if len(val) > 6 {
			return "i64"
		}
		return "i32"
	}

	if valType == "struct" && val == "null" {
		return valType
	}
	return valType
}
