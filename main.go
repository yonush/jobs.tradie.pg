package main

//TODO
// - use CLI or config to permit alternate IP & port binding
import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/jackc/pgx/v5/stdlib" //use pgx in database/sql mode
)

// PostgreSQl configuration if not passed as env variabbles
const (
	host     = "localhost" //127.0.0.1
	port     = 5432
	user     = "postgres"
	password = "postgres"
	dbname   = "ESD"
)

var db *sql.DB

var bindport string = "80" //change to suit requirement
var wait time.Duration

var index *View
var jobs *View

// setup some route endpoints
func newRouter() *mux.Router {
	r := mux.NewRouter()

	//default handler
	r.HandleFunc("/", indexHandler).Methods("GET")

	// setup static content route - strip ./assets/assets/[resource]
	// to keep /assets/[resource] as a route
	staticFileDirectory := http.Dir("./assets/")
	staticFileHandler := http.StripPrefix("/assets/", http.FileServer(staticFileDirectory))
	r.PathPrefix("/assets/").Handler(staticFileHandler).Methods("GET")

	// setup routes to handle job updates and job notes
	r.HandleFunc("/jobs", getJobsHandler).Methods("GET")
	r.HandleFunc("/jobs/{s:[0-9]+}", getJobsHandler).Methods("GET")
	r.HandleFunc("/jobs/{s:[0-9]+}/{f:[s,a,i,t,c]?}", getJobsHandler).Methods("GET")

	r.HandleFunc("/job/{id:[0-9]+}", getJobHandler).Methods("GET")
	r.HandleFunc("/job/{id:[0-9]+}", editJobHandler).Methods("POST")

	r.HandleFunc("/notes/{id:[0-9]+}", getJobNoteHandler).Methods("GET")
	r.HandleFunc("/notes/{id:[0-9]+}", editJobNoteHandler).Methods("POST")
	return r
}

// used to auto detect the active local IP address - not used yet
func GetOutboundIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String()
}

func main() {

	//check if a different bind port was passed from the CLI
	if len(os.Args) > 1 {
		s := os.Args[1]

		if _, err := strconv.ParseInt(s, 10, 64); err == nil {
			bindport = s
		}
	}

	var err error

	// Create a string that will be used to make a connection later
	// Note Password has been left out, which is best to avoid issues when using null password
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	log.Println("Connecting to PostgreSQL on")
	log.Println(psqlInfo)
	db, err := sql.Open("pgx", psqlInfo)
	if err != nil {
		log.Println("Invalid DB arguments, or github.com/lib/pq not installed")
		log.Fatal(err)
	}

	defer func() {
		db.Close()
		log.Printf("Database connection closed")
	}()

	// Ping database (connection is only established at this point, open only validates arguments passed to it)
	//optional code
	err = db.Ping()
	if err != nil {
		log.Fatal("Connection to specified database failed: ", err)
	}

	log.Println("Connected successfully")

	//check data import status
	_, err = os.Stat("./imported.txt")
	if os.IsNotExist(err) {
		log.Println("Importing demo data")
		loadFromJson(db, "./jobs.json")
	}

	// load templates
	log.Println("Loading templates")
	index = NewView("bootstrap", "views/index.gohtml")
	jobs = NewView("bootstrap", "views/jobs.gohtml")

	log.Println("Starting HTTP service on " + bindport)
	r := newRouter()

	// setup HTTP on gorilla mux for a gracefull shutdown
	srv := &http.Server{
		Addr: "0.0.0.0:" + bindport,

		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r,
	}

	// HTTP listener is in a goroutine as its blocking
	go func() {
		if err = srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	// setup a ctrl-c trap to ensure a graceful shutdown
	// this would also allow shutting down other pipes/connections. eg DB
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	srv.Shutdown(ctx)
	log.Println("shutting down")
	os.Exit(0)
}
