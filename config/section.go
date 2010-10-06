// Copyright 2010  The "goconfig" Authors
//
// Use of this source code is governed by the Simplified BSD License
// that can be found in the LICENSE file.
//
// This software is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES
// OR CONDITIONS OF ANY KIND, either express or implied. See the License
// for more details.

package config


// AddSection adds a new section to the configuration.
//
// If the section is nil then uses the section by default which it's already
// created.
//
// It returns true if the new section was inserted, and false if the section
// already existed.
func (self *Config) AddSection(section string) bool {
	// _DEFAULT_SECTION
	if section == "" {
		return false
	}

	if _, ok := self.data[section]; ok {
		return false
	}

	self.data[section] = make(map[string]*tValue)

	// Section order
	self.idSection[section] = self.lastIdSection
	self.lastIdSection++

	return true
}

// RemoveSection removes a section from the configuration.
// It returns true if the section was removed, and false if section did not exist.
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

	self.lastIdOption[section] = 0, false
	self.idSection[section] = 0, false

	return true
}

// HasSection checks if the configuration has the given section.
// (The default section always exists.)
func (self *Config) HasSection(section string) bool {
	_, ok := self.data[section]

	return ok
}

// Sections returns the list of sections in the configuration.
// (The default section always exists.)
func (self *Config) Sections() (sections []string) {
	sections = make([]string, len(self.idSection))
	pos := 0 // Position in sections

	for i := 0; i < self.lastIdSection; i++ {
		for section, id := range self.idSection {
			if id == i {
				sections[pos] = section
				pos++
			}
		}
	}

	return sections
}

