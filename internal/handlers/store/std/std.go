package std

import (
	"sync"

	"github.com/thebeginner86/hippocampus/internal/security"
	"github.com/thebeginner86/hippocampus/persistance/aof"
)

type StdStoreHandler struct {
	store map[string]string
	mu    sync.RWMutex

	securityH *security.Security
	aofH 	 *aof.Aof
}

func NewStdStoreHandler(secH *security.Security, aofH *aof.Aof) *StdStoreHandler {
	return &StdStoreHandler{
		store: map[string]string{},
		mu: sync.RWMutex{},
		securityH: secH,
		aofH: aofH,
	}
}

