// Copyright 2009  The "goconfig" Authors
//
// Use of this source code is governed by the Simplified BSD License
// that can be found in the LICENSE file.
//
// This software is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES
// OR CONDITIONS OF ANY KIND, either express or implied. See the License
// for more details.

package config

import (
	"regexp"
	"strings"
)


var (
	DefaultSection = "default" // Default section name (must be lower-case).
	DepthValues    = 200       // Maximum allowed depth when recursively substituing variable names.

	// Strings accepted as bool.
	BoolStrings = map[string]bool{
		"t":     true,
		"true":  true,
		"y":     true,
		"yes":   true,
		"on":    true,
		"1":     true,
		"f":     false,
		"false": false,
		"n":     false,
		"no":    false,
		"off":   false,
		"0":     false,
	}

	varRegExp = regexp.MustCompile(`%\(([a-zA-Z0-9_.\-]+)\)s`)
)


/* ConfigFile is the representation of configuration settings.
The public interface is entirely through methods.
*/
type ConfigFile struct {
	data map[string]map[string]string // Maps sections to options to values.
}

/* NewConfigFile creates an empty configuration representation.
This representation can be filled with AddSection and AddOption and then saved
to a file using WriteConfigFile.
*/
func NewConfigFile() *ConfigFile {
	c := new(ConfigFile)
	c.data = make(map[string]map[string]string)

	c.AddSection(DefaultSection) // default section always exists

	return c
}


// === Utility
// ===

func firstIndex(s string, delim []byte) int {
	for i := 0; i < len(s); i++ {
		for j := 0; j < len(delim); j++ {
			if s[i] == delim[j] {
				return i
			}
		}
	}
	return -1
}

func stripComments(l string) string {
	// comments are preceded by space or TAB
	for _, c := range []string{" ;", "\t;", " #", "\t#"} {
		if i := strings.Index(l, c); i != -1 {
			l = l[0:i]
		}
	}
	return l
}

