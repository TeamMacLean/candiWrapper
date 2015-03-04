package main

import (
"fmt"
"log"
"math/rand"
"net/http"
"github.com/ttacon/chalk"
"github.com/garyburd/redigo/redis"
)

type Session struct {
	Id		string		`param:"id"`
}

var Sessions = []Session{}

type test_struct struct {
	Test string
}

var (
	port = ":8080"
	redisConnection = initRedis()
	letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	randomStringLength = 10
	rootURL = "http://candisnp.tsl.ac.uk/martin/"
	keyName = "galaxyData"
)

func randSeq() string {
	b := make([]rune, randomStringLength)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func initRedis() redis.Conn {
	redisConnection, err := redis.Dial("tcp", ":6379")
	if err != nil {
		panic(err)
	}
	// defer redisConnection.Close()
	return redisConnection
}

func addToRedis(key string, value string){
	fmt.Println(key+" "+value)
	redisConnection.Do("SET", key, value)
}

func getFromRedis(key string) string{
	value, err := redis.String(redisConnection.Do("GET", key))
	if err != nil {
		fmt.Println(chalk.Red,"key not found", chalk.Reset)
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
				fmt.Fprintf(w, rootURL+"?session="+id)

			}
		}

		} else if(reqType == "GET") {

			shortcode := req.URL.Path[1:]

			if len(shortcode) > 0 {
				log.Println(chalk.Green,"shortcode: "+shortcode,chalk.Reset)
				fromRedis := getFromRedis(shortcode)

				if len(fromRedis) > 0 {

					js, err := json.Marshal(fromRedis)
					w.Header().Set("Content-Type", "application/json")
					w.Write(js)				
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