/* Level of indirection between model and view packages.

Copyright (C) Sarah Mount, 2011.

This program is free software; you can redistribute it and/or
modify it under the terms of the GNU General Public License
as published by the Free Software Foundation; either version 2
of the License, or (at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

package controller

import (
	"appengine"
    "appengine/user"
	"appengine/datastore"
	"http"
	"log"
//	"reflect"
	"time"
)

import (
    "model"
)


// TODO Write this
func SetLastLoggedIn(context appengine.Context) {
	return
}

/* Retreive a user object in the model, from an appengine user object.
 */
func GetUserObject(context appengine.Context, w http.ResponseWriter) (*model.HowlUser) {
	hu := new(model.HowlUser)
	key := datastore.NewKey("HowlUser", user.Current(context).String(), 0, nil)

	log.Println("Looking for user with Id " + user.Current(context).String())

	if err := datastore.Get(context, key, hu); err != nil {
		log.Println("Error fetching HowlUser object: " + err.String())
        return nil
    }
	return hu
}

/* Store new user object.
 */
func PutUserObject (hu model.HowlUser, context appengine.Context, w http.ResponseWriter) {
	// Set values known to the datastore
	hu.LastLogin = datastore.SecondsToTime(time.Seconds())
	hu.Id = user.Current(context).String()
	hu.Email = user.Current(context).Email
	key := datastore.NewKey("HowlUser", hu.Id, 0, nil)
	// Make persistent
    _, err := datastore.Put(context, key, &hu)
	if err != nil {
        http.Error(w, "Error storing new user profile: " + err.String(), http.StatusInternalServerError)
        return
    }
	log.Println("Made persistent new user profile for " + hu.Name + " with id " + hu.Id) 
	return
}

func PutStreamConfigurationObject (sc model.StreamConfiguration, context appengine.Context, w http.ResponseWriter) {
	key := datastore.NewIncompleteKey("StreamConfiguration")
    _, err := datastore.Put(context, key, &sc)
	if err != nil {
        http.Error(w, "Error storing stream configuration: " + err.String(), http.StatusInternalServerError)
        return
    }
	return
}

func PutDataStreamObject(ds model.DataStream, context appengine.Context, w http.ResponseWriter) {
	key := datastore.NewKey("DataStream", user.Current(context).String() + ds.Name, 0, nil)
    _, err := datastore.Put(context, key, &ds)
	if err != nil {
        http.Error(w, "Error storing data stream: " + err.String(), http.StatusInternalServerError)
        return
    }
	return
}


func GetStreamsOwnedByUser (user *model.HowlUser, context appengine.Context, w http.ResponseWriter) ([]model.DataStream) {
	if user == nil {
		return nil
	}
	streams := make([]model.DataStream, 0, 100) // FIXME: Magic number
	q := datastore.NewQuery("DataStream").Filter("Owner=", user.Id).Limit(10) 
	if _, err := q.GetAll(context, &streams); err != nil {
		log.Println("Error fetching DataStream objects for user " + user.Id + ": " + err.String())
        http.Error(w, err.String(), http.StatusInternalServerError)
        return nil
    }
	return streams
}


func GetAllStreamsUserCanAccess (user model.HowlUser, w http.ResponseWriter) ([]model.DataStream) {
	return nil
}

