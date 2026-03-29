//////////////////////////////////////////////////////////////////////
//
// Your video processing service has a freemium model. Everyone has 10
// sec of free processing time on your service. After that, the
// service will kill your process, unless you are a paid premium user.
//
// Beginner Level: 10s max per request
// Advanced Level: 10s max per user (accumulated)
//

package main

import (
	"time"
	"sync"
)

const MAX_TIME_ALLOWED = 10

// User defines the UserModel. Use this to check whether a User is a
// Premium user or not
type User struct {
	ID        int
	IsPremium bool
	TimeUsed  int64 // in seconds
	MULock 	  sync.Mutex
}

func (u *User) AddTime(seconds int64) (int64) {
	u.MULock.Lock()
	defer u.MULock.Unlock()
	u.TimeUsed += seconds
	return u.TimeUsed

}

// HandleRequest runs the processes requested by users. Returns false
// if process had to be killed
func HandleRequest(process func(), u *User) bool {
	if u.IsPremium {
		process()
		return true
	}

	u.MULock.Lock()
	if u.TimeUsed >= MAX_TIME_ALLOWED && !u.IsPremium {
		u.MULock.Unlock()
		return false
	}
	u.MULock.Unlock()
	
	done := make(chan bool)
	ticker := time.Tick(time.Second)

	go func() {
		process()
		done <- true
	} ()

	
	for {
		select {
			case <- done:
				return true
			case <- ticker:
				if currTime := u.AddTime(1); currTime > MAX_TIME_ALLOWED {
					return false
				}
		}
	}

}

func main() {
	RunMockServer()
}
