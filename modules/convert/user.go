// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package convert

import (
	"code.gitea.io/gitea/models"
	"code.gitea.io/gitea/modules/markup"
	api "code.gitea.io/gitea/modules/structs"
)

// ToUser convert models.User to api.User
// if doer is set, private information is added if the doer has the permission to see it
func ToUser(user, doer *models.User) *api.User {
	if user == nil {
		return nil
	}
	authed := false
	signed := false
	if doer != nil {
		signed = true
		authed = doer.ID == user.ID || doer.IsAdmin
	}
	return toUser(user, signed, authed)
}

// ToUserWithAccessMode convert models.User to api.User
// AccessMode is not none show add some more information
func ToUserWithAccessMode(user *models.User, accessMode models.AccessMode) *api.User {
	if user == nil {
		return nil
	}
	return toUser(user, accessMode != models.AccessModeNone, false)
}

// toUser convert models.User to api.User
// signed shall only be set if requester is logged in. authed shall only be set if user is site admin or user himself
func toUser(user *models.User, signed, authed bool) *api.User {
	result := &api.User{
		ID:         user.ID,
		UserName:   user.Name,
		FullName:   markup.Sanitize(user.FullName),
		Email:      user.GetEmail(),
		AvatarURL:  user.AvatarLink(),
		Created:    user.CreatedUnix.AsTime(),
		Restricted: user.IsRestricted,
	}
	// hide primary email if API caller is anonymous or user keep email private
	if signed && (!user.KeepEmailPrivate || authed) {
		result.Email = user.Email
	}
	// only site admin will get these information and possibly user himself
	if authed {
		result.IsAdmin = user.IsAdmin
		result.LastLogin = user.LastLoginUnix.AsTime()
		result.Language = user.Language
	}
	return result
}
