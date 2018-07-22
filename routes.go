package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
)

var router *httprouter.Router
var jsonRes []byte

type apiResponse struct {
	Success bool     `json:"success"`
	Errors  string   `json:"errors,omitempty"`
	Friends []string `json:"friends,omitempty"`
	Count   int      `json:"count,omitempty"`
}

func init() {
	router = httprouter.New()
	router.POST("/api/friends", createFriendsHandler)
	router.GET("/api/friends", getFriendsListHandler)
	router.GET("/api/friends/common", getCommonFriendsListHandler)
	router.POST("/api/friends/subscribe", subscribeUpdatesHandler)
}

func createFriendsHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if err := r.ParseForm(); err != nil {
		jsonRes, _ = json.Marshal(&apiResponse{Success: false, Errors: err.Error()})
		w.Write(jsonRes)
		return
	}

	friendsData := r.FormValue("friends")
	if friendsData == "" {
		jsonRes, _ = json.Marshal(&apiResponse{Success: false, Errors: "no users were provided"})
		w.Write(jsonRes)
		return
	}

	users := []string{}
	if err := json.Unmarshal([]byte(friendsData), &users); err != nil {
		jsonRes, _ = json.Marshal(&apiResponse{Success: false, Errors: err.Error()})
		w.Write(jsonRes)
	}

	if err := createFriends(users); err != nil {
		jsonRes, _ = json.Marshal(&apiResponse{Success: false, Errors: err.Error()})
		w.Write(jsonRes)
		return
	}

	jsonRes, _ = json.Marshal(&apiResponse{Success: true})
	w.Write(jsonRes)
}

func getFriendsListHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	bodyBytes, _ := ioutil.ReadAll(r.Body)
	a := struct {
		Email string
	}{}
	if err := json.Unmarshal(bodyBytes, &a); err != nil {
		jsonRes, _ = json.Marshal(&apiResponse{Success: false, Errors: err.Error()})
		w.Write(jsonRes)
		return
	}
	if a.Email == "" {
		jsonRes, _ = json.Marshal(&apiResponse{Success: false, Errors: "no email was provided"})
		w.Write(jsonRes)
		return
	}

	friendsList, count, err := getFriendsList(a.Email)
	if err != nil {
		jsonRes, _ = json.Marshal(&apiResponse{Success: false, Errors: err.Error(), Friends: friendsList, Count: count})
		w.Write(jsonRes)
		return
	}

	jsonRes, _ = json.Marshal(&apiResponse{Success: true, Errors: "", Friends: friendsList, Count: count})
	w.Write(jsonRes)
}

func getCommonFriendsListHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	bodyBytes, _ := ioutil.ReadAll(r.Body)
	users := struct {
		Friends []string
	}{}
	if err := json.Unmarshal(bodyBytes, &users); err != nil {
		jsonRes, _ = json.Marshal(&apiResponse{Success: false, Errors: err.Error()})
		w.Write(jsonRes)
		return
	}
	if len(users.Friends) == 0 {
		jsonRes, _ = json.Marshal(&apiResponse{Success: false, Errors: "no users were provided"})
		w.Write(jsonRes)
		return
	}

	friendsList, count, err := getCommonFriendsList(users.Friends)
	if err != nil {
		jsonRes, _ = json.Marshal(&apiResponse{Success: false, Errors: err.Error(), Friends: friendsList, Count: count})
		w.Write(jsonRes)
		return
	}

	jsonRes, _ = json.Marshal(&apiResponse{Success: true, Errors: "", Friends: friendsList, Count: count})
	w.Write(jsonRes)
}

func subscribeUpdatesHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	bodyBytes, _ := ioutil.ReadAll(r.Body)

	users := relationship{}
	if err := json.Unmarshal(bodyBytes, &users); err != nil {
		jsonRes, _ = json.Marshal(&apiResponse{Success: false, Errors: err.Error()})
		w.Write(jsonRes)
		return
	}

	errors := []string{}
	if users.Requestor == "" {
		errors = append(errors, "no requestor was provided")
	}
	if users.Target == "" {
		errors = append(errors, "no target was provided")
	}
	if len(errors) > 0 {
		jsonRes, _ = json.Marshal(&apiResponse{Success: false, Errors: strings.Join(errors, ",")})
		w.Write(jsonRes)
		return
	}

	if err := subscribeUpdates(users.Requestor, users.Target); err != nil {
		jsonRes, _ = json.Marshal(&apiResponse{Success: false, Errors: err.Error()})
		w.Write(jsonRes)
		return
	}

	jsonRes, _ = json.Marshal(&apiResponse{Success: true})
	w.Write(jsonRes)
}
