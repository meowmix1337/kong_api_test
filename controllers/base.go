package controllers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"

	_ "github.com/jinzhu/gorm/dialects/mysql" // mysql database driver

	"kong_api/models"
	"kong_api/service"
)

// Server .
type Server struct {
	Router          *mux.Router
	ServicesService service.IServicesService
	VersionService  service.IVersionsService
}

const (
	APIV1 = "/v1"
)

// Initialize .
func (server *Server) Initialize(Dbdriver, DbUser, DbPassword, DbPort, DbHost, DbName string) {

	var err error
	var db *gorm.DB

	if Dbdriver == "mysql" {
		DBURL := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", DbUser, DbPassword, DbHost, DbPort, DbName)
		db, err = gorm.Open(Dbdriver, DBURL)
		if err != nil {
			fmt.Printf("Cannot connect to %s database", Dbdriver)
			log.Fatal("This is the error:", err)
		} else {
			fmt.Printf("We are connected to the %s database", Dbdriver)
		}
	}

	serviceCore := service.NewServiceCore(db)
	servicesService := service.NewServicesService(serviceCore)
	versionsService := service.NewVerionsService(serviceCore)

	server.ServicesService = servicesService
	server.VersionService = versionsService

	serviceCore.DB.AutoMigrate(&models.Service{}, &models.Version{}) // database migration

	server.Router = mux.NewRouter()

	server.initializeRoutes()
}

// Run .
func (server *Server) Run(addr string) {
	fmt.Println("Listening to port 8080")
	log.Fatal(http.ListenAndServe(addr, server.Router))
}
