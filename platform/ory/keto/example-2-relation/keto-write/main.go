package main

import (
	"context"
	"fmt"
	"log"

	keto "github.com/ory/keto-client-go"
)

var (
	writeApiURL = "http://localhost:4467" // The URL for the Keto write API
)

func main() {
	client := keto.NewAPIClient(&keto.Configuration{
		Servers: keto.ServerConfigurations{
			{URL: writeApiURL},
		},
	})

	// Define the ABAC relations for the back-office system
	relations := []keto.CreateRelationshipBody{
		// Users Namespace: Defining roles in the "users" namespace
		{
			Namespace: keto.PtrString("users"),
			Object:    keto.PtrString("manager"),
			Relation:  keto.PtrString("member"),
			SubjectId: keto.PtrString("user:rezky"),
		},
		{
			Namespace: keto.PtrString("users"),
			Object:    keto.PtrString("employee"),
			Relation:  keto.PtrString("member"),
			SubjectId: keto.PtrString("user:meita"),
		},
		{
			Namespace: keto.PtrString("users"),
			Object:    keto.PtrString("employee"),
			Relation:  keto.PtrString("member"),
			SubjectId: keto.PtrString("user:dio"),
		},

		// User-Unit Relations
		{
			Namespace: keto.PtrString("units"),
			Object:    keto.PtrString("unit:accounting"),
			Relation:  keto.PtrString("member"),
			SubjectId: keto.PtrString("user:rezky"),
		},
		{
			Namespace: keto.PtrString("units"),
			Object:    keto.PtrString("unit:accounting"),
			Relation:  keto.PtrString("role"),
			SubjectSet: &keto.SubjectSet{
				Namespace: "users",
				Object:    "user:rezky",
				Relation:  "manager",
			},
		},
		{
			Namespace: keto.PtrString("units"),
			Object:    keto.PtrString("unit:budgeting"),
			Relation:  keto.PtrString("member"),
			SubjectId: keto.PtrString("user:dio"),
		},
		{
			Namespace: keto.PtrString("units"),
			Object:    keto.PtrString("unit:budgeting"),
			Relation:  keto.PtrString("role"),
			SubjectSet: &keto.SubjectSet{
				Namespace: "users",
				Object:    "user:dio",
				Relation:  "employee",
			},
		},

		// Module Access for Unit Based on Role
		{
			Namespace: keto.PtrString("modules"),
			Object:    keto.PtrString("module:payment"),
			Relation:  keto.PtrString("manage"),
			SubjectSet: &keto.SubjectSet{
				Namespace: "units",
				Object:    "unit:accounting",
				Relation:  "manager",
			},
		},
		{
			Namespace: keto.PtrString("modules"),
			Object:    keto.PtrString("module:payment"),
			Relation:  keto.PtrString("view"),
			SubjectSet: &keto.SubjectSet{
				Namespace: "units",
				Object:    "unit:accounting",
				Relation:  "employee",
			},
		},
		{
			Namespace: keto.PtrString("modules"),
			Object:    keto.PtrString("module:ledger"),
			Relation:  keto.PtrString("manage"),
			SubjectSet: &keto.SubjectSet{
				Namespace: "units",
				Object:    "unit:budgeting",
				Relation:  "manager",
			},
		},
		{
			Namespace: keto.PtrString("modules"),
			Object:    keto.PtrString("module:ledger"),
			Relation:  keto.PtrString("view"),
			SubjectSet: &keto.SubjectSet{
				Namespace: "units",
				Object:    "unit:budgeting",
				Relation:  "employee",
			},
		},
	}

	ctx := context.Background()
	// Iterate over the relations and create them using the Keto SDK
	for _, relation := range relations {
		_, _, err := client.RelationshipApi.CreateRelationship(ctx).CreateRelationshipBody(relation).Execute()
		if err != nil {
			log.Printf("Error creating relation: %v", err)
		} else {
			fmt.Printf("Successfully created relation: %+v\n", relation)
		}
	}
}
