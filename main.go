package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"slices"
	"strconv"

	"github.com/fiatjaf/eventstore/sqlite3"
	"github.com/fiatjaf/khatru"
	"github.com/nbd-wtf/go-nostr"
)

func main() {
	acceptableKinds := []int{0, 3}
	fmt.Println(acceptableKinds)

	db := sqlite3.SQLite3Backend{DatabaseURL: "/data/historelay.sqlite?cache=shared&mode=rwc&_journal_mode=WAL"}
	if err := db.Init(); err != nil {
		panic(err)
	}

	relay := khatru.NewRelay()

	relay.Info.Name = "Historelay"
	relay.Info.Description = "Keep replaceable events history"

	relay.RejectEvent = append(relay.RejectEvent, func(ctx context.Context, event *nostr.Event) (reject bool, msg string) {
		if slices.Contains(acceptableKinds, event.Kind) {
			return false, ""
		} else {
			return true, "kind " + strconv.Itoa(event.Kind) + " is not supported"
		}
	})

	relay.RejectFilter = append(relay.RejectFilter, func(ctx context.Context, filter nostr.Filter) (reject bool, msg string) {
		if len(filter.Authors) > 0 {
			return false, ""
		} else {
			return true, "authors is required"
		}
	})

	relay.StoreEvent = append(relay.StoreEvent, db.SaveEvent)
	relay.QueryEvents = append(relay.QueryEvents, db.QueryEvents)

	mux := relay.Router()
	mux.HandleFunc("/", indexHandler)

	fmt.Println("Running")
	http.ListenAndServe("", relay)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	t, err := template.ParseFiles("./index.html")
	if err != nil {
		log.Fatal("Parse: ", err)
		return
	}
	t.Execute(w, nil)
}
