package main

import (
	"encoding/json"
	"net/http"

	log "gopkg.in/inconshreveable/log15.v2"
)

func readJSON(req *http.Request, payload interface{}) error {
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(payload)
	if err != nil {
		log.Error("failed to parse body", "err", err)
	}
	return err
}

func writeJSON(rw http.ResponseWriter, payload interface{}) error {
	b, err := json.Marshal(payload)
	if err != nil {
		log.Error("failed to marshal response", "err", err)
		return err
	}
	_, err = rw.Write(b)
	return err
}
