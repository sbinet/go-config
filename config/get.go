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


/* GetBool has the same behaviour as GetString but converts the response to bool.
See constant BoolStrings for string values converted to bool.
*/
func (c *ConfigFile) GetBool(section string, option string) (value bool, err os.Error) {
	sv, err := c.GetString(section, option)
	if err != nil {
		return false, err
	}

	value, ok := BoolStrings[strings.ToLower(sv)]
	if !ok {
		return false, os.NewError("could not parse bool value: " + sv)
	}

	return value, nil
}

/* GetFloat has the same behaviour as GetString but converts the response to
float.
*/
func (c *ConfigFile) GetFloat(section string, option string) (value float, err os.Error) {
	sv, err := c.GetString(section, option)
	if err == nil {
		value, err = strconv.Atof(sv)
	}

	return value, err
}

/* GetInt has the same behaviour as GetString but converts the response to int. */
func (c *ConfigFile) GetInt(section string, option string) (value int, err os.Error) {
	sv, err := c.GetString(section, option)
	if err == nil {
		value, err = strconv.Atoi(sv)
	}

	return value, err
}

/* GetRawString gets the (raw) string value for the given option in the section.
The raw string value is not subjected to unfolding, which was illustrated in the
beginning of this documentation.

It returns an error if either the section or the option do not exist.
*/
func (c *ConfigFile) GetRawString(section string, option string) (value string, err os.Error) {
	section = strings.ToLower(section)
	option = strings.ToLower(option)

	if _, ok := c.data[section]; ok {
		if value, ok = c.data[section][option]; ok {
			return value, nil
		}
		return "", os.NewError("option not found")
	}
	return "", os.NewError("section not found")
}

/* GetString gets the string value for the given option in the section.
If the value needs to be unfolded (see e.g. %(host)s example in the beginning of
this documentation), then GetString does this unfolding automatically, up to
DepthValues number of iterations.

It returns an error if either the section or the option do not exist, or the
unfolding cycled.
*/
func (c *ConfigFile) GetString(section string, option string) (value string, err os.Error) {
	value, err = c.GetRawString(section, option)
	if err != nil {
		return "", err
	}

	section = strings.ToLower(section)

	var i int

	for i = 0; i < DepthValues; i++ { // keep a sane depth
		vr := varRegExp.FindString(value)
		if len(vr) == 0 {
			break
		}

		noption := value[vr[2]:vr[3]]
		noption = strings.ToLower(noption)

		nvalue, _ := c.data[DefaultSection][noption] // search variable in default section
		if _, ok := c.data[section][noption]; ok {
			nvalue = c.data[section][noption]
		}
		if nvalue == "" {
			return "", os.NewError("option not found: " + noption)
		}

		// substitute by new value and take off leading '%(' and trailing ')s'
		value = value[0:vr[2]-2] + nvalue + value[vr[3]+2:]
	}

	if i == DepthValues {
		return "", os.NewError("possible cycle while unfolding variables: max depth of " + strconv.Itoa(DepthValues) + " reached")
	}

	return value, nil
}

