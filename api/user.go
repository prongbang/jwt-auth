package api

import (
	"jwt-auth/model"

	"cloud.google.com/go/datastore"
	"google.golang.org/api/iterator"
)

const kindUser = "User"

// FindUser from datastore
func FindUser(username, password string) (*model.User, error) {
	ctx, cancel := getContext()
	defer cancel()

	var user model.User
	q := datastore.
		NewQuery(kindUser).
		Filter("Username =", username).
		Limit(1)
	key, err := client.Run(ctx, q).Next(&user)
	if err == iterator.Done {
		// Not found
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	user.SetKey(key)
	if !user.ComparePassword(password) {
		// wrong password return like user not found
		return nil, nil
	}
	return &user, nil
}

// SaveUser to datastore
func SaveUser(user *model.User) error {
	ctx, cancel := getContext()
	defer cancel()

	var err error
	user.Stamp()
	key := user.Key()
	if key == nil {
		key = datastore.IncompleteKey(kindUser, nil)
	}
	key, err = client.Put(ctx, key, user)
	if err != nil {
		return err
	}
	user.SetKey(key)
	return nil
}
