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
	"os"
	"strings"
)


/* AddOption adds a new option and value to the configuration.
It returns true if the option and value were inserted, and false if the value
was overwritten. If the section does not exist in advance, it is created.
*/
func (c *ConfigFile) AddOption(section string, option string, value string) bool {
	c.AddSection(section) // make sure section exists

	section = strings.ToLower(section)
	option = strings.ToLower(option)

	_, ok := c.data[section][option]
	c.data[section][option] = value

	return !ok
}

/* RemoveOption removes a option and value from the configuration.
It returns true if the option and value were removed, and false otherwise,
including if the section did not exist.
*/
func (c *ConfigFile) RemoveOption(section string, option string) bool {
	section = strings.ToLower(section)
	option = strings.ToLower(option)

	if _, ok := c.data[section]; !ok {
		return false
	}

	_, ok := c.data[section][option]
	c.data[section][option] = "", false

	return ok
}

/* HasOption checks if the configuration has the given option in the section.
It returns false if either the option or section do not exist.
*/
func (c *ConfigFile) HasOption(section string, option string) bool {
	section = strings.ToLower(section)
	option = strings.ToLower(option)

	if _, ok := c.data[section]; !ok {
		return false
	}

	_, okd := c.data[DefaultSection][option]
	_, oknd := c.data[section][option]

	return okd || oknd
}

/* GetOptions returns the list of options available in the given section.
It returns an error if the section does not exist and an empty list if the
section is empty. Options within the default section are also included.
*/
func (c *ConfigFile) GetOptions(section string) (options []string, err os.Error) {
	section = strings.ToLower(section)

	if _, ok := c.data[section]; !ok {
		return nil, os.NewError("section not found")
	}

	options = make([]string, len(c.data[DefaultSection])+len(c.data[section]))
	i := 0
	for s, _ := range c.data[DefaultSection] {
		options[i] = s
		i++
	}
	for s, _ := range c.data[section] {
		options[i] = s
		i++
	}

	return options, nil
}

