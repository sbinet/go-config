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
	"bufio"
	"os"
	"strings"
)


/* Base to read a file and get the configuration representation.
That representation can be queried with GetString, etc.
*/
func _read(fname string, c *Config) (*Config, os.Error) {
	var file *os.File
	var err os.Error

	if file, err = os.Open(fname, os.O_RDONLY, 0); err != nil {
		return nil, err
	}

	if err = c.read(bufio.NewReader(file)); err != nil {
		return nil, err
	}

	if err = file.Close(); err != nil {
		return nil, err
	}

	return c, nil
}

/* ReadDefault reads a configuration file and returns its representation.
All arguments, except `fname`, are related to `New()`
*/
func Read(fname string, comment, separator string, preSpace, postSpace bool) (
*Config, os.Error) {
	return _read(fname, New(comment, separator, preSpace, postSpace))
}

/* ReadDefault reads a configuration file and returns its representation.
It uses values by default.
*/
func ReadDefault(fname string) (*Config, os.Error) {
	return _read(fname, NewDefault())
}

// ===

func (self *Config) read(buf *bufio.Reader) (err os.Error) {
	var section, option string
	for {
		l, err := buf.ReadString('\n') // parse line-by-line
		if err == os.EOF {
			break
		} else if err != nil {
			return err
		}

		l = strings.TrimSpace(l)
		// switch written for readability (not performance)
		switch {
		case len(l) == 0: // empty line
			continue

		case l[0] == '#': // comment
			continue

		case l[0] == ';': // comment
			continue

		// comment (for windows users)
		case len(l) >= 3 && strings.ToLower(l[0:3]) == "rem":
			continue

		case l[0] == '[' && l[len(l)-1] == ']': // new section
			option = "" // reset multi-line value
			section = strings.TrimSpace(l[1 : len(l)-1])
			self.AddSection(section)

		case section == "": // not new section and no section defined so far
			return os.NewError("section not found: must start with section")

		default: // other alternatives
			i := strings.IndexAny(l, "=:")
			switch {
			case i > 0: // option and value
				i := strings.IndexAny(l, "=:")
				option = strings.TrimSpace(l[0:i])
				value := strings.TrimSpace(stripComments(l[i+1:]))
				self.AddOption(section, option, value)

			case section != "" && option != "": // continuation of multi-line value
				prev, _ := self.RawString(section, option)
				value := strings.TrimSpace(stripComments(l))
				self.AddOption(section, option, prev+"\n"+value)

			default:
				return os.NewError("could not parse line: " + l)
			}
		}
	}
	return nil
}

