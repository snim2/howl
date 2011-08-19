/* Model of data streams stored in Howl.
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
	"log"
	"os"
)


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

