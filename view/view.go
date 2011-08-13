/* Web version of the howl view (as in MVC).
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


// Used only for debugging.
import (
	"reflect" 
)



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
	Uid          string
	Name         string
	Url          string
	Description  string
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
	userobj, _ := controller.GetCurrentHowlUser(context)
	if userobj == nil {
		log.Println("Cannot find a HowlUser object for this user.")
		UserHandler(w, r)
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
		log.Println("Result from IsUidUnique is: " + strconv.Btoa(isUnique) + " : " + reflect.TypeOf(isUnique).String())
		if isUnique {
			log.Println("Username " + response.Uid + " is available")
			response.Available = "available"
		} else {
			log.Println("Username " + response.Uid + " is not available")
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


// ************************************************************************** //


/* Route requests relating to the RESTful interfaces. 
 *
 */
func RestHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("RestHandler got request with method: " + r.Method + " and URL: " + r.URL.String())
	if r.URL.String() == "/" {
		DashboardHandler(w, r)
		return
	}
	paths := strings.Split(r.URL.Path, "/", -1)
	if len(paths) == 0 {
		serve404(w)
		return
	}
	log.Println("Rest got initial path: " + paths[1])
	switch paths[1] {
	case "user":
		UserHandler(w, r)
		return
	case "stream":
		StreamHandler(w, r)
		return
	case "provider":
		ProviderHandler(w, r)
		return
	}
	serve404(w)
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
		// Create new objects in model
		_ = controller.CreateDataStream(context, r.FormValue("name"), r.FormValue("description"), r.FormValue("url"), tagnames, r.FormValue("pachubekey"), pkey, r.FormValue("twitteraccount"), r.FormValue("twittertoken"), r.FormValue("twittertokensecret"))
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
 *
 * By default present a view of a given user. 
 *
 * If URL is appended with ?action=edit or ?action=delete
 * then perform the appropriate CRUD action
 *

 */
func UserHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("UserHandler got request with method: " + r.Method)
	if r.Method == "GET" {
		context := appengine.NewContext(r)
		hu, _ := controller.GetCurrentHowlUser(context)
		g_user := user.Current(context)
		logout, _ := user.LogoutURL(context, "/")
		nup := new(NewUserPage)
		if hu == nil {
			nup.User = g_user.String()
			nup.Signout = logout
		} else {
			nup.User = hu.Uid
			nup.Signout = logout
			nup.Uid = hu.Uid
			nup.Name = hu.Name
			nup.Url = hu.Url
			nup.Description = hu.About
		}
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
		log.Println("Created new user profile for " + r.FormValue("uid")) 
		_, _ = hu.Create(context)
		req, _ := http.NewRequest("GET", "/", r.Body)
		http.Redirect(w, req, "/", 302)
		return
	}
}

