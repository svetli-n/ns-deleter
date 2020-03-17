package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"
)

func main() {
	client := NewClient()
	conf := NewConf()
	duration := time.Duration(int(conf.CheckFreq) * int(time.Second))
	stop := make(chan struct{})
	nsToKeep := make(map[string]bool)
	keepChIn := make(chan map[string]bool, 1)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go func() {
		keepChIn <- nsToKeep
	}()

	go loopDeleteNs(conf, client, duration, stop, keepChIn)

	go func(stop chan struct{}) {
		for sig := range c {
			log.Println("Got CTRL+C:", sig)
			stop <- struct{}{}
			os.Exit(0)
		}
	}(stop)

	helloHandler := func(w http.ResponseWriter, req *http.Request) {
		ns, ok := req.URL.Query()["ns"]
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Bad request"))
		}
		go func() {
			nsToKeep := <-keepChIn
			nsToKeep[ns[0]] = true
			keepChIn <- nsToKeep
		}()
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("200 OK: Standard response for successful HTTP requests."))
	}
	http.HandleFunc("/keep", helloHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))

}

func loopDeleteNs(conf appConf, client Client, duration time.Duration, stop chan struct{}, keepChIn chan map[string]bool) {
	for {
		select {
		case <-stop:
			return
		case toKeep := <-keepChIn:
			deleteNs(conf, client, toKeep)
			keepChIn <- toKeep
			time.Sleep(duration)
		}
	}
}

func deleteNs(conf appConf, client Client, nsToKeep map[string]bool) {
	nsMaxAge := conf.NsMaxAge
	namespaces := client.namespaces()
	for _, ns := range namespaces.Items {
		curUser := strings.Split(ns.Name, "-")[0]
		nsAgeHours := time.Now().UTC().Sub(ns.CreationTimestamp.Time).Seconds()
		user, isIn := conf.Users[curUser]
		shoulKeep, notified := nsToKeep[ns.Name]
		if isIn && nsAgeHours > nsMaxAge {
			if !notified || shoulKeep {
				log.Printf("Namespace: %s name: %s created: %f email: %s\n", ns.Name, curUser, nsAgeHours, user.Email)
				emailText1 := fmt.Sprintf("%s will be deleted\n", ns.Name)
				link := fmt.Sprintf("http://localhost:8080/keep?ns=%s", ns.Name)
				portFwd := "kubectl -n monitoring port-forward deployment/ns-deleter 8080:8080"
				emailText2 := fmt.Sprintf("but you can run `%s`\nand follow %s to keep it", portFwd, link)
				send(user.Email, emailText1+emailText2)
				nsToKeep[ns.Name] = false
			} else {
				if !shoulKeep {
					log.Println("Deleting", ns.Name)
					client.delete(ns.Name)
					delete(nsToKeep, ns.Name)
				}
			}
		}
	}
}
