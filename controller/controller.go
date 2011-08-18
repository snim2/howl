/* Level of indirection between model and view packages.
 * 
 * Copyright (C) Sarah Mount, 2011.
 * 
 * This program is free software; you can redistribute it and/or
 * modify it under the terms of the GNU General Public License
 * as published by the Free Software Foundation; either version 2
 * of the License, or (at your option) any later version.
 * 
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 * 
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package controller

import (
	"appengine"
    "appengine/user"
	"appengine/datastore"
	"http"
	"log"
	"os"
)

import (
    "model"
)

/* Retreive a user object in the model, from a uid.
 */
func GetUserFromUid (context appengine.Context, uid string) (*model.HowlUser, *datastore.Key) {
	hus := make([]model.HowlUser, 0, 1)
	query := datastore.NewQuery("HowlUser").Filter("Uid =", uid).Limit(1)
	log.Println("Looking for user with id " + uid)
	keys, err := new(model.HowlUser).Query(context, query, &hus)
	if err != nil || len(keys) == 0 {
		return nil, nil
	}
	return &hus[0], keys[0]
}


/* Retreive a user object in the model, from an email address.
 */
func GetUserFromEmail (context appengine.Context, email string) (*model.HowlUser, *datastore.Key) {
	hus := make([]model.HowlUser, 0, 1)
	query := datastore.NewQuery("HowlUser").Filter("Email =", email).Limit(1)
	log.Println("Looking for user with address " + email)
	keys, err := new(model.HowlUser).Query(context, query, &hus)
	if err != nil || len(keys) == 0 {
		return nil, nil
	}
	return &hus[0], keys[0]
}


/* Retreive a user object in the model, from an appengine user object.
 */
func GetCurrentHowlUser(context appengine.Context) (*model.HowlUser, *datastore.Key) {
	email := user.Current(context).Email
	return GetUserFromEmail(context, email)
}


/* Check the uniqueness of a username in the datastore.
 *
 * Returns true if there is no such uid in the datastore, thus the uid will be 
 * unique in the store.
 */
func IsUidUnique(context appengine.Context, uid string) (bool) {
	err := (&model.HowlUser{"", uid, "", "", "", model.Now()}).Read(context)
	if err != nil {
		log.Println("Datastore error retreiving HowlUser with uid " + uid + ": " + err.String())
		return true
	} 
	return false
}


func SetLastLoggedIn(context appengine.Context, w http.ResponseWriter) (os.Error) {
	userobj, err := GetCurrentHowlUser(context)
	if userobj == nil || err != nil {
		return os.NewError("No such user: " + user.Current(context).Id)
	}
	userobj.LastLogin = model.Now()
	return userobj.Update(context)
}


/* Get a list of tags, from a list of strings.
 *
 * For each tag, check if the tag is already in the datastore.
 * If it is, return the existing tag, if not create a new one.
 * A list of keys and corresponding model.Tag objects is returned.
 *
 * It is the responsibility of the calling object to ensure that the strings
 * passed to this function have been trimmed.
 */
func PutTags(context appengine.Context, tagnames []string) ([]model.Tag, []*datastore.Key, os.Error) {
	return model.MakeTags(context, tagnames)
}
 

func CreateDataStream(context appengine.Context, name string, description string, url string, tagnames []string, pachubeKey string, pachubeFeedId int64, twitterName string, twitterToken string, twitterTokenSecret string) (os.Error) {

	key_sc, err_sc := (&model.StreamConfiguration{pachubeKey, pachubeFeedId, twitterName, twitterToken, twitterTokenSecret}).Create(context)	
	if err_sc != nil {
		return err_sc
	}
	_, tagkeys, err_tags := model.MakeTags(context, tagnames)
	if err_tags != nil {
		return err_tags
	}
	_, key_hu := GetCurrentHowlUser(context)
	_, err_ds := (&model.DataStream{key_hu, name, description, url, nil, nil, key_sc, tagkeys, nil}).Create(context)
	if err_ds != nil {
		return err_ds
	}
	return nil
}


func GetStreamsOwnedByUser (context appengine.Context, w http.ResponseWriter) ([]model.DataStream) {
	user, userKey := GetCurrentHowlUser(context)
	log.Println("Looking for datastreams owned by: " + user.Uid)
	streams := make([]model.DataStream, 0, 100) // FIXME: Magic number
	q := datastore.NewQuery("DataStream").Filter("Owner=", userKey).Limit(10) 
	if _, err := q.GetAll(context, &streams); err != nil {
		log.Println("Error fetching DataStream objects for user " + user.Uid + ": " + err.String())
        http.Error(w, err.String(), http.StatusInternalServerError)
        return nil
    }
	return streams
}

