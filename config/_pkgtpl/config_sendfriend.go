// +build ignore

package sendfriend

import (
	"github.com/corestoreio/csfw/config/element"
	"github.com/corestoreio/csfw/store/scope"
)

// ConfigStructure global configuration structure for this package.
// Used in frontend and backend. See init() for details.
var ConfigStructure element.SectionSlice

func init() {
	ConfigStructure = element.MustNewConfiguration(
		element.Section{
			ID:        "sendfriend",
			Label:     `Email to a Friend`,
			SortOrder: 120,
			Scopes:    scope.PermStore,
			Resource:  0, // Magento_Backend::sendfriend
			Groups: element.NewGroupSlice(
				element.Group{
					ID:        "email",
					Label:     `Email Templates`,
					SortOrder: 1,
					Scopes:    scope.PermStore,
					Fields: element.NewFieldSlice(
						element.Field{
							// Path: sendfriend/email/enabled
							ID:        "enabled",
							Label:     `Enabled`,
							Type:      element.TypeSelect,
							SortOrder: 1,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							Default:   true,
							// SourceModel: Magento\Config\Model\Config\Source\Yesno
						},

						element.Field{
							// Path: sendfriend/email/template
							ID:        "template",
							Label:     `Select Email Template`,
							Comment:   text.Long(`Email template chosen based on theme fallback when "Default" option is selected.`),
							Type:      element.TypeSelect,
							SortOrder: 2,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							Default:   `sendfriend_email_template`,
							// SourceModel: Magento\Config\Model\Config\Source\Email\Template
						},

						element.Field{
							// Path: sendfriend/email/allow_guest
							ID:        "allow_guest",
							Label:     `Allow for Guests`,
							Type:      element.TypeSelect,
							SortOrder: 3,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							Default:   false,
							// SourceModel: Magento\Config\Model\Config\Source\Yesno
						},

						element.Field{
							// Path: sendfriend/email/max_recipients
							ID:        "max_recipients",
							Label:     `Max Recipients`,
							Type:      element.TypeText,
							SortOrder: 4,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							Default:   5,
						},

						element.Field{
							// Path: sendfriend/email/max_per_hour
							ID:        "max_per_hour",
							Label:     `Max Products Sent in 1 Hour`,
							Type:      element.TypeText,
							SortOrder: 5,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							Default:   5,
						},

						element.Field{
							// Path: sendfriend/email/check_by
							ID:        "check_by",
							Label:     `Limit Sending By`,
							Type:      element.TypeSelect,
							SortOrder: 6,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							Default:   false,
							// SourceModel: Magento\SendFriend\Model\Source\Checktype
						},
					),
				},
			),
		},
	)
	Backend = NewBackend(ConfigStructure)
}
