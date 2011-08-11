/* Web version of the howl view (as in MVC).
 *
 * FIXME: should just route request to different views according to requests.
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
	"appengine"
    "appengine/user"
	"fmt"
	"http"
	"io"
	"json"
	"log"
	"os"
	"strconv"
	"strings"
	"template"
)


/* // Used only for debugging.
import (
	"reflect" 
)
*/


import (
	"model"
	"controller"
)


/* HTML templates.
 *
 * These are stored in the top level directory of the app.
 */
var ( 
	signInTemplate		= template.MustParseFile("sign.html",      nil)
	streamTemplate		= template.MustParseFile("stream.html",    nil)
	dashTemplate		= template.MustParseFile("dashboard.html", nil)
	newUserTemplate		= template.MustParseFile("profile.html",   nil)
)


/* Struct to store data for the dashTemplate template.
 */
type DashboardPage struct {
	User            string
	Signout         string
	OwnedStreams    []model.DataStream
	SharedStreams   []model.DataStream
	OwnedProviders  []model.DataProvider
	SharedProviders []model.DataProvider
}


/* Struct to store data for the newUserTemplate template.
 */
type NewUserPage struct {
	User		 string
	Signout		 string
}


/* Serve a Not Found page.
 *
 * TODO: Add a styled template.
 */
func serve404(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	io.WriteString(w, "Not Found")
}


/* Serve an error page and response.
 *
 * TODO: Add an error page template.
 */
func serveError(c appengine.Context, w http.ResponseWriter, err os.Error) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	io.WriteString(w, "Internal Server Error")
	c.Errorf("%v", err)
}


/* Apply an HTML template and render.
 *
 * @param templ an HTML template 
 * @tempData usually a struct passed to the template 
 * @w the response writer that should render the template
 */
func renderTemplateFromFile(templ *template.Template, tempData interface{}, w http.ResponseWriter) {
	err := templ.Execute(w, tempData)
    if err != nil {
        http.Error(w, err.String(), http.StatusInternalServerError) // 500
    }
	return
}


/* Verify that the user is logged in, and if so, return an appengine
 * context and the appengine/user.User object corresponding to their
 * identity.
 *
 * If user.User is nil, then the user could not be logged in.
 */
func verifyLoggedIn(w http.ResponseWriter, r *http.Request) (appengine.Context, *model.HowlUser) {
	log.Println("Checking that user is logged in and has a HowlUser object")
    context := appengine.NewContext(r)
    g_user := user.Current(context)
    if g_user == nil {
		url, _ := user.LoginURL(context, r.URL.String())
		log.Println("No username, user not logged in with Google account.")
		renderTemplateFromFile(signInTemplate, url, w)
		return context, nil
    }
	userobj, _ := controller.GetCurrentHowlUser(context, w)
	if userobj == nil {
		log.Println("Cannot find a HowlUser object for this user.")
		ProfileHandler(w, r)
		return nil, nil
	}
	return context, userobj
}


// *** Handlers below ***


/* Handle the dashboard page. 
 *
 * Should go to a page detailing the users details, their data streams
 * and so on.
 *
 * FIXME: Shared steams
 * FIXME: Owned / shared providers
 */
func DashboardHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("DashboardHandler got request with method: " + r.Method)
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
	renderTemplateFromFile(dashTemplate, dp, w)
	return
}


/* Check the uniqueness of a Uid.
 *
 */
func CheckUidHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("CheckUidHandler got request with method: " + r.Method)
	if r.Method == "POST" {
		context := appengine.NewContext(r)
		isUnique := controller.IsUidUnique(context, r.FormValue("uid"))
		encoder := json.NewEncoder(w)
		type UidResponse struct {
			Uid       string
			Available string
		}
		response := new(UidResponse)
		response.Uid = r.FormValue("uid")
		if isUnique {
			response.Available = "available"
		} else {
			response.Available = "not available"
		}
		log_, _ := json.Marshal(&response)
		log.Println("Sending to JQuery: " + string(log_))
		if err := encoder.Encode(&response); err != nil {
			log.Println("Cannot send JSON to AJAX code: " + err.String())
		}
	}
	return
}


/* Handle requests relating to users. 
 *
 * By default present a view of a given user. 
 *
 * If URL is appended with ?action=edit or ?action=delete
 * then perform the appropriate CRUD action
 *
 * FIXME Look at the request code, deal with PUT, DELETE, etc. separately
 */
func UserHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("UserHandler got request with method: " + r.Method)
	return
}


/* Handle requests relating to streams. 
 *
 * By default present a view of a given data stream. 
 *
 * If URL is appended with ?action=edit or ?action=delete
 * then perform the appropriate CRUD action
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
		for i := 0; i < len(tagnames); i++ {		
			tagnames[i] = strings.Trim(tagnames[i], " ")
		}
		// Create new objects in model
		sc := model.StreamConfiguration{PachubeKey:r.FormValue("pachubekey"), PachubeFeedId:pkey, TwitterName:r.FormValue("twitteraccount"), TwitterToken:r.FormValue("twittertoken"), TwitterTokenSecret:r.FormValue("twittertokensecret")}
		ds := model.DataStream{Name:r.FormValue("name"), Description:r.FormValue("description"), Url:r.FormValue("url")}
		controller.PutDataStreamObject(sc, ds, tagnames, context, w)
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	return
}


/* Handle requests relating to data providers. 
 *
 * By default present a view of a given data provider. 
 *
 * If URL is appended with ?action=edit or ?action=delete
 * then perform the appropriate CRUD action
 *
 * FIXME Look at the request code, deal with PUT, DELETE, etc. separately
 */
func ProviderHandler(w http.ResponseWriter, r *http.Request) {
	return
}


/* Handle requests relating to data values. 
 *
 * By default present a view of a given datum. 
 *
 * If URL is appended with ?action=edit or ?action=delete
 * then perform the appropriate CRUD action
 *
 * FIXME Look at the request code, deal with PUT, DELETE, etc. separately
 */
func DatumHandler(w http.ResponseWriter, r *http.Request) {
	return
}


/* Page to configure a user profile.
 */
func ProfileHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("ProfileHandler got request with method: " + r.Method)
	if r.Method == "GET" {
		context := appengine.NewContext(r)
		g_user := user.Current(context)
		logout, _ := user.LogoutURL(context, "/")
		nup := NewUserPage{User:g_user.String(), Signout:logout}
		renderTemplateFromFile(newUserTemplate, nup, w)
		return
	}
	if r.Method == "POST" {
		context := appengine.NewContext(r)
		// Create model object
		hu := model.HowlUser{Name:r.FormValue("name"), 
		                     Uid:r.FormValue("uid"),
	                         Url:r.FormValue("url"), 
	                         About:r.FormValue("about"),}
		log.Println("Created new user profile for " + r.FormValue("name")) 
		controller.PutUserObject(hu, context, w)
		req, _ := http.NewRequest("GET", "/", r.Body)
		http.Redirect(w, req, "/", 302)
		return
	}
}

