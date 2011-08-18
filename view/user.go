/* Handler for dealing with URIs starting with /user/
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
	"http"
	"log"
	"os"
	"strconv"
	"strings"
)


import (
	"model"
	"controller"
)


/* Struct to store data for the profileTemplate template.
 */
type ProfilePage struct {
	User		 string // User name to use in top bar
	Signout		 string // Url to log out
	Uid          string // User ID (if one has been created)
	Name         string // Real name
	Url          string // Profile URL
	Description  string // Profile description
	Action       string // Form action either /user/{Uid}/edit, or /user/new
	Button       string // Label on button, usually Create or Update
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
	log.Println("UserHandler got request with method: " + r.Method + " and path: " + r.URL.Path)
	paths := strings.Split(r.URL.Path, "/", -1)
	log.Println("UserHandler Paths has " + strconv.Itoa(len(paths)) + " items")
	if len(paths) < 2 {
		serve404(w)
		return
	}
	context := appengine.NewContext(r)
	g_user := user.Current(context)
	logged_in_user, _ := controller.GetUserFromEmail(context, g_user.Email)
	switch r.Method  {
	case "POST":
		// User is logged in but has yet to create a profile object
		if logged_in_user == nil {
			hu := model.HowlUser{Name:r.FormValue("name"), 
		                         Uid:r.FormValue("uid"),
	                             Url:r.FormValue("url"), 
	                             About:r.FormValue("about")}
			_, err := hu.Create(context)	
			log.Println("Created new user profile for " + r.FormValue("uid")) 
			if err != nil {
				serveError(context, w, http.StatusInternalServerError, err)
				return
			}
			req, _ := http.NewRequest("GET", "/", r.Body)
			http.Redirect(w, req, "/", http.StatusFound)			
			return
		}
		// PRE: Can only edit username if you are logged in as that user
		if r.FormValue("uid") != logged_in_user.Uid {
			serveError(context, w, http.StatusForbidden, os.NewError("Forbidden: You can only edit your own user profile"))
			return
		}
		// User is logged in and editing existing profile.
		logged_in_user.Name = r.FormValue("name")
		logged_in_user.Uid = r.FormValue("uid")
		logged_in_user.Url =  r.FormValue("url")
		logged_in_user.About = r.FormValue("about")
		err := logged_in_user.Update(context)
		if err != nil {
			serveError(context, w, http.StatusInternalServerError, err)
			return
		}
		req, _ := http.NewRequest("GET", "/", r.Body)
		http.Redirect(w, req, "/", http.StatusFound)
		return
	case "GET":
		logout, _ := user.LogoutURL(context, "/")
		npp := new(ProfilePage)
		if len(paths) > 2 { 
			// User already exists and has a profile object
			uname := paths[2]
			log.Println("Username: " + uname)
			userobj, _ := controller.GetUserFromUid(context, uname)
			npp.User			= logged_in_user.Uid
			npp.Signout			= logout
			npp.Uid				= userobj.Uid
			npp.Name			= userobj.Name
			npp.Url				= userobj.Url
			npp.Description		= userobj.About
			npp.Action          = "/user/" + userobj.Uid + "/edit"
			npp.Button          = "Update"
		} else { 
			// First time user, needs to create profile
			log.Println("First time user.")
			npp.User = g_user.String()
			npp.Action = "/user/new"
			npp.Button = "Create"
			npp.Signout = logout
		}
		renderTemplateFromFile(context, profileTemplate, npp, w)
		return
	default:
		err := os.NewError("Could not understand your request: " + r.URL.Path)
		serveError(context, w, http.StatusBadRequest, err)
		return
	}
}

