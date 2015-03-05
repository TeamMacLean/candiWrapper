package main

import (
"fmt"
"log"
"math/rand"
"encoding/json"
"net/http"
"time"
"github.com/garyburd/redigo/redis"
)

var (
	port = ":8080"
	redisConnection = initRedis()
	letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	randomStringLength = 10
	rootURL = "http://candisnp.tsl.ac.uk/martin/"
	keyName = "data"
	refName = "ref"
	)

func randSeq() string {

	rand.Seed(time.Now().UTC().UnixNano())
	b := make([]rune, randomStringLength)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}

	sOut := string(b)

	if(len(getFromRedis(sOut)) > 0){
		return randSeq()
	} else {
		return sOut
		
	}
	
}

func initRedis() redis.Conn {
	redisConnection, err := redis.Dial("tcp", ":6379")
	if err != nil {
		panic(err)
	}
	return redisConnection
}

func addToRedis(key string, value string){
	fmt.Println(key+" "+value)
	redisConnection.Do("SET", key, value)
}

func getFromRedis(key string) string{
	value, err := redis.String(redisConnection.Do("GET", key))
	if err != nil {
		fmt.Println("key does not exist in DB")
	}
	return value
}

func logIt(w http.ResponseWriter, message string){
	fmt.Fprintf(w, message)
	log.Println(message)
}

func handlePost(w http.ResponseWriter, req *http.Request){
	err := req.ParseForm()

	if err != nil {
		logIt(w, "could not parse form")
	} else {


		log.Println(req.Form)

		galaxyData := req.FormValue(keyName)
		refData := req.FormValue(refName)

		galaxyDataLength := len(galaxyData)
		refDataLength := len(refData)

		if galaxyDataLength > 0 && refDataLength > 0{

			id := randSeq();
			addToRedis(id, string(galaxyData));
			logIt(w, rootURL+"?session="+id+"&species=athalianaTair10")
		}
	}
}
func handleGet(w http.ResponseWriter, req *http.Request){
	shortcode := req.URL.Path[1:]

	if len(shortcode) > 0 {
		log.Println("shortcode: "+shortcode)
		fromRedis := getFromRedis(shortcode)

		if len(fromRedis) > 0 {

			js, err := json.Marshal(fromRedis)
			if err != nil {
				logIt(w, "could not convert to json")
			} else {
				w.Header().Set("Content-Type", "application/json")
				w.Write(js)
			}

		}
	}
}

func handler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	reqType := req.Method
	log.Println("type: "+reqType)

	if reqType == "POST" {
		handlePost(w, req)
	} else if reqType == "GET" {
		handleGet(w,req)			
	} else {
		logIt(w, "did not receive GET or POST")
	}
}

func main() {
	defer redisConnection.Close()
	defer log.Println("Server Stopped")
	http.HandleFunc("/", handler)
	log.Println("Starting server on port "+port)
	http.ListenAndServe(port, nil)

}
