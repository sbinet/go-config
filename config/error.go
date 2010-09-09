// Copyright 2009  The "goconfig" Authors
//
// Use of this source code is governed by the Simplified BSD License
// that can be found in the LICENSE file.
//
// This software is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES
// OR CONDITIONS OF ANY KIND, either express or implied. See the License
// for more details.

package config


type sectionError string

func (self sectionError) String() string {
	return "goconfig: section not found: " + string(self)
}

// ===

type optionError string

func (self optionError) String() string {
	return "goconfig: option not found: " + string(self)
}

// ===

type boolError string

func (self boolError) String() string {
	return "goconfig: could not parse bool value: " + string(self)
}

// ===

type maxDephError string

func (self maxDephError) String() string {
	return "goconfig: possible cycle while unfolding variables: max depth of " +
		string(self) + " reached"
}

// ===

type lineError string

func (self lineError) String() string {
	return "goconfig: could not parse line: " + string(self)
}

