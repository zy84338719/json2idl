package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"github.com/antlr/antlr4/runtime/Go/antlr"
	listener "gitlab.com/jsonToIDL/utils"
	"gitlab.com/jsonToIDL/utils/parser"
	"io/fs"
	"net/http"
	"path"
)

func main() {
	HttpWebServer()
	HttpApiServer()
}

//go:embed web
var addr embed.FS

func HttpWebServer() {
	http.Handle("/", AssetHandler("", addr, "web"))
}

func HttpApiServer() {
	http.HandleFunc("/api/json2IDL", apiJson2IDL)
	err := http.ListenAndServe(":59990", nil)
	if err != nil {
		fmt.Println("exec failed, ", err)
	}
}

func apiJson2IDL(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		res, _ := json.Marshal(resquest{Code: 500})
		_, _ = w.Write(res)
		return
	}
	var jsonData map[string]string
	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&jsonData); err != nil {
		res, _ := json.Marshal(resquest{Code: 500})
		_, _ = w.Write(res)
		return
	}

	lex := parser.NewJSONLexer(antlr.NewInputStream(jsonData["JSON"]))
	stream := antlr.NewCommonTokenStream(lex, antlr.TokenDefaultChannel)
	par := parser.NewJSONParser(stream)
	l := listener.NewJsonToGoListener(&listener.IdlTarget{})
	antlr.ParseTreeWalkerDefault.Walk(l, par.Json())

	c := "\n"
	for _, s := range l.SubStructs {
		c += s
	}
	res, err := json.Marshal(resquest{Code: 0, Data: Data{IDL: c + l.JsonStr}})
	if err != nil {
		res, _ := json.Marshal(resquest{Code: 500})
		_, _ = w.Write(res)
		return
	}
	_, _ = w.Write(res)
}

type resquest struct {
	Code int  `json:"code"`
	Data Data `json:"data"`
}

type Data struct {
	IDL string `json:"idl"`
}

type fsFunc func(name string) (fs.File, error)

func (f fsFunc) Open(name string) (fs.File, error) {
	return f(name)
}

// AssetHandler returns an http.Handler that will serve files from
// the Assets embed.FS. When locating a file, it will strip the given
// prefix from the request and prepend the root to the filesystem.
func AssetHandler(prefix string, assets embed.FS, root string) http.Handler {
	handler := fsFunc(func(name string) (fs.File, error) {
		assetPath := path.Join(root, name)

		// If we can't find the asset, fs can handle the error
		file, err := assets.Open(assetPath)
		if err != nil {
			return nil, err
		}

		// Otherwise assume this is a legitimate request routed correctly
		return file, err
	})

	return http.StripPrefix(prefix, http.FileServer(http.FS(handler)))
}
