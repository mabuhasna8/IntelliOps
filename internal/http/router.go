package httpapi

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/xid"

	"github.com/mabuhasna8/IntelliOps/internal/run"
	"github.com/mabuhasna8/IntelliOps/internal/workflow"
)

type Server struct {
	WFStore  *workflow.Store
	RunStore *run.Store
}

func NewServer(wfs *workflow.Store, rs *run.Store) *Server {
	return &Server{WFStore: wfs, RunStore: rs}
}

func (s *Server) Router() http.Handler {
	r := mux.NewRouter()

	r.HandleFunc("/api/v1/workflows", s.CreateWorkflow).Methods("POST")
	r.HandleFunc("/api/v1/workflows/{id}", s.GetWorkflow).Methods("GET")

	r.HandleFunc("/api/v1/runs", s.CreateRun).Methods("POST")
	r.HandleFunc("/api/v1/runs/{id}", s.GetRun).Methods("GET")

	return r
}

type createWorkflowReq struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Version  string `json:"version"`
	SpecYAML string `json:"spec_yaml"`
}

func (s *Server) CreateWorkflow(w http.ResponseWriter, r *http.Request) {
	var req createWorkflowReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	wf := &workflow.Workflow{
		ID:       req.ID,
		Name:     req.Name,
		Version:  req.Version,
		SpecYAML: req.SpecYAML,
	}
	s.WFStore.Create(wf)

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(wf)
}

func (s *Server) GetWorkflow(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	wf, ok := s.WFStore.Get(id)
	if !ok {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(wf)
}

type createRunReq struct {
	WorkflowID string                 `json:"workflow_id"`
	Version    string                 `json:"version"`
	EnvID      string                 `json:"env_id"`
	Params     map[string]interface{} `json:"params"`
}

func (s *Server) CreateRun(w http.ResponseWriter, r *http.Request) {
	var req createRunReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	id := xid.New().String()
	runObj := &run.Run{
		ID:         id,
		WorkflowID: req.WorkflowID,
		Version:    req.Version,
		EnvID:      req.EnvID,
		Status:     run.StatusPending,
		Params:     req.Params,
	}
	s.RunStore.Create(runObj)

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(runObj)
}

func (s *Server) GetRun(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	obj, ok := s.RunStore.Get(id)
	if !ok {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(obj)
}
