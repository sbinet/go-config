// Copyright 2010  The "goconfig" Authors
//
// Use of this source code is governed by the Simplified BSD License
// that can be found in the LICENSE file.
//
// This software is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES
// OR CONDITIONS OF ANY KIND, either express or implied. See the License
// for more details.

package config

import (
	"errors"
	"strconv"
	"strings"
)

// Bool has the same behaviour as String but converts the response to bool.
// See "boolString" for string values converted to bool.
func (self *Config) Bool(section string, option string) (value bool, err error) {
	sv, err := self.String(section, option)
	if err != nil {
		return false, err
	}

	value, ok := boolString[strings.ToLower(sv)]
	if !ok {
		return false, errors.New("could not parse bool value: " + sv)
	}

	return value, nil
}

// Float has the same behaviour as String but converts the response to float.
func (self *Config) Float(section string, option string) (value float64, err error) {
	sv, err := self.String(section, option)
	if err == nil {
		value, err = strconv.ParseFloat(sv, 64)
	}

	return value, err
}

// Int has the same behaviour as String but converts the response to int.
func (self *Config) Int(section string, option string) (value int, err error) {
	sv, err := self.String(section, option)
	if err == nil {
		value, err = strconv.Atoi(sv)
	}

	return value, err
}

// RawString gets the (raw) string value for the given option in the section.
// The raw string value is not subjected to unfolding, which was illustrated in
// the beginning of this documentation.
//
// It returns an error if either the section or the option do not exist.
func (self *Config) RawString(section string, option string) (value string, err error) {
	if _, ok := self.data[section]; ok {
		if tValue, ok := self.data[section][option]; ok {
			return tValue.v, nil
		}
		return "", errors.New(optionError(option).String())
	}
	return "", errors.New(sectionError(section).String())
}

// String gets the string value for the given option in the section.
// If the value needs to be unfolded (see e.g. %(host)s example in the beginning
// of this documentation), then String does this unfolding automatically, up to
// _DEPTH_VALUES number of iterations.
//
// It returns an error if either the section or the option do not exist, or the
// unfolding cycled.
func (self *Config) String(section string, option string) (value string, err error) {
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
		if nvalue.v == "" {
			return "", errors.New(optionError(noption).String())
		}

		// substitute by new value and take off leading '%(' and trailing ')s'
		value = strings.Replace(value, vr, nvalue.v, -1)
	}

	if i == _DEPTH_VALUES {
		return "", errors.New("possible cycle while unfolding variables: " +
			"max depth of " + strconv.Itoa(_DEPTH_VALUES) + " reached")
	}

	return value, nil
}
