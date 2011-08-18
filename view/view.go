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
			UserHandler(w, r) // In file user.go
			return
		}
	case "stream":
		switch responseTy {
			default: 
			StreamHandler(w, r) // In file stream.go
			return
		}
	case "provider":
		switch responseTy {
			default: 
			ProviderHandler(w, r) // In file provider.go
			return
		}
	case "datum":
		switch responseTy {
			default: 
			DatumHandler(w, r) // In file datum.go
			return
		}
	}
	serve404(w)
	return
}

