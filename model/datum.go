/* Model of data points stored in Howl.
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
	"strconv"
)


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
	Annotation   *datastore.Key
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

