package dotenv

import (
	"bytes"
	"fmt"
	"strings"
)

func Unmarshal(in []byte, out *map[string]interface{}) error {
	if *out == nil {
		*out = make(map[string]interface{})
	}
	for _, line := range bytes.Split(in, []byte("\n")) {
		if len(line) == 0 {
			continue
		}
		if line[0] == '#' {
			continue
		} else {
			pos := bytes.Index(line, []byte("="))
			if pos == -1 {
				return fmt.Errorf("invalid dotenv input line: %s", line)
			}
			(*out)[string(line[:pos])] = strings.Replace(string(line[pos+1:]), "\\n", "\n", -1)
		}
	}

	return nil
}
