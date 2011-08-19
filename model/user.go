/* Model of users stored in Howl. 
 * 
 * Note we use the name HowlUser so as not to shadow the name User in the 
 * appengine/user package.
 *
 * TODO: Use memcache
 * TODO: Document how keys are generated and used.
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

package model


import (
	"appengine"
	"appengine/datastore"
	"appengine/user"
	"log"
	"os"
	"time"
)


/* Data and metadata for a user account.
 *
 * The Uid field is used as the entity key.
 * A username (HowlUser.Uid) must be unique in the datastore.
 */
type HowlUser struct {
	Name				 string // "Real" name.
	Uid			         string // Unique username, checked by jquery, used as key.
	Email				 string // Taken from Google account.
	About				 string
	Url                  string
	LastLogin            datastore.Time
	Streams              []*datastore.Key
	Providers            []*datastore.Key
	SharedStreams        []*datastore.Key
	SharedProviders      []*datastore.Key
}


func (huser *HowlUser) Create(context appengine.Context) (*datastore.Key, os.Error) {
	log.Println("Creating user profile for " + huser.Name + " with id " + huser.Uid) 
	huser.Email = user.Current(context).Email
	huser.LastLogin = datastore.SecondsToTime(time.Seconds())
	huser.Streams = make([]*datastore.Key, 0, max_streams)
	huser.Providers = make([]*datastore.Key, 0, max_providers)
	huser.SharedStreams = make([]*datastore.Key, 0, max_shared_streams)
	huser.SharedProviders = make([]*datastore.Key, 0, max_shared_providers)
	key := datastore.NewKey("HowlUser", huser.Uid, 0, nil)
	err_s := "Error storing new user profile for " + huser.Uid
	return put(context, key, err_s, huser)
}


func (huser *HowlUser) Read(context appengine.Context) (os.Error) {
	hu := new(HowlUser)
	key := datastore.NewKey("HowlUser", huser.Uid, 0, nil) 
	obj, err := get(context, key, "Error fetching HowlUser object: ", hu)
	if obj, ok := obj.(HowlUser); ok {
		huser = &obj
	}
	return err
}


func (*HowlUser) Query(context appengine.Context, dsquery *datastore.Query, hus *[]HowlUser) ([]*datastore.Key, os.Error) {
	err_s := "Error fetching HowlUser object: "
	keys, err := query(context, dsquery, hus, err_s)
	return keys, err
}


func (huser *HowlUser) Update(context appengine.Context) (os.Error) {
	key := datastore.NewKey("HowlUser", huser.Uid, 0, nil) 
	_, err := put(context, key, "Error storing new HowlUser: ", huser)
	return err
}


func (huser *HowlUser) Delete(context appengine.Context) (os.Error) {
	key := datastore.NewKey("HowlUser", huser.Uid, 0, nil) 
	return delete(context, key)
}


func (huser *HowlUser) String() (string) {
	return "HowlUser object for user: " + huser.Uid + " real name: " + huser.Name
}
