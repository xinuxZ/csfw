// +build ignore

package shipping

import "github.com/corestoreio/csfw/config"

var PackageConfiguration = config.NewConfiguration(
	&config.Section{
		ID:        "shipping",
		Label:     "Shipping Settings",
		SortOrder: 310,
		Scope:     config.NewScopePerm(config.IDScopeDefault, config.IDScopeWebsite),
		Groups: config.GroupSlice{
			&config.Group{
				ID:        "origin",
				Label:     `Origin`,
				Comment:   ``,
				SortOrder: 1,
				Scope:     config.NewScopePerm(config.IDScopeDefault, config.IDScopeWebsite),
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `shipping/origin/country_id`,
						ID:           "country_id",
						Label:        `Country`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    10,
						Visible:      config.VisibleYes,
						Scope:        config.NewScopePerm(config.IDScopeDefault, config.IDScopeWebsite),
						Default:      `US`,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Directory\Model\Config\Source\Country
					},

					&config.Field{
						// Path: `shipping/origin/region_id`,
						ID:           "region_id",
						Label:        `Region/State`,
						Comment:      ``,
						Type:         config.TypeText,
						SortOrder:    20,
						Visible:      config.VisibleYes,
						Scope:        config.NewScopePerm(config.IDScopeDefault, config.IDScopeWebsite),
						Default:      12,
						BackendModel: nil,
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `shipping/origin/postcode`,
						ID:           "postcode",
						Label:        `ZIP/Postal Code`,
						Comment:      ``,
						Type:         config.TypeText,
						SortOrder:    30,
						Visible:      config.VisibleYes,
						Scope:        config.NewScopePerm(config.IDScopeDefault, config.IDScopeWebsite),
						Default:      90034,
						BackendModel: nil,
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `shipping/origin/city`,
						ID:           "city",
						Label:        `City`,
						Comment:      ``,
						Type:         config.TypeText,
						SortOrder:    40,
						Visible:      config.VisibleYes,
						Scope:        config.NewScopePerm(config.IDScopeDefault, config.IDScopeWebsite),
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `shipping/origin/street_line1`,
						ID:           "street_line1",
						Label:        `Street Address`,
						Comment:      ``,
						Type:         config.TypeText,
						SortOrder:    50,
						Visible:      config.VisibleYes,
						Scope:        config.NewScopePerm(config.IDScopeDefault, config.IDScopeWebsite),
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `shipping/origin/street_line2`,
						ID:           "street_line2",
						Label:        `Street Address Line 2`,
						Comment:      ``,
						Type:         config.TypeText,
						SortOrder:    60,
						Visible:      config.VisibleYes,
						Scope:        config.NewScopePerm(config.IDScopeDefault, config.IDScopeWebsite),
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil,
					},
				},
			},
		},
	},
	&config.Section{
		ID:        "carriers",
		Label:     "Shipping Methods",
		SortOrder: 320,
		Scope:     config.ScopePermAll,
		Groups:    config.GroupSlice{},
	},
)