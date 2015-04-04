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

package tax

import "github.com/corestoreio/csfw/eav"

var (
	_ eav.AttributeSourceModeller = (*todoASTCP)(nil)
	_ eav.AttributeSourceModeller = (*todoASTCC)(nil)
)

type (
	todoASTCP struct {
		*eav.AttributeSource
	}
	todoASTCC struct {
		*eav.AttributeSource
	}
)

// AttributeSourceTaxClassProduct @todo
// @see magento2/site/app/code/Magento/Tax/Model/TaxClass/Source/Product.php
func AttributeSourceTaxClassProduct() *todoASTCP {
	return &todoASTCP{}
}

// AttributeSourceTaxClassCustomer @todo
// @see magento2/site/app/code/Magento/Tax/Model/TaxClass/Source/Customer.php
func AttributeSourceTaxClassCustomer() *todoASTCC {
	return &todoASTCC{}
}