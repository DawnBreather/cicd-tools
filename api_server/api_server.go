package api_server

import (
	"crypto/tls"
	"github.com/DawnBreather/go-commons/ssl"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

// ApiServer has router
type ApiServer struct {
	Ssl ssl.Ssl
	Router *mux.Router
	Config map[string]interface{}
}

func (a *ApiServer) SetSslMetadata(locality, country, organization, province, postalCode, streetAddress string){
	a.Ssl.
		SetLocality(locality).
		SetCountry(country).
		SetOrganization(organization).
		SetProvince(province).
		SetPostalCode(postalCode).
		SetStreetAddress(streetAddress)
}

// Initialize initializes the app with predefined configuration
func (a *ApiServer) Initialize(config map[string]interface{}) *ApiServer{
//func (a *ApiServer) Initialize() *ApiServer{
	a.Ssl.InitializeCertificateAuthority()
	a.Router = mux.NewRouter()
	return a
}

//func (a *ApiServer) setRouters() {
//	// Routing for handling the projects
//	a.Get("/projects", a.handleRequest(handler.GetAllProjects))
//	a.Post("/projects", a.handleRequest(handler.CreateProject))
//	a.Get("/projects/{title}", a.handleRequest(handler.GetProjects))
//	a.Put("/projects/{title}", a.handleRequest(handler.UpdateProject))
//	a.Delete("/projects/{title}", a.handleRequest(handler.DeleteProject))
//	a.Put("/projects/{title}/archive", a.handleRequest(handler.ArchiveProject))
//	a.Delete("/projects/{title}/archive", a.handleRequest(handler.RestoreProject))
//
//	// Routing for handling the tasks
//	a.Get("/projects/{title}/tasks", a.handleRequest(handler.GetAllTasks))
//	a.Post("/projects/{title}/tasks", a.handleRequest(handler.CreateTask))
//	a.Get("/projects/{title}/tasks/{id:[0-9]+}", a.handleRequest(handler.GetTask))
//	a.Put("/projects/{title}/tasks/{id:[0-9]+}", a.handleRequest(handler.UpdateTask))
//	a.Delete("/projects/{title}/tasks/{id:[0-9]+}", a.handleRequest(handler.DeleteTask))
//	a.Put("/projects/{title}/tasks/{id:[0-9]+}/complete", a.handleRequest(handler.CompleteTask))
//	a.Delete("/projects/{title}/tasks/{id:[0-9]+}/complete", a.handleRequest(handler.UndoTask))
//}

// Get wraps the router for GET method
func (a *ApiServer) Get(path string, f func(w http.ResponseWriter, r *http.Request)) *ApiServer{
	a.Router.HandleFunc(path, f).Methods("GET")
	return a
}

// Post wraps the router for POST method
func (a *ApiServer) Post(path string, f func(w http.ResponseWriter, r *http.Request)) *ApiServer{
	a.Router.HandleFunc(path, f).Methods("POST")
	return a
}

// Put wraps the router for PUT method
func (a *ApiServer) Put(path string, f func(w http.ResponseWriter, r *http.Request)) *ApiServer{
	a.Router.HandleFunc(path, f).Methods("PUT")
	return a
}

// Delete wraps the router for DELETE method
func (a *ApiServer) Delete(path string, f func(w http.ResponseWriter, r *http.Request)) *ApiServer{
	a.Router.HandleFunc(path, f).Methods("DELETE")
	return a
}

// Run the app on it's router
func (a *ApiServer) Run(host string) {
	log.Fatal(http.ListenAndServe(host, a.Router))
}

func (a *ApiServer) RunSslSelfSigned(host string, DNSNames []string) {
	_, _, keyPair := a.Ssl.GenerateSignedCertificate(DNSNames)
	server := &http.Server{
		Addr:              host,
		Handler:           a.Router,
		TLSConfig:         &tls.Config{
			Certificates: []tls.Certificate{keyPair},
		},
	}

	//log.Fatal(http.ListenAndServeTLS(host, "", "", a.Router))
	log.Fatal(server.ListenAndServeTLS("", ""))
}

type RequestHandlerFunction func(w http.ResponseWriter, r *http.Request)

func (a *ApiServer) handleRequest(handler RequestHandlerFunction) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handler(w, r)
	}
}