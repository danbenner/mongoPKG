package mongodb

import (
	"github.com/globalsign/mgo"
)

// Service ... baseSession is initial session which all other sessions are copies of; queue is a Integer Channel for upto *20; URL is where API service is accessible; Open is current number of copies of session running;
type Service struct {
	baseSession *mgo.Session
	queue       chan int
	URL         string
	Open        int
}

// service ... locally global to the mongodb package
var service Service

// New ... service.New() -> Only used once during startup; Sets service.queue to MaxPool size (20); Creates BASE session via mgo.Dial()
func (s *Service) New() error {
	var err error
	s.queue = make(chan int, MaxPool)
	for i := 0; i < MaxPool; i = i + 1 {
		s.queue <- 1
	}
	s.Open = 0
	s.baseSession, err = mgo.Dial(s.URL)
	return err
}

// Session ... updates queue; returns copy of base session
func (s *Service) Session() *mgo.Session {
	<-s.queue
	s.Open++
	return s.baseSession.Copy()
}

// Close ... update queue; calls Database.Close()
func (s *Service) Close(c *Collection) {
	c.db.s.Close()
	s.queue <- 1
	s.Open--
}
