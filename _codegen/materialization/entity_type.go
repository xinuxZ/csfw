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

package main

import (
	"fmt"
	"io/ioutil"
	"text/template"

	"github.com/corestoreio/csfw/_codegen"
	"github.com/corestoreio/csfw/eav"
	"github.com/corestoreio/csfw/storage/dbr"
	"github.com/corestoreio/csfw/util"
	"github.com/juju/errgo"
)

// materializeEntityType writes the data from eav_entity_type into a Go file and transforms
// Magento classes and config strings into Go functions.
// Depends on generated code from tableToStruct.
func materializeEntityType(ctx *context) {
	defer ctx.wg.Done()
	type dataContainer struct {
		ETypeData     eav.TableEntityTypeSlice
		ImportPaths   []string
		Package, Tick string
	}

	etData, err := getEntityTypeData(ctx.dbc.NewSession())
	_codegen.LogFatal(err)

	tplData := &dataContainer{
		ETypeData:   etData,
		ImportPaths: getImportPaths(),
		Package:     _codegen.ConfigMaterializationEntityType.Package,
		Tick:        "`",
	}

	addFM := template.FuncMap{
		"extractFuncType": _codegen.ExtractFuncType,
	}

	formatted, err := _codegen.GenerateCode(_codegen.ConfigMaterializationEntityType.Package, tplEav, tplData, addFM)
	if err != nil {
		fmt.Printf("\n%s\n", formatted)
		_codegen.LogFatal(err)
	}

	_codegen.LogFatal(ioutil.WriteFile(_codegen.ConfigMaterializationEntityType.OutputFile, formatted, 0600))
}

// getEntityTypeData retrieves all EAV models from table eav_entity_type but only those listed in variable
// _codegen.ConfigEntityType. It then applies the mapping data from _codegen.ConfigEntityType to the entity_type struct.
// Depends on generated code from tableToStruct.
func getEntityTypeData(dbrSess *dbr.Session) (etc eav.TableEntityTypeSlice, err error) {

	s, err := eav.TableCollection.Structure(eav.TableIndexEntityType)
	if err != nil {
		return nil, errgo.Mask(err)
	}

	_, err = dbrSess.
		Select(s.AllColumnAliasQuote(s.Name)...).
		From(s.Name).
		Where("entity_type_code IN ?", _codegen.ConfigEntityType.Keys()).
		LoadStructs(&etc)
	if err != nil {
		return nil, errgo.Mask(err)
	}

	for typeCode, mapData := range _codegen.ConfigEntityType {
		// map the fields from the config struct to the data retrieved from the database.
		et, err := etc.GetByCode(typeCode)
		_codegen.LogFatal(err)
		et.EntityModel = _codegen.ParseString(mapData.EntityModel, et)
		et.AttributeModel.String = _codegen.ParseString(mapData.AttributeModel, et)
		et.EntityTable.String = _codegen.ParseString(mapData.EntityTable, et)
		et.IncrementModel.String = _codegen.ParseString(mapData.IncrementModel, et)
		et.AdditionalAttributeTable.String = _codegen.ParseString(mapData.AdditionalAttributeTable, et)
		et.EntityAttributeCollection.String = _codegen.ParseString(mapData.EntityAttributeCollection, et)
	}

	return etc, nil
}

func getImportPaths() []string {
	var paths util.StringSlice

	var getPath = func(s string) string {
		ps, err := _codegen.ExtractImportPath(s)
		_codegen.LogFatal(err)
		return ps
	}

	for _, et := range _codegen.ConfigEntityType {
		paths.Append(
			getPath(et.EntityModel),
			getPath(et.AttributeModel),
			getPath(et.EntityTable),
			getPath(et.IncrementModel),
			getPath(et.AdditionalAttributeTable),
			getPath(et.EntityAttributeCollection),
		)
	}
	return paths.Unique().ToString()
}
