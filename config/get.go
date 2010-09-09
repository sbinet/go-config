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
	"strconv"
	"strings"
)


/* Bool has the same behaviour as String but converts the response to bool.
See "boolString" for string values converted to bool.
*/
func (self *Config) Bool(section string, option string) (value bool, err os.Error) {
	sv, err := self.String(section, option)
	if err != nil {
		return false, err
	}

	value, ok := boolString[strings.ToLower(sv)]
	if !ok {
		return false, os.NewError(boolError(sv).String())
	}

	return value, nil
}

/* Float has the same behaviour as String but converts the response to float. */
func (self *Config) Float(section string, option string) (value float, err os.Error) {
	sv, err := self.String(section, option)
	if err == nil {
		value, err = strconv.Atof(sv)
	}

	return value, err
}

/* Int has the same behaviour as String but converts the response to int. */
func (self *Config) Int(section string, option string) (value int, err os.Error) {
	sv, err := self.String(section, option)
	if err == nil {
		value, err = strconv.Atoi(sv)
	}

	return value, err
}

/* RawString gets the (raw) string value for the given option in the section.
The raw string value is not subjected to unfolding, which was illustrated in the
beginning of this documentation.

It returns an error if either the section or the option do not exist.
*/
func (self *Config) RawString(section string, option string) (value string, err os.Error) {
	if _, ok := self.data[section]; ok {
		if value, ok = self.data[section][option]; ok {
			return value, nil
		}
		return "", os.NewError(optionError(option).String())
	}
	return "", os.NewError(sectionError(section).String())
}

/* String gets the string value for the given option in the section.
If the value needs to be unfolded (see e.g. %(host)s example in the beginning of
this documentation), then String does this unfolding automatically, up to
_DEPTH_VALUES number of iterations.

It returns an error if either the section or the option do not exist, or the
unfolding cycled.
*/
func (self *Config) String(section string, option string) (value string, err os.Error) {
	value, err = self.RawString(section, option)
	if err != nil {
		return "", err
	}

	var i int

	for i = 0; i < _DEPTH_VALUES; i++ { // keep a sane depth
		vr := varRegExp.FindString(value)
		if len(vr) == 0 {
			break
		}

		// Take off leading '%(' and trailing ')s'
		noption := strings.TrimLeft(vr, "%(")
		noption = strings.TrimRight(noption, ")s")

		// Search variable in default section
		nvalue, _ := self.data[_DEFAULT_SECTION][noption]
		if _, ok := self.data[section][noption]; ok {
			nvalue = self.data[section][noption]
		}
		if nvalue == "" {
			return "", os.NewError(optionError(noption).String())
		}

		// substitute by new value and take off leading '%(' and trailing ')s'
		value = strings.Replace(value, vr, nvalue, -1)
	}

	if i == _DEPTH_VALUES {
		return "", os.NewError(maxDephError(strconv.Itoa(_DEPTH_VALUES)).String())
	}

	return value, nil
}

