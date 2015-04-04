// Copyright 2015 CoreStore Authors
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

package eav

var (
	_ AttributeBackendModeller = (*todoABDT)(nil)
	_ AttributeBackendModeller = (*todoABTC)(nil)
	_ AttributeBackendModeller = (*todoABTU)(nil)
)

type (
	todoABDT struct {
		*AttributeBackend
	}
	todoABTC struct {
		*AttributeBackend
	}
	todoABTU struct {
		*AttributeBackend
	}
)

// AttributeBackendDatetime handles date times @todo
// @see magento2/site/app/code/Magento/Eav/Model/Entity/Attribute/Backend/Datetime.php
func AttributeBackendDatetime() *todoABDT {
	return &todoABDT{}
}

// AttributeBackendTimeCreated @todo
// @see magento2/site/app/code/Magento/Eav/Model/Entity/Attribute/Backend/Time/Created.php
func AttributeBackendTimeCreated() *todoABTC {
	return &todoABTC{}
}

// AttributeBackendTimeUpdated @todo
// @see magento2/site/app/code/Magento/Eav/Model/Entity/Attribute/Backend/Time/Updated.php
func AttributeBackendTimeUpdated() *todoABTU {
	return &todoABTU{}
}