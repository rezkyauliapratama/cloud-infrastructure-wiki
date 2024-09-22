package main

import (
	"context"
	"fmt"
	"log"

	keto "github.com/ory/keto-client-go"
)

var (
	readApiURL = "http://localhost:4466" // Adjust to your Ory Keto Read API URL
)

func main() {
	// Initialize the Keto Read API client
	clientConfig := keto.NewConfiguration()
	clientConfig.Servers = keto.ServerConfigurations{
		{URL: readApiURL},
	}
	readClient := keto.NewAPIClient(clientConfig)

	ctx := context.TODO()
	fmt.Print("connection establish")
	// Read all relations in the "folders" namespace
	readRelations(ctx, readClient, "folders", "folder:EA", "manage", "user:rezky")
}

// Function to read relations using the Keto Read API
func readRelations(ctx context.Context, client *keto.APIClient, namespace, object, relation, subject string) {
	// Make the request to the Read API
	request := client.RelationshipApi.GetRelationships(ctx).Namespace(namespace).Object(object).Relation(relation).SubjectId(subject)

	response, _, err := request.Execute()
	if err != nil {
		log.Fatalf("Error reading relations: %v", err)
	}
	fmt.Printf("has relation ? %v", response.HasRelationTuples())
	// Print the found relations
	fmt.Printf("tuples ? %v", response.GetRelationTuples())

	for _, tuple := range response.RelationTuples {
		object := tuple.GetObject()
		relation := tuple.GetRelation()
		subject := tuple.GetSubjectSet().Object

		fmt.Printf("Relation found: %s#%s@%s\n", object, relation, subject)

	}
}
