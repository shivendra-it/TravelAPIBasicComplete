package main

import (
	"fmt"
	"github.com/jmcvetta/neoism"
	"strconv"
)

var (
	db *neoism.Database
)

func panicErr(err error) {
	if err != nil {
		panic(err)
	}
}


func startdb(){
	var err error
  db, err = neoism.Connect("http://neo4j:12345@localhost:7474")
  if err != nil {
    panic(err)
  }
}


func createCityNode(c1 string,c2 string) *neoism.Node{
	n, err := db.CreateNode(neoism.Props{"city_id": c1,"city_name": c2})
	if err != nil {
		panic(err)
	}
	n.AddLabel("CityInfo")
	return n
}


func createHotelNode(s1 string,s2 string,s3 string,s4 string,s5 string,s6 string,s7 string,s8 string,s9 string,s10 string) *neoism.Node{
	n, err := db.CreateNode(neoism.Props{"hotel_name": s1,"latitude": s2,"longitude":s3,"actual_price":s4,"room_count":s5,"image_url":s6,"offers":s7,"city_id":s8,"discounted_price":s9,"hotel_id":s10})
	if err != nil {
		panic(err)
	}
	n.AddLabel("HotelInfo")
	return n
}

func LinkRatingWithHotel(s1 string,s2 string){
	n1,_ := db.NodesByLabel(s1)
	n2,_ := db.NodesByLabel("Fruit")
n1[0].Relate("eats", n2[0].Id(), neoism.Props{})
}



func queryAllNodes() {
	// query results
	res := []struct {
		N neoism.Node // Column "n" gets automagically unmarshalled into field N
	}{}

	// construct query
	cq := neoism.CypherQuery{
		Statement: "MATCH (n) RETURN n",
		Result:    &res,
	}
	// execute query
	err := db.Cypher(&cq)
	panicErr(err)

	fmt.Printf("queryAllNodes(%d)\n", len(res))
	for i, _ := range res {
		n := res[i].N // Only one row of data returned
		fmt.Printf("  Node[%d] %+v\n", i, n.Data)
	}

}


func SearchCountryNode(country string) {
	dbn, _ := neoism.Connect("http://neo4j:12345@localhost:7474")
	res1 := []struct {
	    A string `json:"c.cid"`
	}{}
	cq1 := neoism.CypherQuery{
	    Statement: `
	      MATCH (c:Country {CountryName: {CountryName}})
	       RETURN c.cid
	    `,
	    Parameters: neoism.Props{"CountryName": country},
	    Result:     &res1,
	}
	dbn.Cypher(&cq1)
	fmt.Println(len(res1))
	fmt.Println(res1[0].A)
}

func top_five_by_rating(rating string){
	dbn, _ := neoism.Connect("http://neo4j:12345@localhost:7474")
	res := []struct {
			HotelN string `json:"h.hotel_name"`
		}{}
stmt := `MATCH (h:HotelInfo)-[:has_rating]->(r:Rating) WHERE r.rating="`+rating+`" RETURN h.hotel_name ORDER BY  h.discounted_price DESC LIMIT 5`
		cq := neoism.CypherQuery{
			Statement:  stmt,
			Result:     &res,
		}
		dbn.Cypher(&cq)
		for i:=0 ; i<len(res);i++{
		fmt.Printf(res[i].HotelN+"\n")
}
}


func top_five_by_price(){

	dbn, _ := neoism.Connect("http://neo4j:12345@localhost:7474")

	res := []struct {
			HotelN string `json:"h.hotel_name"`
			HotelP string `json:"h.discounted_price"`
		}{}

		cq := neoism.CypherQuery{
			Statement:  `MATCH (h:HotelInfo) RETURN h.discounted_price,h.hotel_name
ORDER BY  h.discounted_price DESC
LIMIT 5`,
			Result:     &res,
		}

		dbn.Cypher(&cq)

		for i:=0 ; i<len(res);i++{
		fmt.Printf(res[i].HotelN + "			"+res[i].HotelP+"\n")
}
}


func hotel_has_amenities(amenity string) {
dbn, _ := neoism.Connect("http://neo4j:12345@localhost:7474")
res := []struct {
		Prince string `json:"h.hotel_name"`
	}{}
stmt := `MATCH (h:HotelInfo)-[:has_amenities]->(ft:facilities) WHERE ft.facility="`+amenity+`" RETURN h.hotel_name`
	cq := neoism.CypherQuery{
		Statement:  stmt,
		Result:     &res,
	}
	dbn.Cypher(&cq)
	fmt.Println(len(res))
	fmt.Println(res)
}


func top_five_by_offer(){

	dbn, _ := neoism.Connect("http://neo4j:12345@localhost:7474")

	res := []struct {
			HotelN string `json:"h.hotel_name"`
			HotelP string `json:"h.offer"`
		}{}

		cq := neoism.CypherQuery{
			Statement:  `MATCH (h:HotelInfo) RETURN h.hotel_name,h.offer ORDER BY(h.offer) DESC LIMIT 5`,
			Result:     &res,
		}

		dbn.Cypher(&cq)
		for i:=0 ; i<len(res);i++{
		fmt.Printf(res[i].HotelN + "			"+res[i].HotelP+"\n")
}
}


func top_five_by_price_rating(){
	dbn, _ := neoism.Connect("http://neo4j:12345@localhost:7474")

	res := []struct {
			HotelN string `json:"h.hotel_name"`
		}{}

		cq := neoism.CypherQuery{
			Statement:  `MATCH (h:HotelInfo)-[rt:has_rating]->(r:Rating) WHERE r.rating = "5" RETURN h.hotel_name ORDER BY(h.discounted_price) DESC LIMIT 5`,
			Result:     &res,
		}

		dbn.Cypher(&cq)
		for i:=0 ; i<len(res);i++{
		fmt.Println(res[i].HotelN)
	}

}




func NearestHotel(lat string, long string, cityid string) {
	la, _ := strconv.ParseFloat(lat, 64)
	lo, _ := strconv.ParseFloat(long, 64)

	Res := []struct {
		Latitude  string `json:"ft.latitude"`
		Longitude string `json:"ft.longitude"`
		Hotelname string `json:"ft.hotel_name"`
	}{}
	cq := neoism.CypherQuery{
		Statement: `MATCH (h:CityInfo)-[:has_hotels]->(ft:HotelInfo)
		            WHERE h.city_id="` + cityid + `"
					RETURN ft.latitude,ft.longitude,ft.hotel_name ORDER BY ft.room_count `,
		Result: &Res,
	}
	err1 := db.Cypher(&cq)
	panicErr(err1)
	nearest := 100000.00
	ans := "abc"
	for i := range Res {
		lah, _ := strconv.ParseFloat(Res[i].Latitude, 64)
		loh, _ := strconv.ParseFloat(Res[i].Longitude, 64)
		dis := ((la-lah)*(la-lah) + (lo-loh)*(lo-loh))
		if dis <= nearest {
			nearest = dis
			ans = Res[i].Hotelname
		}
	}
	fmt.Println("Nearest Hotel-->")
	fmt.Println(ans)
}





func NearestCity(lat string, long string){
	la, _ := strconv.ParseFloat(lat, 64)
	lo, _ := strconv.ParseFloat(long, 64)

	Res := []struct {
		Cityname string `json:"c.city_name"`
		Latitude string `json:"l.la"`
		Longitude string `json:"l.lon"`
	}{}
	cq := neoism.CypherQuery{
		Statement: `MATCH (c:CityInfo)-[:city_lat_long]->(l:LATLONG) RETURN c.city_name,l.la,l.lon`,
		Result: &Res,
	}
	err1 := db.Cypher(&cq)
	panicErr(err1)
	nearest := 100000.00
	ans := "abc"
	for i := range Res {
		lah, _ := strconv.ParseFloat(Res[i].Latitude, 64)
		loh, _ := strconv.ParseFloat(Res[i].Longitude, 64)
		dis := ((la-lah)*(la-lah) + (lo-loh)*(lo-loh))
		if dis <= nearest {
			nearest = dis
			ans = Res[i].Cityname
		}
	}
	fmt.Println("Nearest City-->")
	fmt.Println(ans)
	//return ans
}




func createRelation(node *neoism.Node, nody *neoism.Node, distance int) {
	rels, _ := node.Relationships("target")

	for _, current := range rels {
		relation, _ := current.End()

		if nody.Id() == relation.Id() {
			return
		}
	}

	node.Relate("target", nody.Id(), neoism.Props{"value": distance})
}
