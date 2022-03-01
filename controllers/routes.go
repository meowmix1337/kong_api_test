package controllers

func (server *Server) initializeRoutes() {
	// Services routes
	server.initializeServicesRoutes()
	server.initializeVersionsRoutes()
}
