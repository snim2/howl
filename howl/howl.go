/* This package simply redirects clients to the correct handlers.

Copyright (C) Sarah Mount, 2011.

This program is free software; you can redistribute it and/or
modify it under the terms of the GNU General Public License
as published by the Free Software Foundation; either version 2
of the License, or (at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

package howl

import (
    "http"
)

// Howl packages
import (
	"view"
)

func init() {
    http.HandleFunc("/newstream", view.CreateDataStreamHandler)
    http.HandleFunc("/newuser", view.NewUserHandler)
    http.HandleFunc("/createnewuser", view.CreateNewUserHandler)
    http.HandleFunc("/", view.IndexHandler)
}

