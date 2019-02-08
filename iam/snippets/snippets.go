// Copyright 2018 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package snippets

import (
	"log"

	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	iam "google.golang.org/api/iam/v1"
)

var client, cerr = google.DefaultClient(context.Background(), iam.CloudPlatformScope)
var iamService, serr = iam.New(client)

// [START iam_view_grantable_roles]
func viewGrantableRoles(fullResourceName string) {
	request := iam.QueryGrantableRolesRequest{FullResourceName: fullResourceName}
	response, err := iamService.Roles.QueryGrantableRoles(&request).Do()
	if err != nil {
		log.Fatalf("Roles.QueryGrantableRoles: %v", err)
	}
	for _, role := range response.Roles {
		log.Println("Title: " + role.Title)
		log.Println("Name: " + role.Name)
		log.Println("Description: " + role.Description)
		log.Println()
	}
}

// [END iam_view_grantable_roles]

// [START iam_query_testable_permissions]
func queryTestablePermissions(fullResourceName string) []*iam.Permission {
	request := iam.QueryTestablePermissionsRequest{FullResourceName: fullResourceName}
	response, err := iamService.Permissions.QueryTestablePermissions(&request).Do()
	if err != nil {
		log.Fatalf("Permissions.QueryTestablePermissions: %v", err)
	}
	for _, p := range response.Permissions {
		log.Println(p.Name)
	}
	return response.Permissions
}

// [END iam_query_testable_permissions]

// [START iam_get_role]
func getRole(name string) iam.Role {
	role, err := iamService.Roles.Get(name).Do()
	if err != nil {
		log.Fatalf("Roles.Get: %v", err)
	}
	log.Println(role.Name)
	for _, permission := range role.IncludedPermissions {
		log.Println(permission)
	}
	return *role
}

// [END iam_get_role]

// [START iam_create_role]
func createRole(name string, projectID string, title string, description string,
	permissions []string, stage string) iam.Role {
	request := iam.CreateRoleRequest{Role: &iam.Role{
		Title:               title,
		Description:         description,
		IncludedPermissions: permissions,
		Stage:               stage,
	}, RoleId: name}
	role, err := iamService.Projects.Roles.Create("projects/"+projectID, &request).Do()
	if err != nil {
		log.Fatalf("Projects.Roles.Create: %v", err)
	}
	log.Println("Created role: " + role.Name)
	return *role
}

// [END iam_create_role]

// [START iam_edit_role]
func editRole(name string, projectID string, newTitle string, newDescription string,
	newPermissions []string, newStage string) iam.Role {
	resource := "projects/" + projectID + "/roles/" + name
	role, err := iamService.Projects.Roles.Get(resource).Do()
	if err != nil {
		log.Fatalf("Projects.Roles.Get: %v", err)
	}
	role.Title = newTitle
	role.Description = newDescription
	role.IncludedPermissions = newPermissions
	role.Stage = newStage
	role, err = iamService.Projects.Roles.Patch(resource, role).Do()
	if err != nil {
		log.Fatalf("Projects.Roles.Patch: %v", err)
	}
	log.Println("Updated role: " + role.Name)
	return *role
}

// [END iam_edit_role]

// [START iam_disable_role]
func disableRole(name string, projectID string) iam.Role {
	resource := "projects/" + projectID + "/roles/" + name
	role, err := iamService.Projects.Roles.Get(resource).Do()
	if err != nil {
		log.Fatalf("Projects.Roles.Get: %v", err)
	}
	role.Stage = "DISABLED"
	role, err = iamService.Projects.Roles.Patch(resource, role).Do()
	if err != nil {
		log.Fatalf("Projects.Roles.Patch: %v", err)
	}
	log.Println("Disabled role: " + role.Name)
	return *role
}

// [END iam_disable_role]

// [START iam_list_roles]
func listRoles(projectID string) []*iam.Role {
	response, err := iamService.Projects.Roles.List("projects/" + projectID).Do()
	if err != nil {
		log.Fatalf("Projects.Roles.List: %v", err)
	}
	for _, role := range response.Roles {
		log.Println(role.Name)
	}
	return response.Roles
}

// [END iam_list_roles]

// [START iam_delete_role]
func deleteRole(name string, projectID string) {
	_, err := iamService.Projects.Roles.Delete("projects/" + projectID + "/roles/" + name).Do()
	if err != nil {
		log.Fatalf("Projects.Roles.Delete: %v", err)
	}
	log.Println("Deleted role: " + name)
}

// [END iam_delete_role]

// [START iam_undelete_role]
func undeleteRole(name string, projectID string) iam.Role {
	resource := "projects/" + projectID + "/roles/" + name
	request := iam.UndeleteRoleRequest{}
	role, err := iamService.Projects.Roles.Undelete(resource, &request).Do()
	if err != nil {
		log.Fatalf("Projects.Roles.Undelete: %v", err)
	}
	log.Println("Undeleted role: " + role.Name)
	return *role
}

// [END iam_undelete_role]

// [START iam_create_service_account]
func createServiceAccount(projectID string, name string, displayName string) iam.ServiceAccount {
	request := iam.CreateServiceAccountRequest{
		AccountId:      name,
		ServiceAccount: &iam.ServiceAccount{DisplayName: displayName}}
	serviceAccount, err := iamService.Projects.ServiceAccounts.Create("projects/"+projectID, &request).Do()
	if err != nil {
		log.Fatalf("Projects.ServiceAccounts.Create: %v", err)
	}
	log.Println("Created service account: " + serviceAccount.Email)
	return *serviceAccount
}

// [END iam_create_service_account]

// [START iam_list_service_accounts]
func listServiceAccounts(projectID string) []*iam.ServiceAccount {
	response, err := iamService.Projects.ServiceAccounts.List("projects/" + projectID).Do()
	if err != nil {
		log.Fatalf("Projects.ServiceAccounts.List: %v", err)
	}
	for _, account := range response.Accounts {
		log.Println("Name: ", account.Name)
		log.Println("Display Name:", account.DisplayName)
		log.Println("Email:", account.Email)
		log.Println()
	}
	return response.Accounts
}

// [END iam_list_service_accounts]

// [START iam_rename_service_account]
func renameServiceAccount(email string, newDisplayName string) iam.ServiceAccount {
	// First, get a ServiceAccount using List() or Get()
	resource := "projects/-/serviceAccounts/" + email
	serviceAccount, err := iamService.Projects.ServiceAccounts.Get(resource).Do()
	if err != nil {
		log.Fatalf("Projects.ServiceAccounts.Get: %v", err)
	}

	// Then you can update the display name
	serviceAccount.DisplayName = newDisplayName
	serviceAccount, err = iamService.Projects.ServiceAccounts.Update(resource, serviceAccount).Do()
	if err != nil {
		log.Fatalf("Projects.ServiceAccounts.Update: %v", err)
	}

	log.Println("Updated service account: " + serviceAccount.Email)
	return *serviceAccount
}

// [END iam_rename_service_account]

// [START iam_delete_service_account]
func deleteServiceAccount(email string) {
	_, err := iamService.Projects.ServiceAccounts.Delete("projects/-/serviceAccounts/" + email).Do()
	if err != nil {
		log.Fatalf("Projects.ServiceAccounts.Delete: %v", err)
	}
	log.Println("Deleted service account:", email)
}

// [END iam_delete_service_account]

// [START iam_create_key]
func createKey(serviceAccountEmail string) iam.ServiceAccountKey {
	resource := "projects/-/serviceAccounts/" + serviceAccountEmail
	request := iam.CreateServiceAccountKeyRequest{}
	key, err := iamService.Projects.ServiceAccounts.Keys.Create(resource, &request).Do()
	if err != nil {
		log.Fatalf("Projects.ServiceAccounts.Keys.Create: %v", err)
	}
	log.Println("Created key: ", key.Name)
	return *key
}

// [END iam_create_key]

// [START iam_list_keys]
func listKeys(serviceAccountEmail string) []*iam.ServiceAccountKey {
	resource := "projects/-/serviceAccounts/" + serviceAccountEmail
	response, err := iamService.Projects.ServiceAccounts.Keys.List(resource).Do()
	if err != nil {
		log.Fatalf("Projects.ServiceAccounts.Keys.List: %v", err)
	}
	for _, key := range response.Keys {
		log.Println("Key: ", key.Name)
	}
	return response.Keys
}

// [END iam_list_keys]

// [START iam_delete_key]
func deleteKey(fullKeyName string) {
	_, err := iamService.Projects.ServiceAccounts.Keys.Delete(fullKeyName).Do()
	if err != nil {
		log.Fatalf("Projects.ServiceAccounts.Keys.Delete: %v", err)
	}
	log.Println("Deleted key: ", fullKeyName)
}

// [END iam_delete_key]
