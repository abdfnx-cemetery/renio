package main

import (
	"fmt"
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"

	"github.com/abdfnx/shell"

	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/js"
)

func main() {
	files, _ := ioutil.ReadDir("js")

	sort.Sort(ByNumericalFilename(files))

	m := minify.New()
	m.AddFunc("text/javascript", js.Minify)

	var finalSource string

	for _, f := range files {
		log.Printf("Bundling %s\n", f.Name())

		file, err := os.Open(filepath.Join("js", f.Name()))
		if err != nil {
			log.Fatalf("Got error opening %s: %v", f.Name(), err)
		}

		buf := new(bytes.Buffer)

		if err := m.Minify("text/javascript", buf, file); err != nil {
			log.Fatalf("Got error minifying %s: %v", f.Name(), err)
		}

		finalSource += buf.String() + "\n"
	}

	err := os.Mkdir("target", 0750)
	if err != nil {
		log.Fatalf("Error in making directory - %v", err)
	}

	err = ioutil.WriteFile(filepath.Join("target", "renop.js"), []byte(finalSource), 0644)
	if err != nil {
		log.Fatalf("Error writing file %v", err)
	}

	err, cmd, errout := shell.RunOut("go run github.com/go-bindata/go-bindata/go-bindata -pkg core -o ./core/data.go typescript/ target/")

	log.Printf("Running command and waiting for it to finish...")

	if err != nil {
		log.Printf("error: %v\n", err)
		fmt.Print(errout)
	}

	fmt.Print(cmd)

	log.Printf("Command finished with error: %v", err)
}
