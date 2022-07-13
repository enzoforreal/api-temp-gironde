package controllers

import (
	"encoding/json"
	"github.com/revel/revel"
	"github.com/revel/revel/cache"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type Temperature_t struct {
	Count      int      `json:"count"`
	First      string   `json:"first"`
	Last       string   `json:"last"`
	Prev       string   `json:"prev"`
	Next       string   `json:"next"`
	ApiVersion string   `json:"api_version"`
	Data       []Data_t `json:"data"`
}

type Data_t struct {
	CodeStation          string  `json:"code_station"`
	LibelleStation       string  `json:"libelle_station"`
	UriStation           string  `json:"uri_station"`
	Longitude            float64 `json:"longitude"`
	Latitude             float64 `json:"latitude"`
	CodeCommune          string  `json:"code_commune"`
	LibelleCommune       string  `json:"libelle_commune"`
	CodeCoursEau         string  `json:"code_cours_eau"`
	LibelleCoursEau      string  `json:"libelle_cours_eau"`
	UriCoursEau          string  `json:"uri_cours_eau"`
	CodeParametre        int32   `json:"code_parametre"`
	LibelleParametre     string  `json:"libelle_parametre"`
	DateMesureTemp       string  `json:"date_mesure_temp"`
	HeureMesureTemp      string  `json:"heure_mesure_temp"`
	Resultat             float64 `json:"resultat"`
	CodeUnite            int32   `json:"code_unite"`
	SymboleUnite         string  `json:"symbole_unite"`
	CodeQualification    int32   `json:"code_qualification"`
	LibelleQualification string  `json:"libelle_qualification"`
}

type App struct {
	*revel.Controller
}

func checkerr(err error) {
	if err != nil {
		log.Println(err)
	}
}

func GetTemp() []byte {
	var temp []byte
	if err := cache.Get("temp-cache", &temp); err != nil {
		url := "https://hubeau.eaufrance.fr/api/v1/temperature/chronique?code_station=05025600&size=5000&sort=desc&pretty"

		req, _ := http.NewRequest("GET", url, nil)

		res, _ := http.DefaultClient.Do(req)

		defer res.Body.Close()
		body, _ := ioutil.ReadAll(res.Body)
		cache.Set("temp-cache", body, 60*time.Second)
		return body
	}
	return temp
}

func (c App) Index() revel.Result {
	var data_tmp Temperature_t
	temperature := GetTemp()
	err := json.Unmarshal(temperature, &data_tmp)
	checkerr(err)

	c.ViewArgs["data"] = data_tmp.Data
	c.ViewArgs["data_encode"] = string(temperature)
	return c.Render()
}
