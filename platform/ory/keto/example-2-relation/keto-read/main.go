package main

import (
	"context"
	"fmt"
	"log"
	"sync"

	keto "github.com/ory/keto-client-go"
)

var (
	readApiURL = "http://localhost:4466" // URL for the Keto read API
	batchSize  = int64(100)              // Batch size for fetching large datasets
)

func main() {
	// Initialize the Keto Read API client
	clientConfig := keto.NewConfiguration()
	clientConfig.Servers = keto.ServerConfigurations{
		{
			URL: readApiURL,
		},
	}
	client := keto.NewAPIClient(clientConfig)

	// The user to check access for
	userID := "user:rezky"

	// Fetch all relations where rezky has access in the "modules" namespace
	modules, err := getAccessibleModules(client, userID)
	if err != nil {
		log.Fatalf("Error fetching accessible modules: %v", err)
	}

	// Print the accessible modules and their relations
	for module, relations := range modules {
		fmt.Printf("Module: %s, Accessible Relations: %v\n", module, relations)
	}
	// fmt.Printf("namespace: %s, object : %s, relation %s, subject %s\n", tuple.Namespace, tuple.Object, tuple.Relation, *keto.PtrString(*tuple.SubjectId))
	// response, _, err := client.RelationshipApi.GetRelationships(context.Background()).Namespace("units").SubjectId(userID).Execute()

}

// getAccessibleModules fetches modules that a user can access and their access types
func getAccessibleModules(client *keto.APIClient, userID string) (map[string][]string, error) {
	accessibleModules := make(map[string][]string)
	var wg sync.WaitGroup

	// Channels for concurrent data processing
	unitRolesChan := make(chan map[string][]string, 1)
	moduleAccessChan := make(chan map[string][]string, 1)
	errorChan := make(chan error, 2)

	// Step 1: Fetch user roles and unit memberships concurrently
	wg.Add(1)
	go func() {
		defer wg.Done()
		unitRoles, err := fetchUserRolesAndUnits(client, userID)
		if err != nil {
			errorChan <- err
			return
		}
		unitRolesChan <- unitRoles
	}()

	// Step 2: Fetch all module relations concurrently
	wg.Add(1)
	go func() {
		defer wg.Done()
		unitRoles := <-unitRolesChan
		moduleAccess, err := fetchModuleAccess(client, unitRoles)
		if err != nil {
			errorChan <- err
			return
		}
		moduleAccessChan <- moduleAccess
	}()

	// Wait for all goroutines to finish
	wg.Wait()
	close(errorChan)
	close(moduleAccessChan)

	// Check for errors in the error channel
	if len(errorChan) > 0 {
		return nil, <-errorChan
	}

	// Step 3: Collect module access results
	accessibleModules = <-moduleAccessChan

	return accessibleModules, nil
}

// fetchUserRolesAndUnits fetches the roles of the user from the "users" namespace and their unit memberships
func fetchUserRolesAndUnits(client *keto.APIClient, userID string) (map[string][]string, error) {
	unitRoles := make(map[string][]string)
	var user keto.SubjectSet

	response, _, err := client.RelationshipApi.GetRelationships(context.Background()).
		Namespace("users").
		SubjectId(userID).
		Execute()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user roles and units: %w", err)
	}

	// Process user relations to extract roles and units
	for _, relation := range response.RelationTuples {
		fmt.Printf("[fetchUserRolesAndUnits] namespace: %s, object : %s, relation %s, subject %s\n", relation.Namespace, relation.Object, relation.Relation, *keto.PtrString(*relation.SubjectId))

		role := relation.Object

		// Fetch unit roles where this user plays a part
		unitRoleResponse, _, err := client.RelationshipApi.GetRelationships(context.Background()).
			Namespace("units").
			Relation("role").
			SubjectSetNamespace("users").
			SubjectSetObject(userID).
			SubjectSetRelation(role).
			Execute()
		if err != nil {
			log.Printf("Error fetching roles in units for role %s: %v", role, err)
			continue
		}

		// Store unit-role relationships
		for _, unitRoleTuple := range unitRoleResponse.RelationTuples {
			unit := unitRoleTuple.Object
			unitRoles[unit] = append(unitRoles[unit], role)
		}

	}

	fmt.Printf("usesRoles:%v , units: %v\n", unitRoles, user)

	return unitRoles, nil
}

// fetchModuleAccess fetches all module relations and filters based on user roles in units
func fetchModuleAccess(client *keto.APIClient, unitRoles map[string][]string) (map[string][]string, error) {
	accessibleModules := make(map[string][]string)

	// Fetch all module relations
	response, _, err := client.RelationshipApi.GetRelationships(context.Background()).
		Namespace("modules").
		Execute()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch module relations: %w", err)
	}

	// Preprocess and filter based on unit roles
	for _, moduleTuple := range response.RelationTuples {
		for unit, roles := range unitRoles {
			for _, role := range roles {
				if moduleTuple.SubjectSet != nil &&
					moduleTuple.SubjectSet.Namespace == "units" &&
					moduleTuple.SubjectSet.Object == unit &&
					moduleTuple.SubjectSet.Relation == role {
					module := moduleTuple.Object
					accessibleModules[module] = append(accessibleModules[module], moduleTuple.Relation)
				}
			}
		}
	}

	return accessibleModules, nil
}
