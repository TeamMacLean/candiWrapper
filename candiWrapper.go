package main

import (
"fmt"
"log"
"math/rand"
"encoding/json"
"net/http"
"time"
"github.com/ttacon/chalk"
"github.com/garyburd/redigo/redis"
)

var (
	port = ":8080"
	redisConnection = initRedis()
	letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	randomStringLength = 10
	rootURL = "http://candisnp.tsl.ac.uk/martin/"
	keyName = "galaxyData"
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

func handler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	reqType := req.Method
	log.Println(chalk.Blue,"type: "+reqType,chalk.Reset)

	if reqType == "POST" {
		err := req.ParseForm()

		if err != nil {
			log.Println(chalk.Red,"could not parse form", chalk.Reset)
		} else {

			galaxyData := req.FormValue(keyName)
			galaxyDataLength := len(galaxyData)

			if galaxyDataLength > 0 {

				id := randSeq();
				addToRedis(id, galaxyData);
				
				fmt.Fprintf(w, rootURL+"?session="+id+"&species=athalianaTair10")
				log.Println(rootURL+"?session="+id+"&species=athalianaTair10")

			}
		}

		} else if(reqType == "GET") {

			shortcode := req.URL.Path[1:]

			if len(shortcode) > 0 {
				log.Println(chalk.Green,"shortcode: "+shortcode,chalk.Reset)
				fromRedis := getFromRedis(shortcode)

				if len(fromRedis) > 0 {

					js, err := json.Marshal(fromRedis)
					if err != nil {
						http.Error(w, err.Error(), http.StatusInternalServerError)
						return
					}
					w.Header().Set("Content-Type", "application/json")
					w.Write(js)
					log.Println("sent data")

				}
			}
		}
	}

	func main() {
		defer redisConnection.Close()
		http.HandleFunc("/", handler)
		log.Println(chalk.Green,"Starting server on port "+port, chalk.Reset)
		http.ListenAndServe(port, nil)

	}
