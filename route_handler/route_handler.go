package route_handler

import (
	"context"
	"encoding/json"
	"net/http"

	"cloud.google.com/go/firestore"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"google.golang.org/api/iterator"
)

// Handler for HTTP requests
type Handler struct {
	logger   *zap.SugaredLogger
	router   *mux.Router
	database *firestore.Client
}

// New creates a new Handler
func New(logger *zap.SugaredLogger, router *mux.Router, database *firestore.Client) *Handler {
	h := Handler{logger, router, database}
	h.registerRoutes()

	return &h
}

// RegisterRoutes for all http endpoints
func (h *Handler) registerRoutes() {
	h.router.HandleFunc("/beans", h.getBeans).Methods("GET")
	h.router.HandleFunc("/beans", h.addBean).Methods("POST")
}

// Bean represents a coffee bean
type Bean struct {
	Flavors []string `json:"flavors"`
	Name    string   `json:"name"`
	Roaster string   `json:"roaster"`
	Shade   string   `json:"shade"`
}

// BeansResp is the response for the beans endpoint
type BeansResp struct {
	Beans []Bean `json:"beans"`
}

// AddBeanReq is the request body for adding a Bean
type AddBeanReq struct {
	Flavors []string `json:"flavors"`
	Name    string   `json:"name"`
	Roaster string   `json:"roaster"`
	Shade   string   `json:"shade"`
}

// AddBeanResp is the response from the POST /beans endpoint
type AddBeanResp struct {
	ID string `json:"id"`
}

func (h *Handler) getBeans(w http.ResponseWriter, r *http.Request) {
	var resp = &BeansResp{}

	// Call Firestore API
	iter := h.database.Collection("beans").Documents(context.TODO())
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			h.logger.Fatalf("Failed to iterate: %v", err)
		}

		var b Bean
		doc.DataTo(&b)
		resp.Beans = append(resp.Beans, b)
	}

	json.NewEncoder(w).Encode(resp)
}
func (h *Handler) addBean(w http.ResponseWriter, r *http.Request) {
	var (
		req  AddBeanReq
		resp = &AddBeanResp{}
		ctx  = context.TODO()
		err  error
	)

	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Make sure roaster exists
	iter := h.database.Collection("roasters").Where("name", "==", req.Roaster).Documents(ctx)
	for {
		doc, err := iter.Next()
		if doc == nil {
			http.Error(w, "invalid roaster", http.StatusBadRequest)
			return
		}

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		break
	}

	// Add the bean
	doc, _, err := h.database.Collection("beans").Add(ctx, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp.ID = doc.ID

	json.NewEncoder(w).Encode(resp)
}
