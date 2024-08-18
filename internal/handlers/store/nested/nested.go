package nested

import (
	"sync"

	"github.com/thebeginner86/hippocampus/internal/security"
	"github.com/thebeginner86/hippocampus/persistance/aof"
)

type NestedStoreHandler struct {
	store map[string]map[string]string
	mu    sync.RWMutex

	securityH *security.Security
	aofH 	 *aof.Aof
}

func NewNestedStoreHandler(secH *security.Security, aofH *aof.Aof) *NestedStoreHandler {
	return &NestedStoreHandler{
		store: map[string]map[string]string{},
		mu: sync.RWMutex{},
		securityH: secH,
		aofH: aofH,
	}
}
