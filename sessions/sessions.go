package sessions

import (
	"github.com/gorilla/sessions" // go get github.com/gorilla/sessions
)

var Store = sessions.NewCookieStore( []byte("t0p-S3cr3t") )