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
	"bufio"
	"fmt"
	"os"
	"strings"
)


/* WriteFile saves the configuration representation to a file.
The desired file permissions must be passed as in os.Open. The header is a
string that is saved as a comment in the first line of the file.
*/
func (self *File) WriteFile(fname string, perm uint32, header string) (err os.Error) {
	var file *os.File

	if file, err = os.Open(fname, os.O_WRONLY|os.O_CREAT|os.O_TRUNC, perm); err != nil {
		return err
	}

	buf := bufio.NewWriter(file)
	if err = self.write(buf, header); err != nil {
		return err
	}
	buf.Flush()

	return file.Close()
}

func (self *File) write(buf *bufio.Writer, header string) (err os.Error) {
	if header != "" {
		if i := strings.Index(header, "\n"); i != -1 {
			header = strings.Replace(header, "\n", "\n"+_DEFAULT_COMMENT, -1)
		}

		if _, err = buf.WriteString(fmt.Sprint(
			_DEFAULT_COMMENT, header, "\n")); err != nil {
			return err
		}
	}

	for section, sectionmap := range self.data {
		if section == _DEFAULT_SECTION && len(sectionmap) == 0 {
			continue // Skips default section if empty.
		}
		if _, err = buf.WriteString(fmt.Sprint("\n[", section, "]\n")); err != nil {
			return err
		}
		for option, value := range sectionmap {
			if _, err = buf.WriteString(fmt.Sprint(
				option, _DEFAULT_SEPARATOR, value, "\n")); err != nil {
				return err
			}
		}
	}
	if _, err = buf.WriteString("\n"); err != nil {
		return err
	}

	return nil
}

