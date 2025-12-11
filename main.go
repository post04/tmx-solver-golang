package main

import (
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/obfio/tmx-solver-golang/mongo"
	allsites "github.com/obfio/tmx-solver-golang/sites/allSites"
	"github.com/obfio/tmx-solver-golang/tmx/websocket"
)

var (
	adminKey = "ADMIN KEY (should be replaced by config.json)"
)

type config struct {
	AdminKey string `json:"adminKey"`
}

func init() {
	rand.Seed(time.Now().UnixNano())
	f, err := os.ReadFile("./config.json")
	if err != nil {
		panic(err)
	}
	var c config
	err = json.Unmarshal(f, &c)
	if err != nil {
		panic(err)
	}
	adminKey = c.AdminKey
	// wipe ./errorsWithScripts
	os.RemoveAll("./errorsWithScripts")
	os.MkdirAll("./errorsWithScripts", 0755)

	for name := range allsites.Sites {
		possibleSites = append(possibleSites, name)
	}
}

var (
	possibleSites       = []string{}
	possibleSitesMobile = []string{}
)

// Handler function for the POST endpoint
func handleProxyRequest(w http.ResponseWriter, r *http.Request) {
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
	response := allsites.GetCookies(req)
	if response.Error != "" {
		json.NewEncoder(w).Encode(response)
		return
	}
	json.NewEncoder(w).Encode(response)
	fmt.Printf("Successfully processed request for proxy (2): %s\n", req.Proxy)
	mongo.UpdateUsesCount(req.APIKey)
	return
}

type newAPIKeyRequest struct {
	AdminKey string `json:"adminKey"`
	Uses     int64  `json:"uses"`
	Duration int64  `json:"duration"`
	Name     string `json:"name"`
}

func handleNewAPIKey(w http.ResponseWriter, r *http.Request) {
	// Set JSON content type
	w.Header().Set("Content-Type", "application/json")
	// Only allow POST requests
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req newAPIKeyRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&req); err != nil {
		response := allsites.Response{
			Error: "Invalid request: " + err.Error(),
		}
		json.NewEncoder(w).Encode(response)
		return
	}
	if req.AdminKey != adminKey {
		response := allsites.Response{
			Error: "Invalid admin key",
		}
		json.NewEncoder(w).Encode(response)
		return
	}
	if req.Duration <= 0 {
		req.Duration = 7
	}
	if req.Uses <= 0 {
		req.Uses = math.MaxInt64
	}
	if req.Name == "" {
		req.Name = "no name defined"
	}
	f := &mongo.Object{}
	f.Key = allsites.RandStringBytesMaskImprSrc1(20)
	f.CreatedAt = time.Now().UnixMilli()
	f.ExpiresAt = time.Now().Add(time.Hour * 24 * time.Duration(req.Duration)).UnixMilli()
	f.MaxUses = req.Uses
	f.Name = req.Name
	err := mongo.AddUser(f)
	if err != nil {
		response := allsites.Response{
			Error: "Failed to add user: " + err.Error(),
		}
		json.NewEncoder(w).Encode(response)
		return
	}
	w.Write([]byte(f.Key))
	//w.WriteHeader(200)
	return
}

type editAPIKeyRequest struct {
	AdminKey string  `json:"adminKey"`
	Key      string  `json:"key"`
	Uses     *int64  `json:"uses"`
	Duration *int64  `json:"duration"`
	Name     *string `json:"name"`
	Action   string  `json:"action"`
}

func handleEditAPIKey(w http.ResponseWriter, r *http.Request) {
	// Set JSON content type
	w.Header().Set("Content-Type", "application/json")
	// Only allow POST requests
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req editAPIKeyRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&req); err != nil {
		response := allsites.Response{
			Error: "Invalid request: " + err.Error(),
		}
		json.NewEncoder(w).Encode(response)
		return
	}
	if req.AdminKey != adminKey {
		response := allsites.Response{
			Error: "Invalid admin key",
		}
		json.NewEncoder(w).Encode(response)
		return
	}
	if !mongo.ValidateAPIKeyWithoutUses(req.Key) {
		response := allsites.Response{
			Error: "Invalid API key",
		}
		json.NewEncoder(w).Encode(response)
		return
	}
	switch req.Action {
	case "delete":
		mongo.DeleteUser(req.Key)
		break
	case "editTime":
		if req.Duration == nil {
			response := allsites.Response{}
			response.Error = "Invalid duration"
			json.NewEncoder(w).Encode(response)
			return
		}
		updatedTime := time.Now().Add(time.Hour * 24 * time.Duration(*req.Duration)).UnixMilli()
		if updatedTime <= 0 {
			updatedTime = math.MaxInt64
		}
		mongo.UpdateTime(req.Key, updatedTime)
		break
	case "editUses":
		if req.Uses == nil {
			response := allsites.Response{}
			response.Error = "Invalid uses"
			json.NewEncoder(w).Encode(response)
			return
		}
		mongo.UpdateMaxUses(req.Key, *req.Uses)
		break
	case "editName":
		if req.Name == nil {
			response := allsites.Response{}
			response.Error = "Invalid name"
			json.NewEncoder(w).Encode(response)
			return
		}
		mongo.UpdateName(req.Key, *req.Name)
		break
	}
	w.Write([]byte(req.Key))
}

func main() {
	go handleTasks()
	go websocket.HandleTasks()
	// Setup POST HTTP API that takes in a proxy and returns user-agent + cookie
	http.HandleFunc("/api/get-cookies", handleProxyRequest)
	http.HandleFunc("/api/get-cookies-polling", handlePollingNew)
	http.HandleFunc("/api/get-cookies-polling-status", handlePollingStatus)
	// Setup POST HTTP API that takes in an admin key and returns a new API key
	http.HandleFunc("/admin/new", handleNewAPIKey)

	http.HandleFunc("/admin/edit", handleEditAPIKey)

	// websocket server
	http.HandleFunc("/ws", websocket.HandleNewConnection)

	port := ":8081"
	fmt.Printf("Starting server on port %s...\n", port)

	// Start the HTTP server
	if err := http.ListenAndServe(port, nil); err != nil {
		fmt.Printf("Server failed to start: %v\n", err)
		os.Exit(1)
	}
}
