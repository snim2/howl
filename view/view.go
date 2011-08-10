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
	"log"
//	"reflect"
	"strconv"
	"strings"
	"template"
)

import (
	"model"
	"controller"
)


var ( // HTML templates
	signInTemplate		= template.MustParseFile("sign.html", nil)
	streamTemplate		= template.MustParseFile("stream.html", nil)
	dashTemplate		= template.MustParseFile("dashboard.html", nil)
	newUserTemplate		= template.MustParseFile("newuser.html", nil)
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


/* Apply an html template and render.
 *
 * @param templ an html template 
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
func verifyLoggedIn(w http.ResponseWriter, r *http.Request) (appengine.Context, *user.User) {
    context := appengine.NewContext(r)
    uname := user.Current(context)
    if uname == nil {
		_, err := user.LoginURL(context, r.URL.String())
        if err != nil {
            http.Error(w, err.String(), http.StatusInternalServerError) // 500
			return context, nil
        }
		CreateNewUserHandler(w, r)
//		renderTemplateFromFile(signInTemplate, url, w)
		return context, nil
    }
	return context, uname
}


/* Handle the index page. 
 *
 * This is used when the user first logs in, to set their 
 * LastLoggedIn field, then pass control to the DashboardHandler.
 */
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	context, uname := verifyLoggedIn(w, r)
	if uname == nil { 
		login_url, _ := user.LoginURL(context, r.URL.String())
		http.Error(w, "You may not access this page until you are <a href=\"" + login_url + "\"logged in.</a> ", http.StatusForbidden) // 403
	}
	// Kernel PANIC! FIXME
	controller.SetLastLoggedIn(context, w)
	DashboardHandler(w, r)
	return
}


/* Handle the dashboard page. 
 *
 * Should go to a page detailing the users details, their data streams
 * and so on.
 *
 * FIXME: Shared steams
 * FIXME: Owned / shared providers
 */
func DashboardHandler(w http.ResponseWriter, r *http.Request) {
	context, uname := verifyLoggedIn(w, r)
	if uname == nil { 
		login_url, _ := user.LoginURL(context, r.URL.String())
		http.Error(w, "You may not access this page until you are <a href=\"" + login_url + "\"logged in.</a> ", http.StatusForbidden) // 403
	}
	// Get logout URL
	logout, _ := user.LogoutURL(context, "/")
	// Get streams owned by this user
	streams := controller.GetStreamsOwnedByUser(context, w)
	log.Println(fmt.Sprintf("Found %v data streams for current user.", len(streams)))
//	log.Println(reflect.TypeOf(streams).String())
	// TODO: Get streams shared with user
	// TODO: Get providers owned by user
	// TODO: Get providers shared with user
	// Render.
	dp := DashboardPage{User:uname.String(), Signout:logout, OwnedStreams:streams, SharedStreams:nil, OwnedProviders:nil, SharedProviders:nil}
	renderTemplateFromFile(dashTemplate, dp, w)
	return
}


/* Page to configure a new user profile.
 */
func NewUserHandler(w http.ResponseWriter, r *http.Request) {
	context, uname := verifyLoggedIn(w, r)
	if uname == nil {
		login_url, _ := user.LoginURL(context, r.URL.String())
		http.Error(w, "You may not access this page until you are <a href=\"" + login_url + "\"logged in.</a> ", http.StatusForbidden) // 403
	}
	// Get logout URL
	logout, _ := user.LogoutURL(context, "/")
	// Render page
	nup := NewUserPage{User:uname.String(), Signout:logout}
	renderTemplateFromFile(newUserTemplate, nup, w)
	return
}


/* Create a new user, usually in response to a POST request.
 */
func CreateNewUserHandler(w http.ResponseWriter, r *http.Request) {
	context, uname := verifyLoggedIn(w, r)
	if uname.Id == "" {
		login_url, _ := user.LoginURL(context, r.URL.String())
		http.Error(w, "You may not access this page until you are <a href=\"" + login_url + "\"logged in.</a> ", http.StatusForbidden) // 403
	}
	// Get value of checkbox
	docs := false;
	if r.FormValue("startupdocs") == "docs" {
		docs = true
	} 
	// Create model object
	hu := model.HowlUser{Name:r.FormValue("name"), 
	                     Url:r.FormValue("url"), About:r.FormValue("about"),
	                     DisplayStartupDocs:docs}
	log.Println("Created new user profile for " + r.FormValue("name")) 
	// Make persistant 
	controller.PutUserObject(hu, context, w)
	// Go to homepage
	http.Redirect(w, r, "/", http.StatusFound)
	return
}


/* Create a new data stream, usually in response to a POST request.
 */
func CreateDataStreamHandler(w http.ResponseWriter, r *http.Request) {
	context, uname := verifyLoggedIn(w, r)
	if uname.Id == "" {
		login_url, _ := user.LoginURL(context, r.URL.String())
		http.Error(w, "You may not access this page until you are <a href=\"" + login_url + "\"logged in.</a> ", http.StatusForbidden) // 403
	}
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
	// Make persistent
	controller.PutDataStreamObject(sc, ds, tagnames, context, w)
	// Return to home page
	http.Redirect(w, r, "/", http.StatusFound)
	return
}

