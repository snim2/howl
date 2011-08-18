/* Model of data stored in Howl.
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
	"strconv"
	"strings"
	"time"
)


const (
	max_streams					= 100
	max_providers				= 100
	max_shared_streams			= 100
	max_shared_providers		= 100
)


/* Return the current time since the UNIX epoch.
 */
func Now() datastore.Time {
	return datastore.SecondsToTime(time.Seconds())
}


/* Place an object in the datastore.
 * TODO: Use memcache.
 *
 * @param context for this particular appengine session
 * @param key datastore key for the object to be stored
 * @param error message to be printed to the log / user in case of error
 * @param object the object to be made persistent
 * @return the key returned by the persistent store (if there is one, nil otherwise) and an error report (if there is one, nil otherwise)
 */
func put(context appengine.Context, key *datastore.Key, error string, object interface{}) (*datastore.Key, os.Error) {
    key_, err := datastore.Put(context, key, object)
	if err != nil {
		log.Println(error + " " + err.String())
        return nil, err
    }
	return key_, nil
}


/* Place an entity in the datastore only if no entity with the same fields eists.
 * TODO: Use memcache.
 *
 * @param context for this particular appengine session
 * @param key datastore key for the object to be stored
 * @param error message to be printed to the log / user in case of error
 * @param object the object to be made persistent
 * @return the key returned by the persistent store (if there is one, nil other
 */
func putSingleton(context appengine.Context, key *datastore.Key, error string, object interface{}) (*datastore.Key, os.Error) {
	_, err_get := get(context, key, "", object)
	if err_get == nil { // Object is already in the store.
		return key, nil
	}
	return put(context, key, error, object)
}


/* Get an object from the datastore.
 * TODO: Use memcache.
 *
 * @param context for this particular appengine session
 * @param key datastore key for the object to be retrieved
 * @param object the object to be made persistent
 * @param error message to be printed to the log / user in case of error
 * @return the retrieved object and any errors generated by the datastore
 */
func get(context appengine.Context, key *datastore.Key, error string, object interface{}) (interface{}, os.Error) {
	err := datastore.Get(context, key, object)
	if err != nil {
		log.Println(error + err.String())		
		return nil, err
	} 
	return object, nil
}


/* Query the datastore.
 * TODO: Use memcache.
 *
 * @param context for this particular appengine session
 * @param key datastore key for the object to be retrieved
 * @param object the object to be made persistent
 * @param error message to be printed to the log / user in case of error
 * @return a list of keys returned by the query and errors generated by the query (if there are any)
 */
func query(context appengine.Context, query *datastore.Query, objects interface{}, error string) ([]*datastore.Key, os.Error) {
	keys, err := query.GetAll(context, objects)
	if err != nil {
		log.Println(error + err.String())		
		return nil, err
	} 
	return keys, nil
}


/* Delete an entity from the datastore.
 * TODO: Use memcache.
 *
 * @param context for this particular appengine session
 * @param key datastore key for the object to be deleted
 * @return any errors generated by the datastore
 */
func delete(context appengine.Context, key *datastore.Key) (os.Error) {
	err := datastore.Delete(context, key)
	return err
}


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


/* A comment is a message from a (usually human) Howl user.
 */
type Comment struct {
	Text		 string
	Author		 string // HowlUser.Uid
}


func (comment *Comment) Create(context appengine.Context) (*datastore.Key, os.Error) {
	log.Println("Creating new comment")
	key := datastore.NewIncompleteKey("Comment")
	err_s := "Error storing new comment "
	return put(context, key, err_s, comment)
}


func (*Comment) Query(context appengine.Context, dsquery *datastore.Query, comments *[]Comment) ([]*datastore.Key, os.Error) {
	err_s := "Error fetching Comment object: "
	keys, err := query(context, dsquery, comments, err_s)
	return keys, err
}


func (comment *Comment) String() (string) {
	return comment.Text
}


/* A Tag is a user-defined piece of metadata.
 *
 * Tags must always be Singletons in the datastore.
 */
type Tag struct {
	Tag string
}


func MakeTags(context appengine.Context, tagnames []string) ([]Tag, []*datastore.Key, os.Error) {
	tags := make([]Tag, len(tagnames)) 
	keys := make([]*datastore.Key, len(tagnames)) 
	errs := make([]os.Error, len(tagnames))
	for i := 0; i < len(tagnames); i++ {
		tags[i] = Tag{strings.Trim(tagnames[i], " ")}
		keys[i], errs[i] = tags[i].Create(context)
	}
	for i := 0; i < len(tagnames); i++ {		
		if errs[i] != nil {
			return tags, keys, errs[i]
		}
	}
	return tags, keys, nil
}


func TagsToStrings(tags []Tag) []string {
	tagnames := make([]string, len(tags))
	for i := 0; i < len(tagnames); i++ {
		tagnames[i] = tags[i].String()
	}
	return tagnames
}


func (tag *Tag) Create(context appengine.Context) (*datastore.Key, os.Error) {
	log.Println("Creating new tag " + tag.Tag)
	key := datastore.NewKey("Tag", tag.Tag, 0, nil)
	err_s := "Error storing new tag " + tag.Tag
	return putSingleton(context, key, err_s, tag)
}


func (*Tag) Query(context appengine.Context, dsquery *datastore.Query, tags *[]Tag) ([]*datastore.Key, os.Error) {
	err_s := "Error fetching Tag object: "
	keys, err := query(context, dsquery, tags, err_s)
	return keys, err
}


func (tag *Tag) Update(context appengine.Context) (os.Error) {
	key := datastore.NewKey("Tag", tag.Tag, 0, nil) 
	_, err := putSingleton(context, key, "Error storing new Tag: ", tag.Tag)
	return err
}


func (tag *Tag) Delete(context appengine.Context) (os.Error) {
	key := datastore.NewKey("Tag", tag.Tag, 0, nil) 
	return delete(context, key)
}


func (tag *Tag) String() (string) {
	return tag.Tag
}


/* Authentication for external services.
 *
 * An incomplete key is used for this type. It should only be accessed
 * through its corresponding DataStream object.
 */
type StreamConfiguration struct {
	PachubeKey			 string
	PachubeFeedId		 int64
	TwitterName			 string
	TwitterToken		 string
	TwitterTokenSecret	 string
}


func (sc *StreamConfiguration) Create(context appengine.Context) (*datastore.Key, os.Error) {
	log.Println("Creating new stream configuration") 
	key := datastore.NewIncompleteKey("StreamConfiguration")
	err_s := "Error storing new stream configuration"
	key_, err := put(context, key, err_s, sc)
	return key_, err
}


/* A stream is a collection of data providers and related metadata.
 *
 * The key for this object is the (unique) Uid of the owner plus the
 * name of the DataStream.
 */
type DataStream struct { 
	Owner				 *datastore.Key    // Key of HowlUser object
	Name				 string
	Description			 string
	Url                  string
	AccessList			 []*datastore.Key    // HowlUsers with write access.
	Providers			 []*datastore.Key
	Configuration		 *datastore.Key
	Tags				 []*datastore.Key
	Comments			 []*datastore.Key
}


func (ds *DataStream) Create(context appengine.Context) (*datastore.Key, os.Error) {
	log.Println("Creating new data stream") 
	keyname := ds.Owner.StringID() + "@" + ds.Name
	key := datastore.NewKey("DataStream", keyname,  0, nil)
	err_s := "Error storing new data stream"
	key_, err := put(context, key, err_s, ds)
	return key_, err
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
	Data         []*datastore.Key
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

func (datum *Datum) Create(context appengine.Context) (*datastore.Key, os.Error) {
	log.Println("Creating new datum " + datum.String())
	datum.Timestamp = Now()
	key := datastore.NewIncompleteKey("Datum")
	err_s := "Error storing new datum "
	return putSingleton(context, key, err_s, datum)
}


func (datum *Datum) String() (string) {
	return strconv.Ftoa64(datum.Value, 'f', 5)
}

