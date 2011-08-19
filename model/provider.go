/* Model of data providers stored in Howl.
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
	"appengine/datastore"
)


/* A provider is any entity which provides data to Howl.
 *
 * It contains metadata such as its pysical location and owner.
 */
type DataProvider struct {
	Name		 string 
	Description	 string
	Owner		 *datastore.Key
	Url			 string
	AccessList	 []*datastore.Key    // "Shared" users with read/write access.
	Latitude	 float32
	Longditude	 float32
	Elevation	 float32
	Dimension	 string              // Unit of dimension
	Data         []*datastore.Key
}

