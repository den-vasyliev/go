package main

import (
	"fmt"
	"math"
	"net/url"
	"strings"
	"strconv"
	"reflect"
	"github.com/go-redis/redis"
	"net/http"
)

func main() {

	var c float64

	R := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})


	R.HSet("CONF", "rev", "test").Err()


	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var mm map[string]float64
		mm = make(map[string]float64)

		m, _ := url.ParseQuery(strings.ToUpper(r.URL.EscapedPath()))

		for key, value := range m {
			switch key {
			case "DURATION","RATEA","RATEB","BILLSEC","MTC":
				valueI,_ :=strconv.Atoi(value[0])
				valueF :=float64(valueI)
				mm[key]=valueF
				fmt.Println("Key:", key, "Value:", valueF,"Type:", reflect.TypeOf(valueF))
			default : fmt.Println("Key:", key, "Value:", value[0])

			}
		}

		if (m["LASTAPP"][0] == "DIAL") && (m["DISPOSITION"][0] == "ANSWERED") && (mm["BILLSEC"]>0){
			switch m["DCONTEXT"][0] {
			case "CALLME_REG":
				c = math.Ceil( ( (mm["DURATION"]/60) * mm["MTC"] ) + mm["DURATION"] * (mm["RATEA"]/60))  + mm["BILLSEC"] * (mm["RATEB"]/60)
			case "DID_REG":
				c = math.Ceil( ( (mm["BILLSEC"]/60) * mm["MTC"] ) + mm["BILLSEC"] * (mm["RATEA"]/60))
			default:
				c = 0.0
			}
		} else if (m["DCONTEXT"][0] == "CALLME_REG" && m["LASTAPP"][0] == "FORKCDR"){
			c = math.Ceil( (mm["BILLSEC"]/60) * mm["MTC"] )

		}
		fmt.Printf("Call rate:%6.2f\n", c/100)
		fmt.Fprintf(w,"Call rate:%6.2f\n", c/100)

		_, err := R.HSet("CONF", "rev", c).Result()
		fmt.Printf("Result",err)
	})

	http.ListenAndServe(":8080", nil)

}
