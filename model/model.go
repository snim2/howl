/* Model of data stored in Howl.

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

package model

import (
    "appengine/datastore"
)

type HowlUser struct {
	Name string
	Id string
    Email string
	About string
    Url string
    LastLogin datastore.Time
	DisplayStartupDocs bool
}

// Authentication for external services.

type StreamConfiguration struct {
	PachubeKey string
	PachubeFeedId int64
	TwitterName string
	TwitterToken string
	TwitterTokenSecret string
	YfdToken string
}

// Streams and data
// Documented in an ER diagram in the docs/ directory.

type Comment struct {
	Text string
	Author string // HowlUser.Id
}

type Tag struct {
	Tag string
}

type DataStream struct { 
	Owner *datastore.Key
	Name string
	Description string
	AccessList []*datastore.Key    // Users with read/write access.
	Providers []*datastore.Key
	Configuration *datastore.Key
	Tags []*datastore.Key
	Comments []*datastore.Key
}

type DataProvider struct {
	Name string 
	Description string
	Url string
	Owner *datastore.Key
	AccessList []*datastore.Key    // Users with read/write access.
	Latitude float32
	Longditude float32
	Elevation float32
	Dimension string // Unit of dimension
    Data []datastore.Key
}

type Datum struct {
	Timestamp datastore.Time
	Value float64
	Url string
}

