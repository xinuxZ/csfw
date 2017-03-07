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
	"go/build"

	"github.com/corestoreio/csfw/_codegen/tableToStruct/tpl"
	"github.com/corestoreio/csfw/util/slices"
)

type (
	// TableToStructMap uses a string key as easy identifier for maybe later
	// manipulation in config_user.go and a pointer to a TableToStruct struct.
	// Developers can later use the init() func in config_user.go to change
	// the value of variable ConfigTableToStruct.
	TableToStructMap map[string]TableToStruct
	TableToStruct    struct {
		// Package defines the name of the target package
		Package string
		// OutputFile specifies the full path where to write the newly generated code.
		// The file extension .go will be added automatically.
		OutputFile OFile
		// QueryString SQL query to filter all the tables which you desire, e.g.
		// catalog\_% OR select table from information_schema ...
		// This query must specify all tables you need for a package.
		SQLQuery string
		// EntityTypeCodes If provided then eav_entity_type.value_table_prefix
		// will be evaluated for further tables.
		EntityTypeCodes []string
		// GenericsWhiteList config option as a SQL query to select all tables
		// for which it should generate the generic functions. If empty nothing
		// gets generated. If GenericsWhiteList contains the word SQLQuery then
		// the query will be copied from SQLQuery field and all tables receive
		// the generated functions.
		GenericsWhiteList string
		// GenericsFunctions specify which functions you need in the whole
		// package
		GenericsFunctions tpl.Generics
	}

	// AttributeToStructMap contains as key the name of the EAV entity and points to
	// its configuration.
	AttributeToStructMap map[string]*AttributeToStruct
	// AttributeToStruct contains the configuration to materialize all
	// attributes belonging to one EAV model
	AttributeToStruct struct {
		// AttrPkgImp defines the package import path to use: possible ATM:
		// custattr and catattr and your custom EAV package.
		AttrPkgImp string
		// FuncCollection specifies the name of the SetAttributeCollection
		// function name within the AttrPkgImp.
		FuncCollection string
		// FuncGetter specifies the name of the SetAttributeGetter function name
		// within the AttrPkgImp.
		FuncGetter string
		// Package defines the name of the target package, must be external.
		// Default is {{.AttrPkgImp}}_test but for your project you must provide
		// your package name.
		Package string
		// AttrStruct defines the name of the attribute struct in an EAV package
		// like catattr or custattr. This struct will have the prefix of
		// AttrPkgImp and will be embedded the newly generated struct which
		// wraps additional project specific columns.
		AttrStruct string
		// MyStruct is the optional name of your struct from your package. @todo review
		MyStruct string
		// OutputFile specifies the full path where to write the newly generated
		// code
		OutputFile string
	}

	// EntityTypeMap uses a string key as for the EAV entity type, which must
	// also exists in the table eav_EntityTable, for maybe later manipulation in
	// config_user.go and as value a pointer to a EntityType struct. Developers
	// can later use the init() func in config_user.go to change the value of
	// variable ConfigEntityType. The values of EntityType will be uses for
	// materialization in Go code of the eav_entity_type table data.
	EntityTypeMap map[string]*EntityType

	// EntityType is configuration struct which maps the PHP classes to Go
	// types, interfaces and table names. Each struct field has a special import
	// path with a function to make it easier to specify different packages.
	EntityType struct {
		// EntityModel Go type which implements eav.EntityTypeModeller. Will be
		// used as template so you can access the current entity_type from the
		// database.
		EntityModel string
		// AttributeModel Go type which implements
		// eav.EntityTypeAttributeModeller Will be used as template so you can
		// access the current entity_type from the database.
		AttributeModel string
		// EntityTable Go type which implements eav.EntityTypeTabler Will be
		// used as template so you can access the current entity_type from the
		// database.
		EntityTable string
		// IncrementModel Go type which implements
		// eav.EntityTypeIncrementModeller Will be used as template so you can
		// access the current entity_type from the database.
		IncrementModel string
		// AdditionalAttributeTable Go type which implements
		// eav.EntityTypeAdditionalAttributeTabler Will be used as template so
		// you can access the current entity_type from the database.
		AdditionalAttributeTable string
		// EntityAttributeCollection Go type which implements
		// eav.EntityTypeAttributeCollectioner Will be used as template so you
		// can access the current entity_type from the database.
		EntityAttributeCollection string

		// TempAdditionalAttributeTable string which defines the existing table
		// name and specifies more attribute configuration options besides
		// eav_attribute table. This table name is used in a DB query while
		// materializing attribute configuration to Go code.
		// Mage_Eav_Model_Resource_Attribute_Collection::_initSelect()
		TempAdditionalAttributeTable string

		// TempAdditionalAttributeTableWebsite string which defines the existing
		// table name and stores website-dependent attribute parameters. If an
		// EAV model doesn't demand this functionality, let this string empty.
		// This table name is used in a DB query while materializing attribute
		// configuration to Go code.
		// Mage_Customer_Model_Resource_Attribute::_getEavWebsiteTable()
		// Mage_Eav_Model_Resource_Attribute_Collection::_getEavWebsiteTable()
		TempAdditionalAttributeTableWebsite string

		AttributeCoreColumns slices.String
	}

	// AttributeModelDefMap contains data to map the three eav_attribute columns
	// (backend | frontend | source | data)_model to the correct Go function and
	// package. It contains mappings for Magento 1 & 2. A developer has the
	// option to change/extend the value using the file config_user.go with the
	// init() func. Def for Definition to avoid a naming conflict :-( Better
	// name?
	AttributeModelDefMap map[string]*AttributeModelDef
	// AttributeModelDef defines which Go type/func has which import path
	AttributeModelDef struct {
		// GoFunc is a function string which implements, when later executed,
		// one of the n interfaces for (backend|frontend|source|data)_model. The
		// GoFunc expects the fully qualified import path to the final method,
		// e.g.: github.com/corestoreio/csfw/customer.Customer()
		GoFunc string
		// i cached import path
		i string
		// f cached func name
		f string
	}
)

// NewAMD NewAttributeModelDef NewAttributeModelDef creates a new attribute model definition
func NewAMD(GoFuncWithPath string) *AttributeModelDef {
	i, err := ExtractImportPath(GoFuncWithPath)
	LogFatal(err)
	f, err := ExtractFuncType(GoFuncWithPath)
	LogFatal(err)

	return &AttributeModelDef{
		GoFunc: GoFuncWithPath,
		i:      i,
		f:      f,
	}
}

// Import extracts the import path
func (d *AttributeModelDef) Import() string {
	return d.i
}

// Func extracts the function
func (d *AttributeModelDef) Func() string {
	return d.f
}

// Keys returns all keys from a EntityTypeMap
func (m EntityTypeMap) Keys() []string {
	ret := make([]string, len(m), len(m))
	i := 0
	for k := range m {
		ret[i] = k
		i++
	}
	return ret
}

// BasePath now a hack but should also consider if it has been vendored.
var BasePath = NewOFile(build.Default.GOPATH, "github.com/corestoreio/csfw")

// EavAttributeColumnNameToInterface mapping column name to Go interface name. Do not add attribute_model
// as this column is unused in Magento 1+2. If you have custom column then add it here.
var EavAttributeColumnNameToInterface = map[string]string{
	"backend_model":           "eav.AttributeBackendModeller",
	"frontend_model":          "eav.AttributeFrontendModeller",
	"source_model":            "eav.AttributeSourceModeller",
	"frontend_input_renderer": "eav.FrontendInputRendererIFace",
	"data_model":              "eav.AttributeDataModeller",
}

// TablePrefix defines the global table name prefix. See Magento install tool. Can be overridden via func init()
var TablePrefix string

// TableMapMagento1To2 provides mapping between table names in tableToStruct. If a table name is in
// the map then the struct name will be rewritten to that new Magneto2 compatible table name.
// Do not change entries in this map except you can always append.
// @see Magento2 dev/tools/Magento/Tools/Migration/factory_table_names/replace_ce.php
var TableMapMagento1To2 = map[string]string{
	"admin_role":             "authorization_role",
	"admin_rule":             "authorization_rule",
	"core_cache":             "cache",     // not needed but added in case
	"core_cache_tag":         "cache_tag", // not needed but added in case
	"core_config_data":       "core_config_data",
	"core_design_change":     "design_change", // not needed but added in case
	"core_directory_storage": "media_storage_directory_storage",
	"core_email_template":    "email_template",
	"core_file_storage":      "media_storage_file_storage",
	"core_flag":              "flag",          // not needed but added in case
	"core_layout_link":       "layout_link",   // not needed but added in case
	"core_layout_update":     "layout_update", // not needed but added in case
	"core_resource":          "setup_module",  // not needed but added in case
	"core_session":           "session",       // not needed but added in case
	"core_store":             "store",
	"core_store_group":       "store_group",
	"core_variable":          "variable",
	"core_variable_value":    "variable_value",
	"core_website":           "store_website",
}

// ConfigTableToStruct contains default configuration. Use the file config_user.go with the func init() to change/extend it.
var ConfigTableToStruct = TableToStructMap{
	"authorization": TableToStruct{
		Package:    "authorization",
		OutputFile: BasePath.AppendDir("authorization", "tables_generated"),
		SQLQuery: `SELECT TABLE_NAME FROM information_schema.COLUMNS WHERE
						TABLE_SCHEMA = DATABASE() AND
						TABLE_NAME IN (
							'{{tableprefix}}authorization_role','{{tableprefix}}admin_role',
							'{{tableprefix}}authorization_rule','{{tableprefix}}admin_rule'
						) GROUP BY TABLE_NAME;`,
		EntityTypeCodes:   nil,
		GenericsWhiteList: "SQLQuery",
		GenericsFunctions: tpl.OptSQL | tpl.OptFindBy,
	},
	"user": TableToStruct{
		Package:           "user",
		OutputFile:        BasePath.AppendDir("user", "tables_generated"),
		SQLQuery:          `{{tableprefix}}admin_user`,
		EntityTypeCodes:   nil,
		GenericsWhiteList: "SQLQuery",
		GenericsFunctions: tpl.OptAll,
	},
	"config/storage/ccd": TableToStruct{
		Package:           "ccd",
		OutputFile:        BasePath.AppendDir("config", "storage", "ccd", "tables_generated"),
		SQLQuery:          `{{tableprefix}}core_config_data`,
		EntityTypeCodes:   nil,
		GenericsWhiteList: "",
		// GenericsFunctions: tpl.OptAll, not needed ... or @todo
	},
	"directory": TableToStruct{
		Package:    "directory",
		OutputFile: BasePath.AppendDir("directory", "tables_generated"),
		SQLQuery: `SELECT TABLE_NAME FROM information_schema.COLUMNS WHERE
						TABLE_SCHEMA = DATABASE() AND
						TABLE_NAME LIKE '{{tableprefix}}directory%' GROUP BY TABLE_NAME;`,
		EntityTypeCodes:   nil,
		GenericsWhiteList: "SQLQuery",
		GenericsFunctions: tpl.OptAll,
	},
	"eav": TableToStruct{
		Package:         "eav",
		OutputFile:      BasePath.AppendDir("eav", "tables_generated"),
		SQLQuery:        `{{tableprefix}}eav%`,
		EntityTypeCodes: nil,
		GenericsWhiteList: `SELECT TABLE_NAME FROM information_schema.COLUMNS WHERE
						TABLE_SCHEMA = DATABASE() AND
						TABLE_NAME LIKE '{{tableprefix}}eav_entity_type' GROUP BY TABLE_NAME;`,
		GenericsFunctions: tpl.OptSQL | tpl.OptFindBy,
	},
	"store": TableToStruct{
		Package:    "store",
		OutputFile: BasePath.AppendDir("store", "tables_generated"),
		SQLQuery: `SELECT TABLE_NAME FROM information_schema.COLUMNS WHERE
					TABLE_SCHEMA = DATABASE() AND
					TABLE_NAME IN (
						'{{tableprefix}}core_store','{{tableprefix}}store',
						'{{tableprefix}}core_store_group','{{tableprefix}}store_group',
						'{{tableprefix}}core_website','{{tableprefix}}store_website'
					) GROUP BY TABLE_NAME;`,
		EntityTypeCodes:   nil,
		GenericsWhiteList: "SQLQuery",
		GenericsFunctions: tpl.OptAll,
	},
	"catalog": TableToStruct{
		// @todo figure out tables which are in both Magneto version present
		Package:    "catalog",
		OutputFile: BasePath.AppendDir("catalog", "tables_generated"),
		SQLQuery: `SELECT TABLE_NAME FROM information_schema.COLUMNS WHERE
						TABLE_SCHEMA = DATABASE() AND
						(TABLE_NAME LIKE '{{tableprefix}}catalog\_%' OR TABLE_NAME LIKE '{{tableprefix}}catalogindex%' ) AND
						TABLE_NAME NOT LIKE '{{tableprefix}}%bundle%' AND
						TABLE_NAME NOT LIKE '{{tableprefix}}%\_flat\_%' GROUP BY TABLE_NAME;`,
		EntityTypeCodes: []string{"catalog_category", "catalog_product"},
		GenericsWhiteList: `SELECT TABLE_NAME FROM information_schema.COLUMNS WHERE
								TABLE_SCHEMA = DATABASE() AND
								(TABLE_NAME LIKE '{{tableprefix}}catalog\_%' OR TABLE_NAME LIKE '{{tableprefix}}catalogindex%' ) AND
								TABLE_NAME NOT LIKE '{{tableprefix}}%bundle%' AND
								TABLE_NAME NOT LIKE '{{tableprefix}}%idx%' AND
								TABLE_NAME NOT LIKE '{{tableprefix}}%tmp%' AND
								TABLE_NAME NOT LIKE '{{tableprefix}}%\_flat\_%' GROUP BY TABLE_NAME;`,
		GenericsFunctions: tpl.OptAll,
	},
	"customer": TableToStruct{
		Package:           "customer",
		OutputFile:        BasePath.AppendDir("customer", "tables_generated"),
		SQLQuery:          `{{tableprefix}}customer%`,
		EntityTypeCodes:   []string{"customer", "customer_address"},
		GenericsWhiteList: "SQLQuery",
		GenericsFunctions: tpl.OptAll,
	},
}

// ConfigMaterializationEntityType configuration for materializeEntityType() to
// write the materialized entity types into a folder. Other fields of the struct
// TableToStruct are ignored. Use the file config_user.go with the func init()
// to change/extend it.
var ConfigMaterializationEntityType = TableToStruct{
	Package:    "testgen",
	OutputFile: BasePath.AppendDir("testgen", "generated_entity_type_test"),
}

// ConfigLocalization temporary integration because missing feature in https://github.com/golang/text
var ConfigLocalization = struct {
	Package       string
	OutputFile    string
	EnabledLocale slices.String
}{
	Package:       "i18n_test",
	OutputFile:    BasePath.AppendDir("i18n", "generated_translation_test").String(),
	EnabledLocale: slices.String{"en", "fr", "de", "de_CH", "nl", "ca_FR"},
}

var (
	// EAVAttributeCoreColumns defines the minimal required columns for table
	// eav_attribute. Developers can extend the table eav_attribute with
	// additional columns but these additional columns with its method receivers
	// must get generated in the attribute materialize function. These core
	// columns are already defined below.
	EAVAttributeCoreColumns = slices.String{
		"attribute_id",
		"entity_type_id",
		"attribute_code",
		"attribute_model", // this column is unused by Mage1+2
		"backend_model",
		"backend_type",
		"backend_table",
		"frontend_model",
		"frontend_input",
		"frontend_label",
		"frontend_class",
		"source_model",
		"is_required",
		"is_user_defined",
		"default_value",
		"is_unique",
		"note",
	}

	// customerAttributeCoreColumns defines the minimal required columns for
	// table customer_eav_attribute. Developers can extend the table
	// customer_eav_attribute with additional columns but these additional
	// columns with its method receivers must get generated in the attribute
	// materialize function. These core columns are already defined below.
	customerAttributeCoreColumns = slices.String{
		"is_visible",
		"input_filter",
		"multiline_count",
		"validate_rules",
		"is_system",
		"sort_order",
		"data_model",
	}

	// catalogAttributeCoreColumns defines the minimal required columns for
	// table catalog_eav_attribute. Developers can extend the table
	// customer_eav_attribute with additional columns but these additional
	// columns with its method receivers must get generated in the attribute
	// materialize function. These core columns are already defined below.
	catalogAttributeCoreColumns = slices.String{
		"frontend_input_renderer",
		"is_global",
		"is_visible",
		"is_searchable",
		"is_filterable",
		"is_comparable",
		"is_visible_on_front",
		"is_html_allowed_on_front",
		"is_used_for_price_rules",
		"is_filterable_in_search",
		"used_in_product_listing",
		"used_for_sort_by",
		"is_configurable",
		"apply_to",
		"is_visible_in_advanced_search",
		"position",
		"is_wysiwyg_enabled",
		"is_used_for_promo_rules",
		"search_weight",
	}
)

// ConfigEntityType contains default configuration to materialize the entity
// types. Use the file config_user.go with the func init() to change/extend it.
// Needed in materializeEntityType()
var ConfigEntityType = EntityTypeMap{
	"customer": &EntityType{
		EntityModel:                         "github.com/corestoreio/csfw/customer.Customer()",
		AttributeModel:                      "github.com/corestoreio/csfw/customer/custattr.HandlerCustomer({{.EntityTypeID}})",
		EntityTable:                         "github.com/corestoreio/csfw/customer.Customer()",
		IncrementModel:                      "github.com/corestoreio/csfw/customer.Customer()",
		AdditionalAttributeTable:            "github.com/corestoreio/csfw/customer.Customer()",
		EntityAttributeCollection:           "github.com/corestoreio/csfw/customer/custattr.HandlerCustomer({{.EntityTypeID}})",
		TempAdditionalAttributeTable:        "{{tableprefix}}customer_eav_attribute",
		TempAdditionalAttributeTableWebsite: "{{tableprefix}}customer_eav_attribute_website",
		AttributeCoreColumns:                customerAttributeCoreColumns,
	},
	"customer_address": &EntityType{
		EntityModel:                         "github.com/corestoreio/csfw/customer.Address()",
		AttributeModel:                      "github.com/corestoreio/csfw/customer/custattr.HandlerAddress({{.EntityTypeID}})",
		EntityTable:                         "github.com/corestoreio/csfw/customer.Address()",
		AdditionalAttributeTable:            "github.com/corestoreio/csfw/customer.Address()",
		EntityAttributeCollection:           "github.com/corestoreio/csfw/customer/custattr.HandlerAddress({{.EntityTypeID}})",
		TempAdditionalAttributeTable:        "{{tableprefix}}customer_eav_attribute",
		TempAdditionalAttributeTableWebsite: "{{tableprefix}}customer_eav_attribute_website",
		AttributeCoreColumns:                customerAttributeCoreColumns,
	},
	"catalog_category": &EntityType{
		EntityModel:                  "github.com/corestoreio/csfw/catalog.Category()",
		AttributeModel:               "github.com/corestoreio/csfw/catalog/catattr.HandlerCategory({{.EntityTypeID}})",
		EntityTable:                  "github.com/corestoreio/csfw/catalog.Category()",
		AdditionalAttributeTable:     "github.com/corestoreio/csfw/catalog.Category()",
		EntityAttributeCollection:    "github.com/corestoreio/csfw/catalog/catattr.HandlerCategory({{.EntityTypeID}})",
		TempAdditionalAttributeTable: "{{tableprefix}}catalog_eav_attribute",
		AttributeCoreColumns:         catalogAttributeCoreColumns,
	},
	"catalog_product": &EntityType{
		EntityModel:                  "github.com/corestoreio/csfw/catalog.Product()",
		AttributeModel:               "github.com/corestoreio/csfw/catalog/catattr.HandlerProduct({{.EntityTypeID}})",
		EntityTable:                  "github.com/corestoreio/csfw/catalog.Product()",
		AdditionalAttributeTable:     "github.com/corestoreio/csfw/catalog.Product()",
		EntityAttributeCollection:    "github.com/corestoreio/csfw/catalog/catattr.HandlerProduct({{.EntityTypeID}})",
		TempAdditionalAttributeTable: "{{tableprefix}}catalog_eav_attribute",
		AttributeCoreColumns:         catalogAttributeCoreColumns,
	},
	// @todo extend for all sales entities
}

// ConfigMaterializationAttributes contains the configuration to materialize all
// attributes for the defined EAV entity types.
var ConfigMaterializationAttributes = AttributeToStructMap{
	"customer": &AttributeToStruct{
		AttrPkgImp:     "github.com/corestoreio/csfw/customer/custattr",
		FuncCollection: "SetCustomerCollection",
		FuncGetter:     "SetCustomerGetter",
		AttrStruct:     "Customer",
		MyStruct:       "",
		Package:        "testgen", // external package name
		OutputFile:     BasePath.AppendDir("testgen", "generated_customer_attribute_test").String(),
	},
	"customer_address": &AttributeToStruct{
		AttrPkgImp:     "github.com/corestoreio/csfw/customer/custattr",
		FuncCollection: "SetAddressCollection",
		FuncGetter:     "SetAddressGetter",
		AttrStruct:     "Customer",
		MyStruct:       "",
		Package:        "testgen",
		OutputFile:     BasePath.AppendDir("testgen", "generated_address_attribute_test").String(),
	},
	"catalog_product": &AttributeToStruct{
		AttrPkgImp:     "github.com/corestoreio/csfw/catalog/catattr",
		FuncCollection: "SetProductCollection",
		FuncGetter:     "SetProductGetter",
		AttrStruct:     "Catalog",
		MyStruct:       "",
		Package:        "testgen",
		OutputFile:     BasePath.AppendDir("testgen", "generated_product_attribute_test.go").String(),
	},
	"catalog_category": &AttributeToStruct{
		AttrPkgImp:     "github.com/corestoreio/csfw/catalog/catattr",
		FuncCollection: "SetCategoryCollection",
		FuncGetter:     "SetCategoryGetter",
		AttrStruct:     "Catalog",
		MyStruct:       "",
		Package:        "testgen",
		OutputFile:     BasePath.AppendDir("testgen", "generated_category_attribute_test.go").String(),
	},
	// extend here for other EAV attributes (not sales* types)
}

// ConfigAttributeModel contains default configuration. Use the file
// config_user.go with the func init() to change/extend it.
var ConfigAttributeModel = AttributeModelDefMap{
	"catalog/product_attribute_frontend_image":                                        NewAMD("github.com/corestoreio/csfw/catalog/catattr.ProductFrontendImage().Config(eav.AttributeFrontendIdx({{.AttributeIndex}}))"),
	"Magento\\Catalog\\Model\\Product\\Attribute\\Frontend\\Image":                    NewAMD("github.com/corestoreio/csfw/catalog/catattr.ProductFrontendImage().Config(eav.AttributeFrontendIdx({{.AttributeIndex}}))"),
	"eav/entity_attribute_frontend_datetime":                                          NewAMD("github.com/corestoreio/csfw/eav.AttributeFrontendDatetime().Config(eav.AttributeFrontendIdx({{.AttributeIndex}}))"),
	"Magento\\Eav\\Model\\Entity\\Attribute\\Frontend\\Datetime":                      NewAMD("github.com/corestoreio/csfw/eav.AttributeFrontendDatetime().Config(eav.AttributeFrontendIdx({{.AttributeIndex}}))"),
	"catalog/attribute_backend_customlayoutupdate":                                    NewAMD(""),
	"Magento\\Catalog\\Model\\Attribute\\Backend\\Customlayoutupdate":                 NewAMD(""),
	"catalog/category_attribute_backend_image":                                        NewAMD("github.com/corestoreio/csfw/catalog/catattr.CategoryBackendImage().Config(eav.AttributeBackendIdx({{.AttributeIndex}}))"),
	"Magento\\Catalog\\Model\\Category\\Attribute\\Backend\\Image":                    NewAMD("github.com/corestoreio/csfw/catalog/catattr.CategoryBackendImage().Config(eav.AttributeBackendIdx({{.AttributeIndex}}))"),
	"catalog/category_attribute_backend_sortby":                                       NewAMD("github.com/corestoreio/csfw/catalog/catattr.CategoryBackendSortby().Config(eav.AttributeBackendIdx({{.AttributeIndex}}))"),
	"Magento\\Catalog\\Model\\Category\\Attribute\\Backend\\Sortby":                   NewAMD("github.com/corestoreio/csfw/catalog/catattr.CategoryBackendSortby().Config(eav.AttributeBackendIdx({{.AttributeIndex}}))"),
	"catalog/category_attribute_backend_urlkey":                                       NewAMD(""),
	"catalog/product_attribute_backend_boolean":                                       NewAMD("github.com/corestoreio/csfw/catalog/catattr.ProductBackendBoolean().Config(eav.AttributeBackendIdx({{.AttributeIndex}}))"),
	"Magento\\Catalog\\Model\\Product\\Attribute\\Backend\\Boolean":                   NewAMD("github.com/corestoreio/csfw/catalog/catattr.ProductBackendBoolean().Config(eav.AttributeBackendIdx({{.AttributeIndex}}))"),
	"Magento\\Catalog\\Model\\Product\\Attribute\\Backend\\Category":                  NewAMD("github.com/corestoreio/csfw/catalog/catattr.ProductBackendCategory().Config(eav.AttributeBackendIdx({{.AttributeIndex}}))"),
	"catalog/product_attribute_backend_groupprice":                                    NewAMD("github.com/corestoreio/csfw/catalog/catattr.ProductBackendGroupPrice().Config(eav.AttributeBackendIdx({{.AttributeIndex}}))"),
	"Magento\\Catalog\\Model\\Product\\Attribute\\Backend\\GroupPrice":                NewAMD("github.com/corestoreio/csfw/catalog/catattr.ProductBackendGroupPrice().Config(eav.AttributeBackendIdx({{.AttributeIndex}}))"),
	"catalog/product_attribute_backend_media":                                         NewAMD("github.com/corestoreio/csfw/catalog/catattr.ProductBackendMedia().Config(eav.AttributeBackendIdx({{.AttributeIndex}}))"),
	"Magento\\Catalog\\Model\\Product\\Attribute\\Backend\\Media":                     NewAMD("github.com/corestoreio/csfw/catalog/catattr.ProductBackendMedia().Config(eav.AttributeBackendIdx({{.AttributeIndex}}))"),
	"catalog/product_attribute_backend_msrp":                                          NewAMD("github.com/corestoreio/csfw/catalog/catattr.ProductBackendPrice().Config(eav.AttributeBackendIdx({{.AttributeIndex}}))"),
	"catalog/product_attribute_backend_price":                                         NewAMD("github.com/corestoreio/csfw/catalog/catattr.ProductBackendPrice().Config(eav.AttributeBackendIdx({{.AttributeIndex}}))"),
	"Magento\\Catalog\\Model\\Product\\Attribute\\Backend\\Price":                     NewAMD("github.com/corestoreio/csfw/catalog/catattr.ProductBackendPrice().Config(eav.AttributeBackendIdx({{.AttributeIndex}}))"),
	"catalog/product_attribute_backend_recurring":                                     NewAMD("github.com/corestoreio/csfw/catalog/catattr.ProductBackendRecurring().Config(eav.AttributeBackendIdx({{.AttributeIndex}}))"),
	"Magento\\Catalog\\Model\\Product\\Attribute\\Backend\\Recurring":                 NewAMD("github.com/corestoreio/csfw/catalog/catattr.ProductBackendRecurring().Config(eav.AttributeBackendIdx({{.AttributeIndex}}))"),
	"catalog/product_attribute_backend_sku":                                           NewAMD("github.com/corestoreio/csfw/catalog/catattr.ProductBackendSku().Config(eav.AttributeBackendIdx({{.AttributeIndex}}))"),
	"Magento\\Catalog\\Model\\Product\\Attribute\\Backend\\Sku":                       NewAMD("github.com/corestoreio/csfw/catalog/catattr.ProductBackendSku().Config(eav.AttributeBackendIdx({{.AttributeIndex}}))"),
	"catalog/product_attribute_backend_startdate":                                     NewAMD("github.com/corestoreio/csfw/catalog/catattr.ProductBackendStartDate().Config(eav.AttributeBackendIdx({{.AttributeIndex}}))"),
	"Magento\\Catalog\\Model\\Product\\Attribute\\Backend\\Startdate":                 NewAMD("github.com/corestoreio/csfw/catalog/catattr.ProductBackendStartDate().Config(eav.AttributeBackendIdx({{.AttributeIndex}}))"),
	"Magento\\Catalog\\Model\\Product\\Attribute\\Backend\\Stock":                     NewAMD("github.com/corestoreio/csfw/catalog/catattr.ProductBackendStock().Config(eav.AttributeBackendIdx({{.AttributeIndex}}))"),
	"catalog/product_attribute_backend_tierprice":                                     NewAMD("github.com/corestoreio/csfw/catalog/catattr.ProductBackendTierPrice().Config(eav.AttributeBackendIdx({{.AttributeIndex}}))"),
	"Magento\\Catalog\\Model\\Product\\Attribute\\Backend\\Tierprice":                 NewAMD("github.com/corestoreio/csfw/catalog/catattr.ProductBackendTierPrice().Config(eav.AttributeBackendIdx({{.AttributeIndex}}))"),
	"Magento\\Catalog\\Model\\Product\\Attribute\\Backend\\Weight":                    NewAMD("github.com/corestoreio/csfw/catalog/catattr.ProductBackendWeight().Config(eav.AttributeBackendIdx({{.AttributeIndex}}))"),
	"catalog/product_attribute_backend_urlkey":                                        NewAMD(""),
	"customer/attribute_backend_data_boolean":                                         NewAMD("github.com/corestoreio/csfw/customer/custattr.CustomerBackendDataBoolean().Config(eav.AttributeBackendIdx({{.AttributeIndex}}))"),
	"Magento\\Customer\\Model\\Attribute\\Backend\\Data\\Boolean":                     NewAMD("github.com/corestoreio/csfw/customer/custattr.CustomerBackendDataBoolean().Config(eav.AttributeBackendIdx({{.AttributeIndex}}))"),
	"customer/customer_attribute_backend_billing":                                     NewAMD("github.com/corestoreio/csfw/customer/custattr.CustomerBackendBilling().Config(eav.AttributeBackendIdx({{.AttributeIndex}}))"),
	"Magento\\Customer\\Model\\Customer\\Attribute\\Backend\\Billing":                 NewAMD("github.com/corestoreio/csfw/customer/custattr.CustomerBackendBilling().Config(eav.AttributeBackendIdx({{.AttributeIndex}}))"),
	"customer/customer_attribute_backend_password":                                    NewAMD("github.com/corestoreio/csfw/customer/custattr.CustomerBackendPassword().Config(eav.AttributeBackendIdx({{.AttributeIndex}}))"),
	"Magento\\Customer\\Model\\Customer\\Attribute\\Backend\\Password":                NewAMD("github.com/corestoreio/csfw/customer/custattr.CustomerBackendPassword().Config(eav.AttributeBackendIdx({{.AttributeIndex}}))"),
	"customer/customer_attribute_backend_shipping":                                    NewAMD("github.com/corestoreio/csfw/customer/custattr.CustomerBackendShipping().Config(eav.AttributeBackendIdx({{.AttributeIndex}}))"),
	"Magento\\Customer\\Model\\Customer\\Attribute\\Backend\\Shipping":                NewAMD("github.com/corestoreio/csfw/customer/custattr.CustomerBackendShipping().Config(eav.AttributeBackendIdx({{.AttributeIndex}}))"),
	"customer/customer_attribute_backend_store":                                       NewAMD("github.com/corestoreio/csfw/customer/custattr.CustomerBackendStore().Config(eav.AttributeBackendIdx({{.AttributeIndex}}))"),
	"Magento\\Customer\\Model\\Customer\\Attribute\\Backend\\Store":                   NewAMD("github.com/corestoreio/csfw/customer/custattr.CustomerBackendStore().Config(eav.AttributeBackendIdx({{.AttributeIndex}}))"),
	"customer/customer_attribute_backend_website":                                     NewAMD("github.com/corestoreio/csfw/customer/custattr.CustomerBackendWebsite().Config(eav.AttributeBackendIdx({{.AttributeIndex}}))"),
	"Magento\\Customer\\Model\\Customer\\Attribute\\Backend\\Website":                 NewAMD("github.com/corestoreio/csfw/customer/custattr.CustomerBackendWebsite().Config(eav.AttributeBackendIdx({{.AttributeIndex}}))"),
	"customer/entity_address_attribute_backend_region":                                NewAMD("github.com/corestoreio/csfw/customer/custattr.AddressBackendRegion().Config(eav.AttributeBackendIdx({{.AttributeIndex}}))"),
	"Magento\\Customer\\Model\\Resource\\Address\\Attribute\\Backend\\Region":         NewAMD("github.com/corestoreio/csfw/customer/custattr.AddressBackendRegion().Config(eav.AttributeBackendIdx({{.AttributeIndex}}))"),
	"customer/entity_address_attribute_backend_street":                                NewAMD("github.com/corestoreio/csfw/customer/custattr.AddressBackendStreet().Config(eav.AttributeBackendIdx({{.AttributeIndex}}))"),
	"Magento\\Eav\\Model\\Entity\\Attribute\\Backend\\DefaultBackend":                 NewAMD("github.com/corestoreio/csfw/eav.AttributeBackendDefaultBackend().Config(eav.AttributeBackendIdx({{.AttributeIndex}}))"),
	"eav/entity_attribute_backend_datetime":                                           NewAMD("github.com/corestoreio/csfw/eav.AttributeBackendDatetime().Config(eav.AttributeBackendIdx({{.AttributeIndex}}))"),
	"Magento\\Eav\\Model\\Entity\\Attribute\\Backend\\Datetime":                       NewAMD("github.com/corestoreio/csfw/eav.AttributeBackendDatetime().Config(eav.AttributeBackendIdx({{.AttributeIndex}}))"),
	"eav/entity_attribute_backend_time_created":                                       NewAMD("github.com/corestoreio/csfw/eav.AttributeBackendTimeCreated().Config(eav.AttributeBackendIdx({{.AttributeIndex}}))"),
	"Magento\\Eav\\Model\\Entity\\Attribute\\Backend\\Time\\Created":                  NewAMD("github.com/corestoreio/csfw/eav.AttributeBackendTimeCreated().Config(eav.AttributeBackendIdx({{.AttributeIndex}}))"),
	"eav/entity_attribute_backend_time_updated":                                       NewAMD("github.com/corestoreio/csfw/eav.AttributeBackendTimeUpdated().Config(eav.AttributeBackendIdx({{.AttributeIndex}}))"),
	"Magento\\Eav\\Model\\Entity\\Attribute\\Backend\\Time\\Updated":                  NewAMD("github.com/corestoreio/csfw/eav.AttributeBackendTimeUpdated().Config(eav.AttributeBackendIdx({{.AttributeIndex}}))"),
	"bundle/product_attribute_source_price_view":                                      NewAMD("github.com/corestoreio/csfw/bundle.AttributeSourcePriceView().Config(eav.AttributeSourceIdx({{.AttributeIndex}}))"),
	"Magento\\Bundle\\Model\\Product\\Attribute\\Source\\Price\\View":                 NewAMD("github.com/corestoreio/csfw/bundle.AttributeSourcePriceView().Config(eav.AttributeSourceIdx({{.AttributeIndex}}))"),
	"Magento\\CatalogInventory\\Model\\Source\\Stock":                                 NewAMD("github.com/corestoreio/csfw/cataloginventory.AttributeSourceStock().Config(eav.AttributeSourceIdx({{.AttributeIndex}}))"),
	"catalog/category_attribute_source_layout":                                        NewAMD(""),
	"Magento\\Catalog\\Model\\Category\\Attribute\\Source\\Layout":                    NewAMD(""),
	"catalog/category_attribute_source_mode":                                          NewAMD("github.com/corestoreio/csfw/catalog/catattr.CategorySourceMode().Config(eav.AttributeSourceIdx({{.AttributeIndex}}))"),
	"Magento\\Catalog\\Model\\Category\\Attribute\\Source\\Mode":                      NewAMD("github.com/corestoreio/csfw/catalog/catattr.CategorySourceMode().Config(eav.AttributeSourceIdx({{.AttributeIndex}}))"),
	"catalog/category_attribute_source_page":                                          NewAMD("github.com/corestoreio/csfw/catalog/catattr.CategorySourcePage().Config(eav.AttributeSourceIdx({{.AttributeIndex}}))"),
	"Magento\\Catalog\\Model\\Category\\Attribute\\Source\\Page":                      NewAMD("github.com/corestoreio/csfw/catalog/catattr.CategorySourcePage().Config(eav.AttributeSourceIdx({{.AttributeIndex}}))"),
	"catalog/category_attribute_source_sortby":                                        NewAMD("github.com/corestoreio/csfw/catalog/catattr.CategorySourceSortby().Config(eav.AttributeSourceIdx({{.AttributeIndex}}))"),
	"Magento\\Catalog\\Model\\Category\\Attribute\\Source\\Sortby":                    NewAMD("github.com/corestoreio/csfw/catalog/catattr.CategorySourceSortby().Config(eav.AttributeSourceIdx({{.AttributeIndex}}))"),
	"catalog/entity_product_attribute_design_options_container":                       NewAMD("github.com/corestoreio/csfw/catalog/catattr.ProductSourceDesignOptionsContainer().Config(eav.AttributeSourceIdx({{.AttributeIndex}}))"),
	"Magento\\Catalog\\Model\\Entity\\Product\\Attribute\\Design\\Options\\Container": NewAMD("github.com/corestoreio/csfw/catalog/catattr.ProductSourceDesignOptionsContainer().Config(eav.AttributeSourceIdx({{.AttributeIndex}}))"),
	"catalog/product_attribute_source_countryofmanufacture":                           NewAMD("github.com/corestoreio/csfw/catalog/catattr.ProductSourceCountryOfManufacture().Config(eav.AttributeSourceIdx({{.AttributeIndex}}))"),
	"Magento\\Catalog\\Model\\Product\\Attribute\\Source\\Countryofmanufacture":       NewAMD("github.com/corestoreio/csfw/catalog/catattr.ProductSourceCountryOfManufacture().Config(eav.AttributeSourceIdx({{.AttributeIndex}}))"),
	"catalog/product_attribute_source_layout":                                         NewAMD("github.com/corestoreio/csfw/catalog/catattr.ProductSourceLayout().Config(eav.AttributeSourceIdx({{.AttributeIndex}}))"),
	"Magento\\Catalog\\Model\\Product\\Attribute\\Source\\Layout":                     NewAMD("github.com/corestoreio/csfw/catalog/catattr.ProductSourceLayout().Config(eav.AttributeSourceIdx({{.AttributeIndex}}))"),
	"Magento\\Catalog\\Model\\Product\\Attribute\\Source\\Status":                     NewAMD("github.com/corestoreio/csfw/catalog/catattr.ProductSourceStatus().Config(eav.AttributeSourceIdx({{.AttributeIndex}}))"),
	"catalog/product_attribute_source_msrp_type_enabled":                              NewAMD(""),
	"catalog/product_attribute_source_msrp_type_price":                                NewAMD("github.com/corestoreio/csfw/msrp.NewAttributeSourcePrice().Config(eav.AttributeSourceIdx({{.AttributeIndex}}))"),
	"Magento\\Msrp\\Model\\Product\\Attribute\\Source\\Type\\Price":                   NewAMD("github.com/corestoreio/csfw/msrp.NewAttributeSourcePrice().Config(eav.AttributeSourceIdx({{.AttributeIndex}}))"),
	"catalog/product_status":                                                          NewAMD("github.com/corestoreio/csfw/catalog/catattr.ProductSourceStatus().Config(eav.AttributeSourceIdx({{.AttributeIndex}}))"),
	"catalog/product_visibility":                                                      NewAMD("github.com/corestoreio/csfw/catalog/catattr.ProductSourceVisibility().Config(eav.AttributeSourceIdx({{.AttributeIndex}}))"),
	"Magento\\Catalog\\Model\\Product\\Visibility":                                    NewAMD("github.com/corestoreio/csfw/catalog/catattr.ProductSourceVisibility().Config(eav.AttributeSourceIdx({{.AttributeIndex}}))"),
	"Magento\\Catalog\\Model\\Product\\Attribute\\Source\\Visibility":                 NewAMD("github.com/corestoreio/csfw/catalog/catattr.ProductSourceVisibility().Config(eav.AttributeSourceIdx({{.AttributeIndex}}))"),
	"core/design_source_design":                                                       NewAMD(""),
	"customer/customer_attribute_source_group":                                        NewAMD("github.com/corestoreio/csfw/customer/custattr.CustomerSourceGroup().Config(eav.AttributeSourceIdx({{.AttributeIndex}}))"),
	"Magento\\Customer\\Model\\Customer\\Attribute\\Source\\Group":                    NewAMD("github.com/corestoreio/csfw/customer/custattr.CustomerSourceGroup().Config(eav.AttributeSourceIdx({{.AttributeIndex}}))"),
	"customer/customer_attribute_source_store":                                        NewAMD("github.com/corestoreio/csfw/customer/custattr.CustomerSourceStore().Config(eav.AttributeSourceIdx({{.AttributeIndex}}))"),
	"Magento\\Customer\\Model\\Customer\\Attribute\\Source\\Store":                    NewAMD("github.com/corestoreio/csfw/customer/custattr.CustomerSourceStore().Config(eav.AttributeSourceIdx({{.AttributeIndex}}))"),
	"customer/customer_attribute_source_website":                                      NewAMD("github.com/corestoreio/csfw/customer/custattr.CustomerSourceWebsite().Config(eav.AttributeSourceIdx({{.AttributeIndex}}))"),
	"Magento\\Customer\\Model\\Customer\\Attribute\\Source\\Website":                  NewAMD("github.com/corestoreio/csfw/customer/custattr.CustomerSourceWebsite().Config(eav.AttributeSourceIdx({{.AttributeIndex}}))"),
	"customer/entity_address_attribute_source_country":                                NewAMD("github.com/corestoreio/csfw/customer/custattr.AddressSourceCountry().Config(eav.AttributeSourceIdx({{.AttributeIndex}}))"),
	"Magento\\Customer\\Model\\Resource\\Address\\Attribute\\Source\\Country":         NewAMD("github.com/corestoreio/csfw/customer/custattr.AddressSourceCountry().Config(eav.AttributeSourceIdx({{.AttributeIndex}}))"),
	"customer/entity_address_attribute_source_region":                                 NewAMD("github.com/corestoreio/csfw/customer/custattr.AddressSourceRegion().Config(eav.AttributeSourceIdx({{.AttributeIndex}}))"),
	"Magento\\Customer\\Model\\Resource\\Address\\Attribute\\Source\\Region":          NewAMD("github.com/corestoreio/csfw/customer/custattr.AddressSourceRegion().Config(eav.AttributeSourceIdx({{.AttributeIndex}}))"),
	"eav/entity_attribute_source_boolean":                                             NewAMD("github.com/corestoreio/csfw/eav.AttributeSourceBoolean().Config(eav.AttributeSourceIdx({{.AttributeIndex}}))"),
	"Magento\\Eav\\Model\\Entity\\Attribute\\Source\\Boolean":                         NewAMD("github.com/corestoreio/csfw/eav.AttributeSourceBoolean().Config(eav.AttributeSourceIdx({{.AttributeIndex}}))"),
	"eav/entity_attribute_source_table":                                               NewAMD("github.com/corestoreio/csfw/eav.AttributeSourceTable().Config(eav.AttributeSourceIdx({{.AttributeIndex}}))"),
	"Magento\\Eav\\Model\\Entity\\Attribute\\Source\\Table":                           NewAMD("github.com/corestoreio/csfw/eav.AttributeSourceTable().Config(eav.AttributeSourceIdx({{.AttributeIndex}}))"),
	"tax/class_source_product":                                                        NewAMD("github.com/corestoreio/csfw/tax.AttributeSourceTaxClassProduct().Config(eav.AttributeSourceIdx({{.AttributeIndex}}))"),
	"Magento\\Tax\\Model\\TaxClass\\Source\\Product":                                  NewAMD("github.com/corestoreio/csfw/tax.AttributeSourceTaxClassProduct().Config(eav.AttributeSourceIdx({{.AttributeIndex}}))"),
	"Magento\\Tax\\Model\\TaxClass\\Source\\Customer":                                 NewAMD("github.com/corestoreio/csfw/tax.AttributeSourceTaxClassCustomer().Config(eav.AttributeSourceIdx({{.AttributeIndex}}))"),
	"Magento\\Theme\\Model\\Theme\\Source\\Theme":                                     NewAMD(""),
	"customer/attribute_data_postcode":                                                NewAMD("github.com/corestoreio/csfw/customer/custattr.AddressDataPostcode().Config(eav.AttributeDataIdx({{.AttributeIndex}}))"),
	"Magento\\Customer\\Model\\Attribute\\Data\\Postcode":                             NewAMD("github.com/corestoreio/csfw/customer/custattr.AddressDataPostcode().Config(eav.AttributeDataIdx({{.AttributeIndex}}))"),
}
