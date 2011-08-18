/* Web version of the howl view (as in MVC).
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
	profileTemplate		= template.MustParseFile("profile.html",   nil)
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


/* Struct to store data for the profileTemplate template.
 */
type ProfilePage struct {
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
	log.Println("ERROR: 404 Not Found")
}


/* Serve an error page and response.
 *
 * TODO: Add an error page template.
 */
func serveError(c appengine.Context, w http.ResponseWriter, code int, err os.Error) {
	w.WriteHeader(code)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Write([]uint8(err.String()))
	log.Println("ERROR: " + strconv.Itoa(code) + " " + err.String())
}


/* Apply an HTML template and render.
 *
 * @param templ an HTML template 
 * @tempData usually a struct passed to the template 
 * @w the response writer that should render the template
 */
func renderTemplateFromFile(context appengine.Context, templ *template.Template, tempData interface{}, w http.ResponseWriter) {
	err := templ.Execute(w, tempData)
    if err != nil {
        serveError(context, w, http.StatusInternalServerError, err) // 500
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
		renderTemplateFromFile(context, signInTemplate, url, w)
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


/* Check the uniqueness of a Uid.
 *
 * FIXME: RESTify this!
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


type ResponseMedia int
const (
	HTML ResponseMedia = iota
	JSON
	CSV
	XML
	ATOM
)

/* Route requests relating to the RESTful interfaces. 
 *
 * FIXME: Remove nested Switch and replace with strategy pattern.
 */
func RestHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("RestHandler got request with method: " + r.Method + " and URL: " + r.URL.String())
	// Determine what media type to respond with.
	responseTy := HTML
	for _, mime := range r.Header["Accept"]  {
		if mime == "text/html" {
			responseTy = HTML
			break
		}
	}

	if r.URL.String() == "/" {
		switch responseTy {
			default: 
			DashboardHandler(w, r)
			return
		}
	}
	paths := strings.Split(r.URL.Path, "/", -1)
	if len(paths) == 0 {
		serve404(w)
		return
	}
	log.Println("Rest got initial path: " + paths[1])
	switch paths[1] {
	case "user":
		switch responseTy {
			default: 
			UserHandler(w, r)
			return
		}
	case "stream":
		switch responseTy {
			default: 
			StreamHandler(w, r)
			return
		}
	case "provider":
		switch responseTy {
			default: 
			ProviderHandler(w, r)
			return
		}
	case "datum":
		switch responseTy {
			default: 
			DatumHandler(w, r)
			return
		}
	}
	serve404(w)
	return
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


/* Handle requests relating to data providers. 
 *
 * By default present a view of a given data provider. 
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
 * FIXME Look at the request code, deal with PUT, DELETE, etc. separately
 */
func DatumHandler(w http.ResponseWriter, r *http.Request) {
	return
}


/* Page to display or configure a user profile.
 *
 * URLs pointing here should be of the form: /user/$USERNAME/ to display or
 * /user/$USERNAME/edit to edit a profile.
 *
 * TODO: Add cases for PUT, DELETE, HEAD, OPTIONS
 * TODO: What should happen to orphaned streams, providers and data on a user DELETE?
 */
func UserHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("UserHandler got request with method: " + r.Method)
	paths := strings.Split(r.URL.Path, "/", -1)
	if len(paths) < 2 {
		serve404(w)
		return
	}
	context := appengine.NewContext(r)
	g_user := user.Current(context)
	logged_in_user, _ := controller.GetUserFromEmail(context, g_user.Email)

	switch r.Method  {
	case "POST":
		// PRE: Can only edit username if you are logged in as that user
		if r.FormValue("name") != logged_in_user.Uid {
			serveError(context, w, http.StatusForbidden, os.NewError("Forbidden: You can only edit your own user profile"))
			return
		}
		// User who is logged in has no profile, need to create one.
		if logged_in_user == nil {
			hu := model.HowlUser{Name:r.FormValue("name"), 
		                         Uid:r.FormValue("uid"),
	                             Url:r.FormValue("url"), 
	                             About:r.FormValue("about"),}
			log.Println("Created new user profile for " + r.FormValue("uid")) 
			_, _ = hu.Create(context)
			req, _ := http.NewRequest("GET", "/", r.Body)
			http.Redirect(w, req, "/", http.StatusFound)
		} else {
			// User is logged in and editing existing profile.
			logged_in_user.Name = r.FormValue("name")
			logged_in_user.Uid = r.FormValue("uid")
			logged_in_user.Url =  r.FormValue("url")
			logged_in_user.Url = r.FormValue("about")
			err := logged_in_user.Update(context)
			if err != nil {
				serveError(context, w, http.StatusInternalServerError, err)
				return
			}
			req, _ := http.NewRequest("GET", "/", r.Body)
			http.Redirect(w, req, "/", http.StatusFound)
			return
		}
	case "GET":
		uname := paths[2]
		log.Println("Username: " + uname)
		userobj, _ := controller.GetUserFromUid(context, uname)
		logout, _ := user.LogoutURL(context, "/")
		npp := new(ProfilePage)
		if userobj == nil { // User has no profile
			npp.User = g_user.String()
			npp.Signout = logout
		} else {
			npp.User			= logged_in_user.Uid
			npp.Signout			= logout
			npp.Uid				= userobj.Uid
			npp.Name			= userobj.Name
			npp.Url				= userobj.Url
			npp.Description		= userobj.About
		}
		renderTemplateFromFile(context, profileTemplate, npp, w)
		return
	default:
		err := os.NewError("Could not understand your request: " + r.URL.Path)
		serveError(context, w, http.StatusBadRequest, err)
		return
	}
}

