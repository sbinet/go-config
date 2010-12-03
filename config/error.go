// Copyright 2010  The "goconfig" Authors
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
	return "section not found: " + string(self)
}


type optionError string

func (self optionError) String() string {
	return "option not found: " + string(self)
}

