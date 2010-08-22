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
See constant BoolStrings for string values converted to bool.
*/
func (self *File) Bool(section string, option string) (value bool, err os.Error) {
	sv, err := self.String(section, option)
	if err != nil {
		return false, err
	}

	value, ok := BoolStrings[strings.ToLower(sv)]
	if !ok {
		return false, os.NewError("could not parse bool value: " + sv)
	}

	return value, nil
}

/* Float has the same behaviour as String but converts the response to float. */
func (self *File) Float(section string, option string) (value float, err os.Error) {
	sv, err := self.String(section, option)
	if err == nil {
		value, err = strconv.Atof(sv)
	}

	return value, err
}

/* Int has the same behaviour as String but converts the response to int. */
func (self *File) Int(section string, option string) (value int, err os.Error) {
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
func (self *File) RawString(section string, option string) (value string, err os.Error) {
	section = strings.ToLower(section)
	option = strings.ToLower(option)

	if _, ok := self.data[section]; ok {
		if value, ok = self.data[section][option]; ok {
			return value, nil
		}
		return "", os.NewError("option not found")
	}
	return "", os.NewError("section not found")
}

/* String gets the string value for the given option in the section.
If the value needs to be unfolded (see e.g. %(host)s example in the beginning of
this documentation), then String does this unfolding automatically, up to
_DEPTH_VALUES number of iterations.

It returns an error if either the section or the option do not exist, or the
unfolding cycled.
*/
func (self *File) String(section string, option string) (value string, err os.Error) {
	value, err = self.RawString(section, option)
	if err != nil {
		return "", err
	}

	section = strings.ToLower(section)

	var i int

	for i = 0; i < _DEPTH_VALUES; i++ { // keep a sane depth
		vr := varRegExp.FindString(value)
		if len(vr) == 0 {
			break
		}

		// Take off leading '%(' and trailing ')s'
		noption := strings.TrimLeft(vr, "%(")
		noption = strings.TrimRight(noption, ")s")
		noption = strings.ToLower(noption)

		nvalue, _ := self.data[DEFAULT_SECTION][noption] // search variable in default section
		if _, ok := self.data[section][noption]; ok {
			nvalue = self.data[section][noption]
		}
		if nvalue == "" {
			return "", os.NewError("option not found: " + noption)
		}

		// substitute by new value and take off leading '%(' and trailing ')s'
		value = strings.Replace(value, vr, nvalue, -1)
	}

	if i == _DEPTH_VALUES {
		return "", os.NewError("possible cycle while unfolding variables: max depth of " + strconv.Itoa(_DEPTH_VALUES) + " reached")
	}

	return value, nil
}

