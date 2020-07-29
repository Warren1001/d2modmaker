package main

import (
	"log"
	"net/http"

	"github.com/shurcooL/vfsgen"
)

var reactFS http.FileSystem = http.Dir("../../gui/bin/react-ui/build")
var reactOptions vfsgen.Options = vfsgen.Options{Filename: "generated.go", PackageName: "generated", VariableName: "ReactAssets"}

func main() {
	reactErr := vfsgen.Generate(reactFS, reactOptions)
	if reactErr != nil {
		log.Fatalln(reactErr)
	}

}
