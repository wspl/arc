package arc

import (
	"fmt"
	"math/rand"

	"github.com/orcaman/concurrent-map"
)

func createSessionSet() (*SessionSet, error) {
	s := new(SessionSet)
	s.sessions = cmap.New()
	return s, nil
}

type SessionSet struct {
	sessions cmap.ConcurrentMap
}

func (s *SessionSet) Has(id uint32) bool {
	return s.sessions.Has(fmt.Sprint(id))
}

func (s *SessionSet) Set(id uint32, session *Session) {
	s.sessions.Set(fmt.Sprint(id), session)
}

func (s *SessionSet) Get(id uint32) *Session {
	v, _ := s.sessions.Get(fmt.Sprint(id))
	return v.(*Session)
}

func (s *SessionSet) New(conn *ArcConn) (uint32, *Session) {
	id := rand.Uint32()
	if !s.Has(id) && id != 0 {
		session, _ := createSession(conn, id)
		s.Set(id, session)
		return id, session
	} else {
		return s.New(conn)
	}
}