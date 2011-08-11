/* This package routes clients to the correct handlers.
 * 
 * The URI scheme should be RESTful, in the sense of 
 * (Fielding, 2005).
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

package howl

import (
    "http"
)

import (
	"view"
)

/* init() routes requests to the appropriate handlers.
 *
 * The design of this API should be RESTful in the sense of (Fielding, 2005).
 *
 * INVARIANT: This should be the only init() method in the application.
 */
func init() {
	// Especially for web browsers.
    http.HandleFunc("/",						view.DashboardHandler)
	http.HandleFunc("/dashboard",				view.DashboardHandler)
	// RESTful interface.
    http.HandleFunc("/user",					view.UserHandler)
    http.HandleFunc("/user/([^/]+)/profile",	view.ProfileHandler)
    http.HandleFunc("/stream",					view.StreamHandler)
    http.HandleFunc("/provider",                view.ProviderHandler)
    http.HandleFunc("/datum",					view.DatumHandler)
}

