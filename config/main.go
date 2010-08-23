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


const (
	_DEFAULT_SECTION   = "default" // Default section name (must be lower-case).
	_DEFAULT_COMMENT   = "# "      // Another valid option is ";"
	_DEFAULT_SEPARATOR = ": "      // Another valid option is "="
	_DEPTH_VALUES      = 200       // Maximum allowed depth when recursively substituing variable names.
)

var (
	// Strings accepted as boolean.
	boolString = map[string]bool{
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

	varRegExp = regexp.MustCompile(`%\(([a-zA-Z0-9_.\-]+)\)s`) // %(variable)s
)


/* File is the representation of configuration settings. */
type File struct {
	data map[string]map[string]string // Maps sections to options to values.
}

/* NewFile creates an empty configuration representation.
This representation can be filled with AddSection and AddOption and then saved
to a file using WriteFile.
*/
func NewFile() *File {
	c := new(File)
	c.data = make(map[string]map[string]string)

	c.AddSection(_DEFAULT_SECTION) // default section always exists

	return c
}


// === Utility
// ===

func stripComments(l string) string {
	// comments are preceded by space or TAB
	for _, c := range []string{" ;", "\t;", " #", "\t#"} {
		if i := strings.Index(l, c); i != -1 {
			l = l[0:i]
		}
	}
	return l
}

