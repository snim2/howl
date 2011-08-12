/* Model of data stored in Howl.
 * TODO: Document how keys are generated and used.
 * TODO: Move functions from the controller package to methods here.
 * TODO: To make the RESTful interface simpler, each model.type should implement CRUD methods
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
	"os" // TODO decide whether to use os.NewError or a custom type
	"time"
)


/* Data and metadata for a user account.
 *
 * A username must be unique in the datastore.
 */
type HowlUser struct {
	Name				 string // "Real" name.
	Uid			         string // Unique username, checked by jquery, used as key.
	Email				 string // Taken from Google account.
	About				 string
	Url                  string
	LastLogin            datastore.Time
}


func (*HowlUser) Create(context appengine.Context, name string, uid string, about string, url string) (*datastore.Key, *os.Error) {
	email := user.Current(context).Email
	login := datastore.SecondsToTime(time.Seconds())
	hu := &HowlUser{name, uid, email, about, url, login}
	key := datastore.NewKey("HowlUser", hu.Uid, 0, nil)
    _, err := datastore.Put(context, key, &hu)
	if err != nil {
        log.Println("Error storing new user profile: " + err.String())
        return nil, &err
    }
	log.Println("Created persistent new user profile for " + hu.Name + " with id " + hu.Uid) 
	return key, nil
}


func (*HowlUser) Read(context appengine.Context, key *datastore.Key) (*HowlUser, *os.Error) {
	hu := new(HowlUser)
	err := datastore.Get(context, key, hu)
	if err != nil {
		log.Println("Error fetching HowlUser object: " + err.String())		
		return nil, &err
	} 
	return hu, nil
}


func (*HowlUser) Query(context appengine.Context, query *datastore.Query) ([]HowlUser, []*datastore.Key, *os.Error) {
	hus := make([]HowlUser, 0, 100) // FIXME magic number
	keys, err := query.GetAll(context, &hus); 
	if err != nil {
		log.Println("Error fetching HowlUser object: " + err.String())
        return nil, nil, &err
    }
	return hus, keys, nil
}


func (*HowlUser) Update(context appengine.Context, newUser *HowlUser) (*os.Error) {
	key := datastore.NewKey("HowlUser", newUser.Uid, 0, nil) 
    _, err := datastore.Put(context, key, &newUser)
	if err != nil {
        log.Println("Error storing new entity: " + err.String())
        return &err
    }
	return nil
}


func (*HowlUser) Delete(context appengine.Context, key *datastore.Key) (*os.Error) {
	err := datastore.Delete(context, key)
	return &err
}


/* Authentication for external services.
 *
 */
type StreamConfiguration struct {
	PachubeKey			 string
	PachubeFeedId		 int64
	TwitterName			 string
	TwitterToken		 string
	TwitterTokenSecret	 string
	YfdToken			 string
}


/* A comment is a message from a (usually human) Howl user.
 */
type Comment struct {
	Text		 string
	Author		 string // HowlUser.Id
}


/* A Tag is a user-defined piece of metadata.
 *
 * Tags must always be Singletons in the datastore.
 */
type Tag struct {
	Tag string
}


/* A stream is a collection of data providers and related metadata.
 *
 */
type DataStream struct { 
	Owner				 *datastore.Key
	Name				 string
	Description			 string
	Url                  string
	AccessList			 []*datastore.Key    // Users with read/write access.
	Providers			 []*datastore.Key
	Configuration		 *datastore.Key
	Tags				 []*datastore.Key
	Comments			 []*datastore.Key
}


/* A provider is any entity which provides data to Howl.
 *
 * It contains metadata such as its pysical location and owner.
 */
type DataProvider struct {
	Name		 string 
	Description	 string
	Url			 string
	Owner		 *datastore.Key
	AccessList	 []*datastore.Key    // "Shared" users with read/write access.
	Latitude	 float32
	Longditude	 float32
	Elevation	 float32
	Dimension	 string              // Unit of dimension
	Data         []datastore.Key
}


/* A datum is an individual piece of raw data provided to Howl by a
 * DataProvider.
 *
 * Value should be used for numerical data, alternatively a Url may be
 * provided and the data stored on external services such as YFrog.
 */
type Datum struct {
	Timestamp	 datastore.Time
	Value		 float64
	Url			 string
}

