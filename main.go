package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os/exec"
	"strings"
)

type Jail struct {
	Name         string `json:"name"`
	Path         string `json:"path"`
	IP           string `json:"ip"`
	HostHostname string `json:"host_hostname"`
}

func run(cmd string, args ...string) (string, error) {
	c := exec.Command(cmd, args...)
	out, err := c.CombinedOutput()
	return string(out), err
}

func listJails(w http.ResponseWriter, _ *http.Request) {
	out, err := run("jls", "-n")
	if err != nil {
		http.Error(w, out, 500)
		return
	}
	w.Write([]byte(out))
}

func createJail(w http.ResponseWriter, r *http.Request) {
	var j Jail
	if err := json.NewDecoder(r.Body).Decode(&j); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	args := []string{
		"-c",
		"name=" + j.Name,
		"path=" + j.Path,
		"host.hostname=" + j.HostHostname,
		"ip4.addr=" + j.IP,
		"persist",
	}

	out, err := run("jail", args...)
	if err != nil {
		http.Error(w, out, 500)
		return
	}
	w.Write([]byte("jail created\n"))
}

func startJail(w http.ResponseWriter, r *http.Request) {
	name := strings.TrimPrefix(r.URL.Path, "/jails/")
	name = strings.TrimSuffix(name, "/start")
	out, err := run("jail", "-c", "name="+name)
	if err != nil {
		http.Error(w, out, 500)
		return
	}
	w.Write([]byte("jail started\n"))
}

func stopJail(w http.ResponseWriter, r *http.Request) {
	name := strings.TrimPrefix(r.URL.Path, "/jails/")
	name = strings.TrimSuffix(name, "/stop")
	out, err := run("jail", "-r", name)
	if err != nil {
		http.Error(w, out, 500)
		return
	}
	w.Write([]byte("jail stopped\n"))
}

func deleteJail(w http.ResponseWriter, r *http.Request) {
	name := strings.TrimPrefix(r.URL.Path, "/jails/")
	out, err := run("jail", "-r", name)
	if err != nil {
		http.Error(w, out, 500)
		return
	}
	w.Write([]byte("jail removed\n"))
}

func execInJail(w http.ResponseWriter, r *http.Request) {
	name := strings.TrimPrefix(r.URL.Path, "/jails/")
	name = strings.TrimSuffix(name, "/exec")

	var body struct {
		Command []string `json:"command"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	args := append([]string{name}, body.Command...)
	out, err := run("jexec", args...)
	if err != nil {
		http.Error(w, out, 500)
		return
	}
	w.Write([]byte(out))
}

func main() {
	http.HandleFunc("/jails", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			listJails(w, r)
		case http.MethodPost:
			createJail(w, r)
		default:
			http.Error(w, "method not allowed", 405)
		}
	})

	http.HandleFunc("/jails/", func(w http.ResponseWriter, r *http.Request) {
		switch {
		case strings.HasSuffix(r.URL.Path, "/start"):
			startJail(w, r)
		case strings.HasSuffix(r.URL.Path, "/stop"):
			stopJail(w, r)
		case strings.HasSuffix(r.URL.Path, "/exec"):
			execInJail(w, r)
		case r.Method == http.MethodDelete:
			deleteJail(w, r)
		default:
			http.Error(w, "not found", 404)
		}
	})

	log.Println("FreeBSD Jail API listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
