package picker

import (
	"sync"

	"github.com/reubenmiller/c8y-session-bitwarden/pkg/core"
)

type randomItemGenerator struct {
	sessions []*core.CumulocitySession
	index    int
	mtx      *sync.Mutex
}

func (r *randomItemGenerator) Len() int {
	return len(r.sessions)
}

func (r *randomItemGenerator) reset() {
	r.mtx = &sync.Mutex{}
}

func (r *randomItemGenerator) Next() *core.CumulocitySession {
	if r.mtx == nil {
		r.reset()
	}

	r.mtx.Lock()
	defer r.mtx.Unlock()

	i := core.CloneSession(r.sessions[r.index])

	r.index++
	if r.index >= len(r.sessions) {
		r.index = 0
	}

	return i
}
