/* Handler for dealing with the index page.
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
	"appengine/user"
	"fmt"
	"http"
	"log"
)


import (
	"controller"
	"model"
)


/* Struct to store data for the dashTemplate template.
 */
type DashboardPage struct {
	User            string // Usually a UID
	Signout         string
	OwnedStreams    []model.DataStream
	SharedStreams   []model.DataStream
	OwnedProviders  []model.DataProvider
	SharedProviders []model.DataProvider
}


/* Handle the dashboard page (URL: /). 
 *
 * Should go to a page detailing the users details, their data streams
 * and so on.
 *
 * FIXME: Shared steams
 * FIXME: Owned / shared providers
 */
func DashboardHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("DashboardHandler got request with method: " + r.Method)
	log.Println(r.URL.String())
	context, userobj := verifyLoggedIn(w, r)
	// If the user has not created a profile the app Will have already
	// redirected to login page.
	if userobj == nil { return } 
	// Get logout URL
	logout, _ := user.LogoutURL(context, "/")
	// Get streams owned by this user
	log.Println("About to look for streams owned by user.")
	streams := controller.GetStreamsOwnedByUser(context, w)
	log.Println(fmt.Sprintf("Found %v data streams for current user.", len(streams)))
	// TODO: Get streams shared with user
	// TODO: Get providers owned by user
	// TODO: Get providers shared with user
	// Render.
	dp := DashboardPage{User:userobj.Uid, Signout:logout, OwnedStreams:streams, SharedStreams:nil, OwnedProviders:nil, SharedProviders:nil}
	renderTemplateFromFile(context, dashTemplate, dp, w)
	return
}

