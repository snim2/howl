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
 */
func init() {
    http.HandleFunc("/newstream", view.CreateDataStreamHandler)  // Should be /stream?action=create or PUT
    http.HandleFunc("/newuser", view.NewUserHandler)
    http.HandleFunc("/createnewuser", view.CreateNewUserHandler) // Should be /user?action=create or PUT

	// Especially for web browsers
    http.HandleFunc("/dashboard", view.DashboardHandler)
    http.HandleFunc("/", view.IndexHandler)

    // http.HandleFunc("/user", view.UserHandler)
    // http.HandleFunc("/stream", view.StreamHandler)
    // http.HandleFunc("/provider", view.ProviderHandler)
    // http.HandleFunc("/datum", view.DatumHandler)
}

/*
From StackOverflow

application = webapp.WSGIApplication([
    ('/user/([^/]+)/([^/]+)', UserHandler),
    ], debug=True)

class UserHandler(webapp.RequestHandler):
  def get(self, user_id, action_to_consume):
    self.response.out.write("Action %s" % action_to_consume)#Should print History
*/