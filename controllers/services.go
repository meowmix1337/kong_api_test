package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"kong_api/endpoint"
	"kong_api/responses"

	"github.com/gorilla/mux"
)

func (server *Server) initializeServicesRoutes() {
	// Services routes
	server.Router.HandleFunc(APIV1+"/services", server.CreateService).Methods("POST")
	server.Router.HandleFunc(APIV1+"/services", server.GetServices).Methods("GET")
	server.Router.HandleFunc(APIV1+"/services/{id}", server.GetService).Methods("GET")
	server.Router.HandleFunc(APIV1+"/services/{id}", server.UpdateService).Methods("PUT")
	server.Router.HandleFunc(APIV1+"/services/{id}", server.DeleteService).Methods("DELETE")
}

// CreateService .
func (server *Server) CreateService(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.Error(w, http.StatusUnprocessableEntity, err)
	}

	service := endpoint.Service{}
	err = json.Unmarshal(body, &service)
	if err != nil {
		responses.Error(w, http.StatusUnprocessableEntity, err)
		return
	}

	err = service.Validate()
	if err != nil {
		responses.Error(w, http.StatusUnprocessableEntity, err)
		return
	}

	serviceCreated, err := server.ServicesService.Store(service)
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusCreated, serviceCreated)
}

// GetServices .
func (server *Server) GetServices(w http.ResponseWriter, r *http.Request) {

	queryParams := r.URL.Query()
	keyword := queryParams.Get("keyword")
	sortBy := queryParams.Get("sort_by")
	sortOrder := queryParams.Get("sort_order")

	// TODO: add validation
	page, _ := strconv.Atoi(queryParams.Get("page"))
	pageSize, _ := strconv.Atoi(queryParams.Get("page_size"))

	services, total, err := server.ServicesService.All(keyword, sortBy, sortOrder, page, pageSize)
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	paginatedServices := &endpoint.PaginatedServices{
		Total:    total,
		Services: services,
		// TODO: Add next page if available
		// TODO: Add last page available
	}

	responses.JSON(w, http.StatusOK, paginatedServices)
}

// GetService .
func (server *Server) GetService(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	uid, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		responses.Error(w, http.StatusBadRequest, err)
		return
	}

	serviceRetrieved, err := server.ServicesService.ByID(uint32(uid))
	if err != nil {
		if err == endpoint.ErrServiceNotFound {
			responses.Error(w, http.StatusNotFound, err)
			return
		}
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusOK, serviceRetrieved)
}

// UpdateService .
func (server *Server) UpdateService(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		responses.Error(w, http.StatusBadRequest, err)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.Error(w, http.StatusUnprocessableEntity, err)
		return
	}

	service := endpoint.Service{}
	err = json.Unmarshal(body, &service)
	if err != nil {
		responses.Error(w, http.StatusUnprocessableEntity, err)
		return
	}

	err = service.Validate()
	if err != nil {
		responses.Error(w, http.StatusUnprocessableEntity, err)
		return
	}

	updatedService, err := server.ServicesService.Update(uint32(id), service)
	if err != nil {
		if err == endpoint.ErrServiceNotFound {
			responses.Error(w, http.StatusNotFound, err)
			return
		}
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusOK, updatedService)
}

// DeleteService .
func (server *Server) DeleteService(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		responses.Error(w, http.StatusBadRequest, err)
		return
	}

	err = server.ServicesService.DeleteByID(uint32(id))
	if err != nil {
		if err == endpoint.ErrServiceNotFound {
			responses.Error(w, http.StatusNotFound, err)
			return
		}
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	w.Header().Set("Entity", fmt.Sprintf("%d", id))
	responses.JSON(w, http.StatusOK, "")
}
