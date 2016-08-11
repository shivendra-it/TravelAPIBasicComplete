package main

import (
	"io/ioutil"
	"net/http"
)


func handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	jsonbody, err := ioutil.ReadAll(r.Body)
//	fmt.Println(string(jsonbody))
	r.Body.Close()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(jsonbody))
		return
	}
	startdb()

		top_five_by_price_rating()
//	top_five_by_offer()
// 	NearestCity("17.422","78.33")
//	NearestHotel("17.422","78.33","2162254155836171767")
//	top_five_by_price()
//	top_five_by_rating("1")
//	hotel_has_amenities("Internet")
//	city_data_insert([]byte(jsonbody))
//	hotel_data_insert([]byte(jsonbody))
//	SearchCountryNode("India")
//	queryAllNodes()
//	LinkRatingWithHotel()
//	InsertHotelData([]byte(jsonbody))

}
