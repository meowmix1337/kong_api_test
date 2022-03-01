package controllers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"kong_api/endpoint"
	"kong_api/responses"
)

func (server *Server) initializeVersionsRoutes() {
	// Services routes
	server.Router.HandleFunc(APIV1+"/versions", server.AssociateVersion).Methods("POST")
}

// AssociateVersion .
func (server *Server) AssociateVersion(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.Error(w, http.StatusUnprocessableEntity, err)
	}

	version := endpoint.Version{}
	err = json.Unmarshal(body, &version)
	if err != nil {
		responses.Error(w, http.StatusUnprocessableEntity, err)
		return
	}

	err = version.Validate()
	if err != nil {
		responses.Error(w, http.StatusUnprocessableEntity, err)
		return
	}

	versionCreated, err := server.VersionService.Store(version)
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusCreated, versionCreated)
}
