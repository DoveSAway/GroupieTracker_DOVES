package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"text/template"
)

const portLocalHost = ":8080"

func setArtist() []Artist {
	reponse, _ := http.Get("https://groupietrackers.herokuapp.com/api/artists")
	var artist []Artist
	body, _ := ioutil.ReadAll(reponse.Body)
	json.Unmarshal(body, &artist)
	return artist
}

func SetInfo() Relation {
	response, _ := http.Get("https://groupietrackers.herokuapp.com/api/relation")
	data, _ := ioutil.ReadAll(response.Body)
	var relation Relation
	json.Unmarshal([]byte(data), &relation)
	return relation
}

type Artist struct {
	Id           int      `json: "id"`
	Image        string   `json: "image"`
	Name         string   `json: "name"`
	Members      []string `json: "members"`
	CreationDate int      `json: "creationDate"`
	FirstAlbum   string   `json: "firstAlbum"`
}

type Relation struct {
	M_Index []struct {
		M_Id  int                 `json:"id"`
		M_Rel map[string][]string `json:"datesLocations"`
		Name  string
	} `json:"index"`
}

type RelaGroup struct {
	M_Id  int
	M_Rel map[string][]string
	Name  string
}

type Groupie struct {
	Datartist []Artist
	Datarela  []RelaGroup
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("usingAPI/Home.html")
	if err != nil {
		fmt.Fprint(w, err)
	}
	var exGroupie Groupie
	exGroupie.Datartist = setArtist()
	tmp := SetInfo()
	for i := range tmp.M_Index {
		exGroupie.Datarela = append(exGroupie.Datarela, RelaGroup(tmp.M_Index[i]))
		exGroupie.Datarela[i].Name = exGroupie.Datartist[i].Name
	}
	errExe := t.Execute(w, exGroupie)
	fmt.Println(errExe)
}

func main() {
	http.FileServer(http.Dir("usingAPI"))
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/", HomeHandler)
	http.ListenAndServe(portLocalHost, nil)

}
