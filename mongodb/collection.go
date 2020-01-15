package mongodb

import (
	"log"

	"github.com/globalsign/mgo"
)

//Collection ..
type Collection struct {
	db      *Database
	name    string
	Session *mgo.Collection
}

//Connect ...
func (c *Collection) Connect() {
	session := *c.db.session.C(c.name)
	c.Session = &session
}

//NewCollectionSession ...
func NewCollectionSession(name string) *Collection {
	var c = Collection{
		db:   newDBSession(DBNAME),
		name: name,
	}
	c.Connect()
	return &c
}

// NewCollectionWithIndex ... creates copy of original session with a Collection Name and ANY (including zero) index strings (to limit search);
// NOTE: 7 is the MAX number of index values, otherwise excedes 127 byte max limit for 'EnsureIndex'
func NewCollectionWithIndex(cName string, indexStrings ...string) (c *Collection) {
	c = NewCollectionSession(cName)
	// []string{`status`, `retryScheduledAt`},
	err := c.Session.EnsureIndex(mgo.Index{Key: indexStrings})
	if err != nil {
		if err.Error() != "invalid index key: no fields provided" {
			log.Println("EnsureIndex Error: " + err.Error())
		}
	}
	return
}

//Close ..
func (c *Collection) Close() {
	service.Close(c)
}
