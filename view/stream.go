/* Handler for dealing with URIs starting with /stream/
 *
 * FIXME: Remove spaghetti code and replace with Strategy Pattern.
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
 *
 */

package view

import (
	"http"
	"log"
	"strings"
	"strconv"
)


import (
	"controller"
)


/* Handle requests relating to streams. 
 *
 * By default present a view of a given data stream. 
 *
 * FIXME Look at the request code, deal with PUT, DELETE, etc. separately
 */
func StreamHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("StreamHandler got request with method: " + r.Method)
	if r.Method == "POST" {
		context, userobj := verifyLoggedIn(w, r)
		if userobj == nil { return } // Will have already redirected to login page.
		// Reformat form data
		pkey, errkey := strconv.Atoi64(r.FormValue("pachubefeedid"))
		if errkey != nil {
			pkey = 0
		}
		// Get a list of tags
		tagnames := strings.Split(r.FormValue("tags"), ",", -1)
		// Create new objects in model
		_ = controller.CreateDataStream(context, r.FormValue("name"), r.FormValue("description"), r.FormValue("url"), tagnames, r.FormValue("pachubekey"), pkey, r.FormValue("twitteraccount"), r.FormValue("twittertoken"), r.FormValue("twittertokensecret"))
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	return
}
