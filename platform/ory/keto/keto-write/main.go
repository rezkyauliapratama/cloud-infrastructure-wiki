package main

import (
	"context"
	"fmt"
	"log"

	keto "github.com/ory/keto-client-go"
)

var (
	writeApiURL = "http://localhost:4467" // URL for Ory Keto Write API
)

func main() {
	// Initialize the Keto Write API client
	writeClient := keto.NewAPIClient(&keto.Configuration{
		Servers: keto.ServerConfigurations{
			{URL: writeApiURL},
		},
	})

	ctx := context.Background()

	// Create relations for EA group
	createRelation(ctx, writeClient, "group:EA", "manage", "folder:EA")

	// Create relations for PMO group
	createRelation(ctx, writeClient, "group:PMO", "manage", "folder:PMO")

	// Create relations for IT Strategy group
	createRelation(ctx, writeClient, "group:ITStrategy", "manage", "folder:ITStrategy")
	createRelation(ctx, writeClient, "group:ITStrategy", "view", "folder:EA")
	createRelation(ctx, writeClient, "group:ITStrategy", "view", "folder:PMO")

	// Create relations for Group Head (GH) group
	createRelation(ctx, writeClient, "group:GH", "manage", "folder:EA")
	createRelation(ctx, writeClient, "group:GH", "manage", "folder:PMO")
	createRelation(ctx, writeClient, "group:GH", "manage", "folder:ITStrategy")

	// Add users to their respective groups
	addMemberToGroup(ctx, writeClient, "user:rezky", "group:EA")
	addMemberToGroup(ctx, writeClient, "user:meita", "group:EA")
	addMemberToGroup(ctx, writeClient, "user:dio", "group:PMO")
	addMemberToGroup(ctx, writeClient, "user:vicky", "group:ITStrategy")
	addMemberToGroup(ctx, writeClient, "user:ardian", "group:GH")

	fmt.Println("Groups and relationships created successfully.")
}

// Helper function to create a relation in Ory Keto
func createRelation(ctx context.Context, client *keto.APIClient, subject, relation, object string) {
	tuple := keto.CreateRelationshipBody{
		Namespace: keto.PtrString("folders"),
		Object:    keto.PtrString(object),
		Relation:  keto.PtrString(relation),
		SubjectSet: &keto.SubjectSet{
			Namespace: "groups",
			Object:    subject,
		},
	}

	_, r, err := client.RelationshipApi.CreateRelationship(ctx).CreateRelationshipBody(tuple).Execute()
	if err != nil {
		log.Fatalf("Failed to create relation: %v", err)
	}

	fmt.Printf("Created relation: %s can %s %s\nfull response %v\n", subject, relation, object, r)
}

// Helper function to add a user to a group
func addMemberToGroup(ctx context.Context, client *keto.APIClient, user, group string) {
	tuple := keto.CreateRelationshipBody{
		Namespace: keto.PtrString("groups"),
		Object:    keto.PtrString(group),
		Relation:  keto.PtrString("member"),
		SubjectId: keto.PtrString(user),
		SubjectSet: &keto.SubjectSet{
			Namespace: "users",
			Object:    user,
		},
	}

	_, r, err := client.RelationshipApi.CreateRelationship(ctx).CreateRelationshipBody(tuple).Execute()
	if err != nil {
		log.Fatalf("Failed to add user to group: %v", err)
	}

	fmt.Printf("Added user: %s to group: %s\nfull response %v\n", user, group, r)
}
