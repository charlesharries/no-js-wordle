package main

import (
	"bufio"
	"encoding/gob"
	"flag"
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
	addr := flag.String("addr", "localhost:4001", "HTTP network address")
	devMode := flag.Bool("dev", false, "Enable debug mode")
	secret := flag.String("secret", "change_this", "Secret key")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	templateCache, err := newTemplateCache("./ui/templates")
	if err != nil {
		log.Fatal(err)
	}

	session := sessions.New([]byte(*secret))
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

	srv := &http.Server{
		Addr:    *addr,
		Handler: app.routes(),
	}

	log.Printf("Starting server at http://%s\n", *addr)
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
