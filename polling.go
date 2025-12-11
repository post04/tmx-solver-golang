package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/obfio/tmx-solver-golang/mongo"
	allsites "github.com/obfio/tmx-solver-golang/sites/allSites"
)

type pollingTask struct {
	ID        string             `json:"id"`
	URL       string             `json:"-"`
	Site      string             `json:"-"`
	UUID      string             `json:"-"`
	Proxy     string             `json:"-"`
	Mobile    bool               `json:"-"`
	Status    string             `json:"status"`
	CreatedAt time.Time          `json:"createdAt"`
	StopAt    int                `json:"stopAt"`
	Solution  *allsites.Response `json:"solution,omitempty"`
}

var (
	currTasks  = make(map[string]*pollingTask)
	taskLocker = sync.RWMutex{}
	newTasks   = make(chan *pollingTask)
)

func handlePollingNew(w http.ResponseWriter, r *http.Request) {
	// Set JSON content type
	w.Header().Set("Content-Type", "application/json")

	// Only allow POST requests
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse the request body
	var req *allsites.ProxyRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&req); err != nil {
		response := allsites.Response{
			Error: "Invalid request: " + err.Error(),
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	// validate API key
	if !mongo.ValidateAPIKey(req.APIKey) {
		response := allsites.Response{
			Error: "Invalid API key",
		}
		json.NewEncoder(w).Encode(response)
		return
	}
	if !req.Mobile && !sliceContains(possibleSites, req.Site) {
		response := allsites.Response{
			Error: "Invalid site: " + req.Site + ". Must be one of: " + strings.Join(possibleSites, ", "),
		}
		json.NewEncoder(w).Encode(response)
		return
	}
	if req.Mobile && !sliceContains(possibleSitesMobile, req.Site) {
		response := allsites.Response{
			Error: "Invalid site: " + req.Site + ". Must be one of: " + strings.Join(possibleSitesMobile, ", "),
		}
		json.NewEncoder(w).Encode(response)
		return
	}
	if req.Mobile {
		req.Site += "MOBILE"
	}
	task := &pollingTask{
		ID:        allsites.RandStringBytesMaskImprSrc1(11),
		URL:       req.URL,
		Site:      req.Site,
		UUID:      req.SessionID,
		Proxy:     req.Proxy,
		Mobile:    req.Mobile,
		Status:    "pending",
		CreatedAt: time.Now(),
		StopAt:    req.StopAt,
	}
	taskLocker.Lock()
	currTasks[task.ID] = task
	newTasks <- task
	taskLocker.Unlock()
	json.NewEncoder(w).Encode(task)
}

type pollingStatusRequest struct {
	APIKey string `json:"apiKey"`
	TaskID string `json:"taskID"`
}

func handlePollingStatus(w http.ResponseWriter, r *http.Request) {
	// Set JSON content type
	w.Header().Set("Content-Type", "application/json")

	// Only allow POST requests
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse the request body
	var req pollingStatusRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&req); err != nil {
		response := allsites.Response{
			Error: "Invalid request: " + err.Error(),
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	// validate API key
	if !mongo.ValidateAPIKey(req.APIKey) {
		response := allsites.Response{
			Error: "Invalid API key",
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	// get task
	taskLocker.RLock()
	task, ok := currTasks[req.TaskID]
	taskLocker.RUnlock()
	if !ok {
		response := allsites.Response{
			Error: "Task not found",
		}
		json.NewEncoder(w).Encode(response)
		return
	}
	if task.Status == "error" || task.Status == "success" {
		if task.Status == "success" {
			mongo.UpdateUsesCount(req.APIKey)
		}
		// remove task from currTasks
		taskLocker.Lock()
		delete(currTasks, req.TaskID)
		taskLocker.Unlock()
	}
	// return task status
	json.NewEncoder(w).Encode(task)
}

func handleTasks() {
	go func() {
		for {
			task := <-newTasks
			go func(t *pollingTask) {
				response := allsites.GetCookies(&allsites.ProxyRequest{
					Proxy:     t.Proxy,
					Site:      t.Site,
					URL:       t.URL,
					SessionID: t.UUID,
					Mobile:    t.Mobile,
					StopAt:    t.StopAt,
				})
				if response.Error != "" {
					t.Status = "error"
					t.Solution = response
					return
				}
				t.Status = "success"
				t.Solution = response
			}(task)
		}
	}()

	// clear tasks that are older than 10 minutes
	go func() {
		for {
			time.Sleep(10 * time.Minute)
			taskLocker.Lock()
			for id, task := range currTasks {
				if task.CreatedAt.Before(time.Now().Add(-10 * time.Minute)) {
					delete(currTasks, id)
				}
			}
			taskLocker.Unlock()
		}
	}()
}
