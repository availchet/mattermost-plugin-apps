package apps

import (
	"net/url"
	"unicode"

	"github.com/pkg/errors"

	"github.com/mattermost/mattermost-server/v5/model"
)

// AppID is a globally unique identifier that represents a Mattermost App.
// Allowed characters are letters, numbers, underscores and hyphens.
type AppID string

// AppVersion is the version of a Mattermost App.
// Allowed characters are letters, numbers, underscores and hyphens.
type AppVersion string

type AppType string

// default is HTTP
const (
	AppTypeHTTP      AppType = "http"
	AppTypeAWSLambda AppType = "aws_lambda"
	AppTypeBuiltin   AppType = "builtin"
)

// Function describes app's function mapping
// For now Function can be either AWS Lambda or HTTP function
type Function struct {
	Path    string `json:"path"`
	Name    string `json:"name"`
	Handler string `json:"handler"`
	Runtime string `json:"runtime"`
}

type Manifest struct {
	AppID   AppID      `json:"app_id"`
	Type    AppType    `json:"app_type"`
	Version AppVersion `json:"version"`

	DisplayName string `json:"display_name,omitempty"`
	Description string `json:"description,omitempty"`
	HomepageURL string `json:"homepage_url,omitempty"`

	// See DefaultInstallCall, DefaultBindingCall, etc. for the defaults. The
	// App developer can override the defaults by providing explicit Path,
	// Expand values.
	OnDisable        *Call `json:"on_disable,omitempty"`
	OnEnable         *Call `json:"on_enable,omitempty"`
	OnInstall        *Call `json:"on_install,omitempty"`
	OnVersionChanged *Call `json:"on_version_changed,omitempty"`
	OnUninstall      *Call `json:"on_uninstall,omitempty"`
	Bindings         *Call `json:"bindings,omitempty"`

	// For HTTP Apps all paths are relative to the RootURL.
	HTTPRootURL string `json:"root_url,omitempty"`

	RequestedPermissions Permissions `json:"requested_permissions,omitempty"`

	// RequestedLocations is the list of top-level locations that the
	// application intends to bind to, e.g. `{"/post_menu", "/channel_header",
	// "/command/apptrigger"}``.
	RequestedLocations Locations `json:"requested_locations,omitempty"`

	// Functions must be included by the developer in the published manifest for
	// AWS apps. These declarations are used to:
	// - create AWS Lambda functions that will service requests in Mattermost
	// Cloud;
	// - define path->function mappings, aka "routes". The function with the
	// path matching as the longest prefix is used to handle a Call request.
	Functions []Function
}

// App describes an App installed on a Mattermost instance. App should be
// abbreviated as `app`.
type App struct {
	Manifest

	Disabled bool `json:"disabled,omitempty"`

	// Secret is used to issue JWT
	Secret string `json:"secret,omitempty"`

	OAuth2ClientID     string `json:"oauth2_client_id,omitempty"`
	OAuth2ClientSecret string `json:"oauth2_client_secret,omitempty"`
	OAuth2TrustedApp   bool   `json:"oauth2_trusted_app,omitempty"`

	BotUserID      string `json:"bot_user_id,omitempty"`
	BotUsername    string `json:"bot_username,omitempty"`
	BotAccessToken string `json:"bot_access_token,omitempty"`

	// Grants should be scopable in the future, per team, channel, post with
	// regexp.
	GrantedPermissions Permissions `json:"granted_permissions,omitempty"`

	// GrantedLocations contains the list of top locations that the
	// application is allowed to bind to.
	GrantedLocations Locations `json:"granted_locations,omitempty"`
}

// ListedApp is a Mattermost App listed in the Marketplace containing metadata.
type ListedApp struct {
	Manifest  *Manifest                `json:"manifest"`
	Installed bool                     `json:"installed"`
	Enabled   bool                     `json:"enabled"`
	Labels    []model.MarketplaceLabel `json:"labels,omitempty"`
}

const (
	DefaultInstallCallPath  = "/install"
	DefaultBindingsCallPath = "/bindings"
)

var DefaultInstallCall = &Call{
	Path: DefaultInstallCallPath,
	Expand: &Expand{
		App:              ExpandAll,
		AdminAccessToken: ExpandAll,
	},
}

var DefaultBindingsCall = &Call{
	Path: DefaultBindingsCallPath,
}

func (m *Manifest) IsValid() error {
	if m.AppID == "" {
		return errors.New("empty AppID")
	}
	if m.Type == "" {
		return errors.New("app_type is empty, must be specified, e.g. `aws_lamda`")
	}
	if !m.Type.IsValid() {
		return errors.Errorf("invalid type: %s", m.Type)
	}

	if m.Type == AppTypeHTTP {
		_, err := url.Parse(m.HTTPRootURL)
		if err != nil {
			return errors.Wrapf(err, "invalid manifest URL %q", m.HTTPRootURL)
		}
	}
	return nil
}

func (id AppID) IsValid() error {
	for _, c := range id {
		if unicode.IsLetter(c) {
			continue
		}

		if unicode.IsNumber(c) {
			continue
		}

		if c == '-' || c == '_' {
			continue
		}

		return errors.Errorf("invalid character %v in appID", c)
	}

	return nil
}

func (v AppVersion) IsValid() error {
	for _, c := range v {
		if unicode.IsLetter(c) {
			continue
		}

		if unicode.IsNumber(c) {
			continue
		}

		if c == '-' || c == '_' {
			continue
		}

		return errors.Errorf("invalid character %v in appVersion", c)
	}

	return nil
}

func (at AppType) IsValid() bool {
	return at == AppTypeHTTP ||
		at == AppTypeAWSLambda ||
		at == AppTypeBuiltin
}
