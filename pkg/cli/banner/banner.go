package banner

import (
	"fmt"
	"io"
)

const asciiArt = `
             __           __           
 _   _______/ /_  _______/ /____  _____
| | / / ___/ / / / / ___/ __/ _ \/ ___/
| |/ / /__/ / /_/ (__  ) /_/  __/ /    
|___/\___/_/\__,_/____/\__/\___/_/     
`

func Print(w io.Writer) {
	fmt.Fprint(w, asciiArt)
}

func PrintWithVersion(w io.Writer, version string) {
	fmt.Fprint(w, asciiArt)
	fmt.Fprintf(w, "                                %s\n", version)
}
