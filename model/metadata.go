/* Model of metadata (e.g. tags, comments, annotations) stored in Howl.
 *
 * TODO: Use memcache
 * TODO: Document how keys are generated and used.
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

package model


import (
	"appengine"
	"appengine/datastore"
	"log"
	"os"
	"strings"
)


/* A annotation describes a datum.
 */
type Annotation struct {
	Text		 string
}


func (annotation *Annotation) Create(context appengine.Context) (*datastore.Key, os.Error) {
	log.Println("Creating new annotation")
	key := datastore.NewIncompleteKey("Annotation")
	err_s := "Error storing new annotation "
	return put(context, key, err_s, annotation)
}


func (*Annotation) Query(context appengine.Context, dsquery *datastore.Query, annotations *[]Annotation) ([]*datastore.Key, os.Error) {
	err_s := "Error fetching Annotation object: "
	keys, err := query(context, dsquery, annotations, err_s)
	return keys, err
}


func (annotation *Annotation) String() (string) {
	return annotation.Text
}


/* A comment is a message from a (usually human) Howl user.
 */
type Comment struct {
	Text		 string
	Author		 string // HowlUser.Uid
}


func (comment *Comment) Create(context appengine.Context) (*datastore.Key, os.Error) {
	log.Println("Creating new comment")
	key := datastore.NewIncompleteKey("Comment")
	err_s := "Error storing new comment "
	return put(context, key, err_s, comment)
}


func (*Comment) Query(context appengine.Context, dsquery *datastore.Query, comments *[]Comment) ([]*datastore.Key, os.Error) {
	err_s := "Error fetching Comment object: "
	keys, err := query(context, dsquery, comments, err_s)
	return keys, err
}


func (comment *Comment) String() (string) {
	return comment.Text
}


/* A Tag is a user-defined piece of metadata.
 *
 * Tags must always be Singletons in the datastore.
 */
type Tag struct {
	Tag string
}


func MakeTags(context appengine.Context, tagnames []string) ([]Tag, []*datastore.Key, os.Error) {
	tags := make([]Tag, len(tagnames)) 
	keys := make([]*datastore.Key, len(tagnames)) 
	errs := make([]os.Error, len(tagnames))
	for i := 0; i < len(tagnames); i++ {
		tags[i] = Tag{strings.Trim(tagnames[i], " ")}
		keys[i], errs[i] = tags[i].Create(context)
	}
	for i := 0; i < len(tagnames); i++ {		
		if errs[i] != nil {
			return tags, keys, errs[i]
		}
	}
	return tags, keys, nil
}


func TagsToStrings(tags []Tag) []string {
	tagnames := make([]string, len(tags))
	for i := 0; i < len(tagnames); i++ {
		tagnames[i] = tags[i].String()
	}
	return tagnames
}


func (tag *Tag) Create(context appengine.Context) (*datastore.Key, os.Error) {
	log.Println("Creating new tag " + tag.Tag)
	key := datastore.NewKey("Tag", tag.Tag, 0, nil)
	err_s := "Error storing new tag " + tag.Tag
	return putSingleton(context, key, err_s, tag)
}


func (*Tag) Query(context appengine.Context, dsquery *datastore.Query, tags *[]Tag) ([]*datastore.Key, os.Error) {
	err_s := "Error fetching Tag object: "
	keys, err := query(context, dsquery, tags, err_s)
	return keys, err
}


func (tag *Tag) Update(context appengine.Context) (os.Error) {
	key := datastore.NewKey("Tag", tag.Tag, 0, nil) 
	_, err := putSingleton(context, key, "Error storing new Tag: ", tag.Tag)
	return err
}


func (tag *Tag) Delete(context appengine.Context) (os.Error) {
	key := datastore.NewKey("Tag", tag.Tag, 0, nil) 
	return delete(context, key)
}


func (tag *Tag) String() (string) {
	return tag.Tag
}

