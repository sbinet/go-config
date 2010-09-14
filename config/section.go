// Copyright 2010  The "goconfig" Authors
//
// Use of this source code is governed by the Simplified BSD License
// that can be found in the LICENSE file.
//
// This software is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES
// OR CONDITIONS OF ANY KIND, either express or implied. See the License
// for more details.

package config


/* AddSection adds a new section to the configuration.

If the section is nil then uses the section by default which it's already
created.

It returns true if the new section was inserted, and false if the section
already existed.
*/
func (self *Config) AddSection(section string) bool {
	// _DEFAULT_SECTION
	if section == "" {
		return false
	}

	if _, ok := self.data[section]; ok {
		return false
	}

	self.data[section] = make(map[string]*tValue)

	return true
}

/* RemoveSection removes a section from the configuration.
It returns true if the section was removed, and false if section did not exist.
*/
func (self *Config) RemoveSection(section string) bool {
	_, ok := self.data[section]

	// Default section cannot be removed.
	if !ok || section == _DEFAULT_SECTION {
		return false
	}

	for o, _ := range self.data[section] {
		self.data[section][o] = nil, false // *value
	}

	self.data[section] = nil, false
	self.idOption[section] = 0, false

/*
	switch _, ok := self.data[section]; {
	case !ok:
		return false
	case section == _DEFAULT_SECTION:
		return false // default section cannot be removed
	default:
		for o, _ := range self.data[section] {
			self.data[section][o] = nil, false // *value
		}
		self.data[section] = nil, false
	}
*/
	return true
}

/* HasSection checks if the configuration has the given section.
(The default section always exists.)
*/
func (self *Config) HasSection(section string) bool {
	_, ok := self.data[section]

	return ok
}

/* Sections returns the list of sections in the configuration.
(The default section always exists.)
*/
func (self *Config) Sections() (sections []string) {
	sections = make([]string, len(self.data))

	i := 0
	for s, _ := range self.data {
		sections[i] = s
		i++
	}

	return sections
}

