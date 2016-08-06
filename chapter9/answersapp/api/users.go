package api

import (
	"crypto/md5"
	"errors"
	"fmt"
	"io"
	"strings"

	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/user"
)

type User struct {
	Key         *datastore.Key `json:"id" datastore:"-"`
	UserID      string         `json:"-"`
	DisplayName string         `json:"display_name"`
	AvatarURL   string         `json:"avatar_url"`
}

// UserFromAEUser gets or creates a User from the currently
// logged in user.
func UserFromAEUser(ctx context.Context) (*User, error) {
	aeuser := user.Current(ctx)
	if aeuser == nil {
		return nil, errors.New("not logged in")
	}
	var appUser User
	appUser.Key = datastore.NewKey(ctx, "User", aeuser.ID, 0, nil)
	err := datastore.Get(ctx, appUser.Key, &appUser)
	if err == datastore.ErrNoSuchEntity {
		appUser.UserID = aeuser.ID
		appUser.DisplayName = aeuser.String()
		appUser.AvatarURL = gravatarURL(aeuser.Email)
		log.Infof(ctx, "saving new user: %s", aeuser.String())
		appUser.Key, err = datastore.Put(ctx, appUser.Key, &appUser)
		if err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}
	return &appUser, nil
}

func gravatarURL(email string) string {
	m := md5.New()
	io.WriteString(m, strings.ToLower(email))
	return fmt.Sprintf("//www.gravatar.com/avatar/%x", m.Sum(nil))
}

func (u *User) Card() UserCard {
	return UserCard{
		Key:         u.Key,
		DisplayName: u.DisplayName,
		AvatarURL:   u.AvatarURL,
	}
}

type UserCard struct {
	Key         *datastore.Key `json:"id"`
	DisplayName string         `json:"display_name"`
	AvatarURL   string         `json:"avatar_url"`
}
