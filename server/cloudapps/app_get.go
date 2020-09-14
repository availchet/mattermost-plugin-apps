// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package cloudapps

import "github.com/mattermost/mattermost-plugin-cloudapps/server/utils"

func (r *registry) GetApp(appID AppID) (*App, error) {
	if appID != "hello" {
		return nil, utils.ErrNotFound
	}
	return &App{
		AppID:       "hello",
		DisplayName: "Hallo სამყარო",
		RootURL:     "https://levb.ngrok.io/plugin/cloudapps/hello",
	}, nil
}
