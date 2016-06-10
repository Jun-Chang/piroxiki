package piroxiki

import (
	"bytes"
	"fmt"
	"net/http"
	"sync"

	"github.com/julienschmidt/httprouter"
)

type Http struct {
	Method     string
	URL        string
	PORT       string
	Body       string
	BodyFormat string
	//filter   []imap.Field TODO
	message chan Message
	errc    chan<- error
}

var httpm sync.RWMutex
var httpz []*Http = make([]*Http, 0)

func (h *Http) HandleInput(policy Policy, filter Filter, errc chan<- error) (<-chan Message, error) {
	h.errc = errc
	h.message = make(chan Message)

	httpInputListen()

	httpm.Lock()
	defer httpm.Unlock()
	httpz = append(httpz, h)

	return h.message, nil
}

func (*Http) HandleOutput(policy Policy, message <-chan Message, errc chan<- error) error {
	go func() {
		for {
			msg := <-message
			fmt.Println(msg)
		}
	}()

	return nil
}

var listenOnce sync.Once

func httpInputListen() {
	all := "/*all"
	listenOnce.Do(func() {
		router := httprouter.New()
		router.GET(all, allHandler)
		router.POST(all, allHandler)
		router.DELETE(all, allHandler)
		router.PATCH(all, allHandler)
		router.PUT(all, allHandler)
		go http.ListenAndServe(":8080", router)
	})
}

func allHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	fmt.Println(r.URL.String())
	fmt.Println(r.Method)
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	fmt.Println(buf.String())

	httpm.RLock()
	defer httpm.RUnlock()

	for _, h := range httpz {
		// TODO regexp
		if h.URL == r.URL.String() {
			h.message <- map[string]interface{}{"url": r.URL.String()}
			return
		}
	}
	fmt.Println("not much", r.URL.String())
}
