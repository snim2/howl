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
	"reflect"
	"time"
)

import (
    "model"
)


// TODO: Factor out serving errors
// TODO: Factor our PUTs and GETs (remember singleton get/put)
// TODO: Turn PUTs / GETs into Memcache calls


/* Place an object in the datastore.
 * TODO: Use memcache.
 *
 * @param context for this particular appengine session
 * @param key datastore key for the object to be stored
 * @param error message to be printed to the log / user in case of error
 * @param object the object to be made persistent
 * @return the key returned by the persistent store (if there is one, nil otherwise) and an error report (if there is one, nil otherwise)
 */
func put(context appengine.Context, key datastore.Key, error string, object interface{}) (*datastore.Key, os.Error) {
    key_, err := datastore.Put(context, &key, &object)
	if err != nil {
		log.Println(error + " " + err.String())
//        http.Error(w, "Error storing new user profile: " + err.String(), http.StatusInternalServerError)
        return nil, err
    }
	return key_, nil
}


// FIXME: This seems to be storing the current time rather than a static timestamp
func SetLastLoggedIn(context appengine.Context, w http.ResponseWriter) (os.Error) {
	userobj, err := GetUserObject(context, w)
	if userobj == nil || err != nil {
		return os.NewError("No such user: " + user.Current(context).Id)
	}
	userobj.LastLogin = datastore.SecondsToTime(time.Seconds())
	PutUserObject(*userobj, context, w)
	return nil
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
func GetTags(tagnames []string, context appengine.Context, w http.ResponseWriter) ([]*datastore.Key, []model.Tag) {
	name := new(string)
	key := new(datastore.Key)
	tag := new(model.Tag)
	tags :=  make([]model.Tag, len(tagnames)) 
	keys :=  make([]*datastore.Key, len(tagnames)) 
	for i := 0; i < len(tagnames); i++ {		
		name = &tagnames[i]
		log.Println(reflect.TypeOf(name))
		key = datastore.NewKey("Tag", *name, 0, nil)
		keys[i] = key
		tag = &model.Tag{Tag:*name}
		if err := datastore.Get(context, key, tag); err != nil {
			log.Println("No such tag: " + *name)
			_, err2 := datastore.Put(context, key, tag)
			if err2 != nil {
				log.Println("Could not store Tag: " + *name + " " + err2.String())
			}
		}
		tags[i] = *tag
    }
	return keys, tags
}


/* Retreive a user object in the model, from an appengine user object.
 */
func GetUserObject(context appengine.Context, w http.ResponseWriter) (*model.HowlUser, *datastore.Key) {
	hu := new(model.HowlUser)
	key := datastore.NewKey("HowlUser", user.Current(context).String(), 0, nil)
	log.Println("Looking for user with Id " + user.Current(context).String())
	if err := datastore.Get(context, key, hu); err != nil {
		log.Println("Error fetching HowlUser object: " + err.String())
        return nil, key
    }
	return hu, key
}


/* Store new user object.
 */
func PutUserObject (hu model.HowlUser, context appengine.Context, w http.ResponseWriter) {
	// Set values known to the datastore
	hu.LastLogin = datastore.SecondsToTime(time.Seconds())
	hu.Uid = user.Current(context).String()
	hu.Email = user.Current(context).Email
	key := datastore.NewKey("HowlUser", hu.Uid, 0, nil)
	// Make persistent
    _, err := datastore.Put(context, key, &hu)
	if err != nil {
        http.Error(w, "Error storing new user profile: " + err.String(), http.StatusInternalServerError)
        return
    }
	log.Println("Made persistent new user profile for " + hu.Name + " with id " + hu.Uid) 
	return
}


/* Store a datastream, and its configurtation in the persistent store.
 *
 * @param sc configuration object for the datastream
 * @param ds new datastream objects
 * @param tagnames list of strings representing tags for the datastream
 * @param context 
 * @param w writer used to write error messages
 */
func PutDataStreamObject(sc model.StreamConfiguration, ds model.DataStream, 
	                     tagnames []string, context appengine.Context, 
	                     w http.ResponseWriter) {
	// Deal with keys
	userObj, userKey := GetUserObject(context, w)
	dsKey := datastore.NewKey("DataStream", userObj.Uid + ds.Name, 0, nil)
	scKey := datastore.NewKey("StreamConfiguration", "Config" + userObj.Uid + ds.Name, 0, nil)
	// Deal with tags.
	tagKeys, _ := GetTags(tagnames, context, w)
	ds.Tags = tagKeys
	// Store stream configuration
    _, err := datastore.Put(context, scKey, &sc)
	if err != nil {
        http.Error(w, "Error storing stream configuration: " + err.String(), http.StatusInternalServerError)
		return
    }
	// Store data stream
	ds.Owner = userKey
	ds.Configuration = scKey
    _, err = datastore.Put(context, dsKey, &ds)
	if err != nil {
        http.Error(w, "Error storing data stream: " + err.String(), http.StatusInternalServerError)
        return
    }
	return
}


func GetStreamsOwnedByUser (context appengine.Context, w http.ResponseWriter) ([]model.DataStream) {
	user, userKey := GetUserObject(context, w)
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


func GetAllStreamsUserCanAccess (user model.HowlUser, w http.ResponseWriter) ([]model.DataStream) {
	return nil
}

