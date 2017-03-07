// Copyright 2015-2017, Cyrill @ Schumacher.fm and the CoreStore contributors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package _codegen

import (
	"bytes"
	"fmt"
	"go/format"
	"math/rand"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"github.com/corestoreio/csfw/util"
	"github.com/corestoreio/errors"
	"github.com/corestoreio/log"
)

var (
	letters   = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	Copyright = []byte(`// Copyright 2015-2017, Cyrill @ Schumacher.fm and the CoreStore contributors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
`)
)

// OFile is a OutputFile
type OFile string

func (of OFile) String() string {
	return string(of) + ".go"
}

func (of OFile) AppendDir(s ...string) OFile {
	sof := string(of)
	parts := strings.Split(sof, string(filepath.Separator))
	parts = append(parts, s...)
	nf := NewOFile(parts...)
	if filepath.IsAbs(sof) {
		nf = OFile(string(filepath.Separator)) + nf
	}
	return nf
}

func (of OFile) AppendName(s ...string) OFile {
	return NewOFile(string(of) + strings.Join(s, ""))
}

// NewOFile creates a new path from parts
func NewOFile(paths ...string) OFile {
	return OFile(filepath.Join(paths...))
}

// GenerateCode uses text/template for create Go code. package name pkg will also be used
// to remove stutter in variable names.
func GenerateCode(pkg, tplCode string, data interface{}, addFM template.FuncMap) ([]byte, error) {

	funcMap := template.FuncMap{
		"quote":           func(s string) string { return "`" + s + "`" },
		"prepareVar":      prepareVar(pkg),
		"camelize":        util.UnderscoreCamelize,
		"toLowerFirst":    toLowerFirst,
		"prepareVarIndex": func(i int, s string) string { return fmt.Sprintf("%03d%s", i, prepareVar(pkg)(s)) },
		"sprintf":         fmt.Sprintf,
	}
	for k, v := range addFM {
		funcMap[k] = v
	}

	codeTpl := template.Must(template.New("tpl_code").Funcs(funcMap).Parse(tplCode))

	var buf = new(bytes.Buffer)
	err := codeTpl.Execute(buf, data)
	if err != nil {
		return nil, errors.Wrap(err, "template execute")
	}

	fmt, err := format.Source(buf.Bytes())
	if err != nil {
		return buf.Bytes(), errors.Wrap(err, "gofmt")
	}
	return fmt, nil
}

func toLowerFirst(s string) string {
	if len(s) == 0 {
		return ""
	}
	return strings.ToLower(s[:1]) + s[1:]
}

// prepareVar converts a string into a Go code variable. Removes the package name if this string
// starts with the package name. Replaces all illegal characters with an underscore.
func prepareVar(pkg string) func(s string) string {

	return func(str string) string {

		l := len(pkg) + 1
		if len(str) > l && str[:l] == pkg+TableNameSeparator {
			str = str[l:]
		}

		return util.UnderscoreCamelize(str)
	}
}

// PrepareVar converts a string into a Go code variable. Removes the package name if this string
// starts with the package name. Replaces all illegal characters with an underscore.
func PrepareVar(pkg, s string) string {
	return prepareVar(pkg)(s)
}

// LogFatal logs an error as fatal with printed location and terminates the program.
func LogFatal(err error) {
	if err == nil {
		return
	}
	PkgLog.Fatal("ERROR", log.String("The_Error", fmt.Sprintf("%+v", err)), log.Err(err))
}

// randSeq returns a random string with a defined length n.
func randSeq(n int) string {
	rand.Seed(time.Now().UTC().UnixNano())
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

// ReplaceTablePrefix replaces the {{tableprefix}} place holder with the configure real TablePrefix
// TablePrefix can be set via init() statement in config_user.go
func ReplaceTablePrefix(query string) string {
	return strings.Replace(query, "{{tableprefix}}", TablePrefix, -1)
}

func extractSplit(s string) (path string, tp []string, hasVersion bool, _ error) {
	hasVersion = false

	pp := strings.Split(s, "/")
	if len(pp) == 1 {
		return path, tp, hasVersion, errors.Errorf("Path %q contains no slashes", s)
	}

	switch pp[0] { // add more domains here
	case "gopkg.in":
		hasVersion = true
		break
	}
	path = strings.Join(pp[:len(pp)-1], "/")
	tp = strings.Split(pp[len(pp)-1:][0], ".")

	if len(tp) == 1 {
		return path, tp, hasVersion, errors.Errorf("Missing . in package.Type")
	}
	return path, tp, hasVersion, nil
}

// ExtractImportPath extracts from an extended import path with a function or type call
// the import path.
// github.com/corestoreio/csfw/customer.Customer() would become
// github.com/corestoreio/csfw/customer
func ExtractImportPath(s string) (string, error) {

	if s == "" {
		return "", nil
	}
	path, tp, hasVersion, err := extractSplit(s)
	if err != nil {
		return "", err
	}

	path = path + "/" + tp[0]
	if hasVersion {
		path = path + "." + strings.Join(tp[1:len(tp)-1], ".")
	}

	return path, nil
}

// ExtractFuncType extracts from an extended import path with a function or type call
// the function or type call.
// github.com/corestoreio/csfw/customer.Customer() would become customer.Customer()
func ExtractFuncType(s string) (string, error) {
	if s == "" {
		return "", nil
	}

	_, tp, hasVersion, err := extractSplit(s)
	if err != nil {
		return "", errors.Wrap(err, "extractSplit")
	}

	if hasVersion {
		return strings.Join(append(tp[:1], tp[len(tp)-1:]...), "."), nil // cut the version out
	}
	return strings.Join(tp, "."), nil
}

// ParseString text/template for a string. Fails on error.
func ParseString(tpl string, data interface{}) string {
	codeTpl := template.Must(template.New("tpl_code").Parse(tpl))
	var buf = &bytes.Buffer{}
	err := codeTpl.Execute(buf, data)
	if err != nil {
		LogFatal(err)
	}
	return buf.String()
}
