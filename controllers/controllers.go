package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"josh/site_status_checker/models"
	"josh/site_status_checker/service"
	"log"
	"net/http"
	"sync"
	"time"
)

var mutex = &sync.Mutex{}
var ctx context.Context

func AddSites(w http.ResponseWriter, r *http.Request) {

	var sd models.PostRequest
	if body, err := ioutil.ReadAll(r.Body); err == nil {
		if err := json.Unmarshal([]byte(body), &sd); err != nil {
			fmt.Println("Dang galti here")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			log.Fatal(err)
		}
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		log.Fatal(err)
	}

	for _, v := range sd.AddedWebsites {
		models.StatusInfo[v] = "unchecked"
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("AddedWebsites added"))
	go startChecking()
}

func GetSiteStatus(w http.ResponseWriter, r *http.Request) {
	site := r.URL.Query().Get("link")
	if site != "" {
		GetSingleSiteStatus(w, r)
		return
	}
	res, err := json.Marshal(models.StatusInfo)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		log.Fatal(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)

}

func GetSingleSiteStatus(w http.ResponseWriter, r *http.Request) {
	site := r.URL.Query().Get("link")
	if _, ok := models.StatusInfo[site]; ok {
		res, err := json.Marshal(models.StatusInfo[site])
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			log.Fatal(err)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(res)
	} else {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(site + " is not stored in the database"))
	}
	return

}

func startChecking() {
	var h service.HttpChecker
	site_channel := make(chan string)

	for k := range models.StatusInfo {
		go func(k string) {
			status, _ := h.CheckSite(ctx, k)
			mutex.Lock()
			if status {
				models.StatusInfo[k] = "working"
			} else {
				models.StatusInfo[k] = "not working"
			}
			mutex.Unlock()
			site_channel <- k
		}(k)

	}

	for k := range site_channel {
		go func(k string) {
			time.Sleep(time.Minute)
			status, _ := h.CheckSite(ctx, k)
			mutex.Lock()
			if status {
				models.StatusInfo[k] = "working"
			} else {
				models.StatusInfo[k] = "not working"
			}
			mutex.Unlock()
			site_channel <- k
		}(k)
	}
}
