/* Model of data stored in Howl.
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
    "appengine/datastore"
)


/* Data and metadata for a user account.
 *
 * A username must be unique in the datastore.
 */
type HowlUser struct {
	Name				 string // "Real" name.
	Uid			         string
	Email				 string // Taken from Google account.
	About				 string
	Url                  string
	LastLogin            datastore.Time // FIXME See controller package.
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

