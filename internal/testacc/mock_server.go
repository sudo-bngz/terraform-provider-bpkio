package testacc

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"sync"
)

type adServer struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type transcodingProfile struct {
	ID         uint   `json:"id"`
	Name       string `json:"name"`
	Content    string `json:"content"`
	InternalID string `json:"internal_id"`
}

/* -------------------------------------------------------------------------- */
/*  Mock Broadpeak IO backend                                                 */
/* -------------------------------------------------------------------------- */
type Mock struct {
	srv                 *httptest.Server
	mu                  sync.RWMutex
	adServers           map[uint]adServer
	transcodingProfiles map[uint]transcodingProfile
	nextID              uint
	validToken          string
}

/* -------------------------- constructor & helpers ------------------------- */

func NewMock() *Mock {
	m := &Mock{
		adServers:           make(map[uint]adServer),
		transcodingProfiles: make(map[uint]transcodingProfile),
		nextID:              1,
		validToken:          "test-token",
	}

	m.transcodingProfiles[1] = transcodingProfile{
		ID:   1,
		Name: "Default Offline Profile",
		// Come from the real default
		Content:    `{"type":"OFFLINE_TRANSCODING","audios":{"common":{"loudnorm":{"i":-23,"tp":-1},"sampling_rate":48000},"audio_0":{"bitrate":128000,"codec_string":"mp4a.40.2","channel_layout":"stereo"}},"videos":{"common":{"gop_size":25,"framerate":{"den":1000,"num":25000},"perf_level":3},"video_0":{"scale":{"width":-2,"height":232},"bitrate":500000,"codec_string":"avc1.42C00D"},"video_1":{"scale":{"width":-2,"height":360},"bitrate":1600000,"codec_string":"avc1.4D401E"},"video_2":{"scale":{"width":-2,"height":480},"bitrate":2300000,"codec_string":"avc1.4D401F"},"video_3":{"scale":{"width":-2,"height":720},"bitrate":3200000,"codec_string":"avc1.640020"},"video_4":{"scale":{"width":-2,"height":1080},"bitrate":5000000,"codec_string":"avc1.640028"}},"version":"02.00.05","packaging":{"hls":{"version":3,"fragment_length":{"den":1,"num":4}},"dash":{"fragment_length":{"den":1,"num":4}}}}`,
		InternalID: "default-offline-profile",
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/v1/adservers", m.handleCollection)
	mux.HandleFunc("/v1/adservers/", m.handleItem)
	mux.HandleFunc("/v1/transcoding-profiles", m.handleTranscodingProfiles)

	m.srv = httptest.NewServer(mux)
	return m
}

func (m *Mock) URL() string   { return m.srv.URL }
func (m *Mock) Token() string { return m.validToken }
func (m *Mock) Close()        { m.srv.Close() }

/* ------------------------ convenience for acceptance --------------------- */

func (m *Mock) ResourceExists(id string) bool {
	uid, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return false
	}
	m.mu.RLock()
	_, ok := m.adServers[uint(uid)]
	m.mu.RUnlock()
	return ok
}

/* -------------------------------------------------------------------------- */
/*  Authorization + Handlers                                                  */
/* -------------------------------------------------------------------------- */

func (m *Mock) authorize(w http.ResponseWriter, r *http.Request) bool {
	log.Println(">>> NEW handler reached")
	if r.Header.Get("Authorization") != "Bearer "+m.validToken {
		fmt.Println("Authorization header:", r.Header.Get("Authorization"))
		http.Error(w,
			`{"message":"token content not correct","error":"Forbidden","statusCode":403}`,
			http.StatusForbidden)
		return false
	}
	return true
}

func (m *Mock) handleCollection(w http.ResponseWriter, r *http.Request) {
	log.Println(">>> NEW handler reached")
	if !m.authorize(w, r) {
		return
	}

	switch r.Method {
	case http.MethodPost:
		var in struct {
			Name string `json:"name"`
		}
		if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		m.mu.Lock()
		id := m.nextID
		m.nextID++
		obj := adServer{ID: id, Name: in.Name}
		m.adServers[id] = obj
		m.mu.Unlock()

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(obj)

	case http.MethodGet:
		m.mu.RLock()
		var list []adServer
		for _, v := range m.adServers {
			list = append(list, v)
		}
		m.mu.RUnlock()

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(list)

	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (m *Mock) handleItem(w http.ResponseWriter, r *http.Request) {
	if !m.authorize(w, r) {
		return
	}

	idStr := r.URL.Path[len("/v1/adservers/"):]
	uid64, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	uid := uint(uid64)

	switch r.Method {
	case http.MethodGet:
		m.mu.RLock()
		obj, ok := m.adServers[uid]
		m.mu.RUnlock()
		if !ok {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(obj)

	case http.MethodDelete:
		m.mu.Lock()
		delete(m.adServers, uid)
		m.mu.Unlock()
		w.WriteHeader(http.StatusNoContent)

	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (m *Mock) handleTranscodingProfiles(w http.ResponseWriter, r *http.Request) {
	if !m.authorize(w, r) {
		return
	}

	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	m.mu.RLock()
	defer m.mu.RUnlock()
	var list []transcodingProfile
	for _, v := range m.transcodingProfiles {
		list = append(list, v)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(list)
}
