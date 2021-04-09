package handler

import (
	"MangaLibrary/src/internal/api"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/mux"

	"go.uber.org/zap"
	"golang.org/x/net/websocket"
)

type SocketHandler struct {
	Logger *zap.Logger
	Query  *http.Request
}

// Heavily based on Kubernetes' (https://github.com/GoogleCloudPlatform/kubernetes) detection code.
var connectionUpgradeRegex = regexp.MustCompile("(^|.*,\\s*)upgrade($|\\s*,)")

func isWebsocketRequest(req *http.Request) bool {
	return connectionUpgradeRegex.MatchString(strings.ToLower(req.Header.Get("Connection"))) && strings.ToLower(req.Header.Get("Upgrade")) == "websocket"
}

func (h SocketHandler) Handle(w http.ResponseWriter, r *http.Request) {
	// Handle websockets if specified.
	h.Query = r
	if isWebsocketRequest(r) {
		websocket.Handler(h.HandleWebSockets).ServeHTTP(w, r)
	} else {
		h.HandleHttp(w, r)
	}

}

func (h SocketHandler) HandleWebSockets(ws *websocket.Conn) {
	dataMap := map[string]string{}
	for k, v := range h.Query.URL.Query() {
		dataMap[k] = v[0]
	}
	siteName := mux.Vars(h.Query)["site_name"]
	Lib := api.NewLibrary()
	Lib.LoadDefaults(h.Logger)
	p := &api.Progress{}
	var wg sync.WaitGroup
	wg.Add(2)
	finished := false
	go func() {
		err := Lib.ProcessSite(siteName, dataMap, p)
		if err != nil {
			//SendError(err.Error(), ws.Resp)
			return
		}
		finished = true
		wg.Done()

	}()
	go func() {
		for !finished {
			h.Logger.Info(fmt.Sprintf("Sending some data: %d/%d", Lib.Progress.Current, Lib.Progress.Total))
			err := websocket.JSON.Send(ws, Lib.Progress)
			if err != nil {
				h.Logger.Error(fmt.Sprintf("Sending some data:"), zap.Error(err))
				return
			}
			// Artificially induce a 1s pause
			time.Sleep(time.Second * 5)
		}
		wg.Done()
	}()
	wg.Wait()

}

func (h SocketHandler) HandleHttp(w http.ResponseWriter, r *http.Request) {
	cn, ok := w.(http.CloseNotifier)
	if !ok {
		http.NotFound(w, r)
		return
	}
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.NotFound(w, r)
		return
	}

	// Send the initial headers saying we're gonna stream the response.
	w.Header().Set("Transfer-Encoding", "chunked")
	w.WriteHeader(http.StatusOK)
	flusher.Flush()

	enc := json.NewEncoder(w)

	dataMap := map[string]string{}
	for k, v := range h.Query.URL.Query() {
		dataMap[k] = v[0]
	}
	siteName := mux.Vars(h.Query)["site_name"]
	Lib := api.NewLibrary()
	Lib.LoadDefaults(h.Logger)

	p := &api.Progress{}
	var wg sync.WaitGroup
	wg.Add(2)
	finished := false
	go func() {
		err := Lib.ProcessSite(siteName, dataMap, p)
		if err != nil {
			//SendError(err.Error(), ws.Resp)
			return
		}
		finished = true
		wg.Done()

	}()
	go func() {
		for !finished {
			//h.Logger.Info(fmt.Sprintf("Sending some data: %d", i))
			// Send some data.
			select {
			case <-cn.CloseNotify():
				return
			default:
				h.Logger.Info(fmt.Sprintf("Sending some data: %d/%d", Lib.Progress.Current, Lib.Progress.Total))
				err := enc.Encode(p)
				if err != nil {
					h.Logger.Error(fmt.Sprintf("Sending some data:"), zap.Error(err))
					break
				}
				flusher.Flush()
				// Artificially induce a 1s pause
				time.Sleep(time.Second * 5)
			}
		}
		wg.Done()
	}()
	wg.Wait()
}
