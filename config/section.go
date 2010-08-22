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
	"strings"
)


/* AddSection adds a new section to the configuration.
It returns true if the new section was inserted, and false if the section
already existed.
*/
func (self *File) AddSection(section string) bool {
	section = strings.ToLower(section)

	if _, ok := self.data[section]; ok {
		return false
	}
	self.data[section] = make(map[string]string)

	return true
}

/* RemoveSection removes a section from the configuration.
It returns true if the section was removed, and false if section did not exist.
*/
func (self *File) RemoveSection(section string) bool {
	section = strings.ToLower(section)

	switch _, ok := self.data[section]; {
	case !ok:
		return false
	case section == DEFAULT_SECTION:
		return false // default section cannot be removed
	default:
		for o, _ := range self.data[section] {
			self.data[section][o] = "", false
		}
		self.data[section] = nil, false
	}

	return true
}

/* HasSection checks if the configuration has the given section.
(The default section always exists.)
*/
func (self *File) HasSection(section string) bool {
	_, ok := self.data[strings.ToLower(section)]

	return ok
}

/* Sections returns the list of sections in the configuration.
(The default section always exists.)
*/
func (self *File) Sections() (sections []string) {
	sections = make([]string, len(self.data))

	i := 0
	for s, _ := range self.data {
		sections[i] = s
		i++
	}

	return sections
}

