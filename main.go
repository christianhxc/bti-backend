package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
)

type ItemInfo struct {
	Plot   string  `json:"plot"`
	Rating float64 `json:"rating"`
}

type Item struct {
	Year  int      `json:"year"`
	Title string   `json:"title"`
	Info  ItemInfo `json:"info"`
}

type Payload struct {
	Movies []Item `json:"movies"`
	Count  int    `json:"count"`
}

func handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	region := "us-west-1"
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region)},
	)

	if err != nil {
		log.Fatal("Got error creating session:", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, "Backend error! Look at the logs for more details.")
	}

	svc := dynamodb.New(sess)

	proj := expression.NamesList(expression.Name("title"), expression.Name("year"), expression.Name("info.rating"))

	expr, err := expression.NewBuilder().WithProjection(proj).Build()

	if err != nil {
		log.Fatal("Got error building expression:", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, "Backend error! Look at the logs for more details.")
	}

	params := &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		ProjectionExpression:      expr.Projection(),
		TableName:                 aws.String("Movies"),
	}

	result, err := svc.Scan(params)

	if err != nil {
		log.Fatal("Query API call failed:", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, "Backend error! Look at the logs for more details.")
	}

	payload := &Payload{
		Movies: nil,
		Count:  len(result.Items),
	}

	for _, i := range result.Items {
		item := Item{}

		err = dynamodbattribute.UnmarshalMap(i, &item)

		if err != nil {
			log.Fatal("Got error unmarshalling:", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			io.WriteString(w, "Backend error! Look at the logs for more details.")
		}

		payload.Movies = append(payload.Movies, item)
	}

	js, err := json.Marshal(payload)
	if err != nil {
		log.Fatal("Got error unmarshalling:", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, "Backend error! Look at the logs for more details.")
	}

	w.Write(js)
	log.Print("Done getting the list of movies ...")
}

func main() {
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
