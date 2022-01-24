package main

import (
	"bufio"
	"encoding/gob"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"text/template"
	"time"

	"github.com/golangcollege/sessions"
)

type application struct {
	devMode       *bool
	errorLog      *log.Logger
	infoLog       *log.Logger
	session       *sessions.Session
	templateCache map[string]*template.Template
	words         []string
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "4001"
	}
	addr := fmt.Sprintf(":%s", port)

	secret := os.Getenv("APP_SECRET")
	if secret == "" {
		secret = "change_this"
	}

	devMode := flag.Bool("dev", false, "Enable debug mode")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	templateCache, err := newTemplateCache("./ui/templates")
	if err != nil {
		log.Fatal(err)
	}

	session := sessions.New([]byte(secret))
	session.Lifetime = 12 * time.Hour
	gob.Register(game{})

	words, err := getWords()
	if err != nil {
		log.Fatal(err)
	}

	app := &application{
		devMode:       devMode,
		infoLog:       infoLog,
		errorLog:      errorLog,
		session:       session,
		templateCache: templateCache,
		words:         words,
	}

	if *devMode {
		addr = fmt.Sprintf("localhost:%s", port)
	}

	srv := &http.Server{
		Addr:    addr,
		Handler: app.routes(),
	}

	log.Printf("Starting server at http://%s\n", addr)
	err = srv.ListenAndServe()
	log.Fatal(err)
}

func getWords() ([]string, error) {
	words := []string{}
	file, err := os.Open("./data/words.txt")
	if err != nil {
		return words, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		words = append(words, scanner.Text())
	}

	return words, scanner.Err()
}
