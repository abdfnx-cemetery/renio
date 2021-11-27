package tools

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func FileX(file string) string {
	tsjs := ""
  
	// check if the file is javascript or typescript
	if strings.HasSuffix(file, ".js") {
	  tsjs = "temp.js"
	} else if strings.HasSuffix(file, ".ts") {
	  tsjs = "temp.ts"
	}
  
	tempDir := filepath.Join(os.TempDir(), tsjs)

	return tempDir
}

func Download(url string) {
  // Create the file
  out, err := os.Create(FileX(url))

  if err != nil  {
    fmt.Println(err)
  }

  defer out.Close()

  // Get the data
  resp, err := http.Get(url)
  if err != nil {
    fmt.Println(err)
  }

  defer resp.Body.Close()

  // Check server response
  if resp.StatusCode != http.StatusOK {
    fmt.Errorf("bad status: %s", resp.Status)
  }

  // Writer the body to file
  _, err = io.Copy(out, resp.Body)
  if err != nil  {
    fmt.Println(err)
  }
}
