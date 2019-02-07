// Copyright 2018 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package snippets

import (
	"fmt"

	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	iam "google.golang.org/api/iam/v1"
)

var client, _ = google.DefaultClient(
	context.Background(),
	iam.CloudPlatformScope)
var iamService, _ = iam.New(client)

// [START iam_view_grantable_roles]
func viewGrantableRoles(fullResourceName string) {
	request := iam.QueryGrantableRolesRequest{FullResourceName: fullResourceName}
	response, _ := iamService.Roles.QueryGrantableRoles(&request).Do()
	for _, role := range response.Roles {
		fmt.Println("Title: " + role.Title)
		fmt.Println("Name: " + role.Name)
		fmt.Println("Description: " + role.Description)
		fmt.Println()
	}
}

// [END iam_view_grantable_roles]

// [START iam_query_testable_permissions]
func queryTestablePermissions(fullResourceName string) []*iam.Permission {
	request := iam.QueryTestablePermissionsRequest{FullResourceName: fullResourceName}
	response, _ := iamService.Permissions.QueryTestablePermissions(
		&request).Do()
	for _, p := range response.Permissions {
		fmt.Println(p.Name)
	}
	return response.Permissions
}

// [END iam_query_testable_permissions]

// [START iam_get_role]
func getRole(name string) iam.Role {
	role, _ := iamService.Roles.Get(name).Do()
	fmt.Println(role.Name)
	for _, permission := range role.IncludedPermissions {
		fmt.Println(permission)
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
	role, _ := iamService.Projects.Roles.Create(
		"projects/"+projectID, &request).Do()
	fmt.Println("Created role: " + role.Name)
	return *role
}

// [END iam_create_role]

// [START iam_edit_role]
func editRole(name string, projectID string, newTitle string, newDescription string,
	newPermissions []string, newStage string) iam.Role {
	resource := "projects/" + projectID + "/roles/" + name
	role, _ := iamService.Projects.Roles.Get(resource).Do()
	role.Title = newTitle
	role.Description = newDescription
	role.IncludedPermissions = newPermissions
	role.Stage = newStage
	role, _ = iamService.Projects.Roles.Patch(resource, role).Do()
	fmt.Println("Updated role: " + role.Name)
	return *role
}

// [END iam_edit_role]

// [START iam_disable_role]
func disableRole(name string, projectID string) iam.Role {
	resource := "projects/" + projectID + "/roles/" + name
	role, er := iamService.Projects.Roles.Get(resource).Do()
	if er != nil {
		fmt.Println("DID NOT GOT ROLE")
		fmt.Println(er)
	}
	role.Stage = "DISABLED"
	role, erx := iamService.Projects.Roles.Patch(resource, role).Do()
	if erx != nil {
		fmt.Println(er)
	}
	fmt.Println("Disabled role: " + role.Name)
	return *role
}

// [END iam_disable_role]

// [START iam_list_roles]
func listRoles(projectID string) []*iam.Role {
	response, _ := iamService.Projects.Roles.List("projects/" + projectID).Do()
	for _, role := range response.Roles {
		fmt.Println(role.Name)
	}
	return response.Roles
}

// [END iam_list_roles]

// [START iam_delete_role]
func deleteRole(name string, projectID string) {
	iamService.Projects.Roles.Delete(
		"projects/" + projectID + "/roles/" + name).Do()
	fmt.Println("Deleted role: " + name)
}

// [END iam_delete_role]

// [START iam_undelete_role]
func undeleteRole(name string, projectID string) iam.Role {
	resource := "projects/" + projectID + "/roles/" + name
	request := iam.UndeleteRoleRequest{}
	role, _ := iamService.Projects.Roles.Undelete(resource, &request).Do()
	fmt.Println("Undeleted role: " + role.Name)
	return *role
}

// [END iam_undelete_role]

// [START iam_create_service_account]
func createServiceAccount(projectID string, name string, displayName string) iam.ServiceAccount {
	request := iam.CreateServiceAccountRequest{
		AccountId:      name,
		ServiceAccount: &iam.ServiceAccount{DisplayName: displayName}}
	serviceAccount, err := iamService.Projects.ServiceAccounts.Create(
		"projects/"+projectID, &request).Do()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Created service account: " + serviceAccount.Email)
	return *serviceAccount
}

// [END iam_create_service_account]

// [START iam_list_service_accounts]
func listServiceAccounts(projectID string) []*iam.ServiceAccount {
	response, _ := iamService.Projects.ServiceAccounts.List("projects/" + projectID).Do()
	for _, account := range response.Accounts {
		fmt.Println("Name: ", account.Name)
		fmt.Println("Display Name:", account.DisplayName)
		fmt.Println("Email:", account.Email)
		fmt.Println()
	}
	return response.Accounts
}

// [END iam_list_service_accounts]

// [START iam_rename_service_account]
func renameServiceAccount(email string, newDisplayName string) iam.ServiceAccount {
	// First, get a ServiceAccount using List() or Get()
	resource := "projects/-/serviceAccounts/" + email
	serviceAccount, _ := iamService.Projects.ServiceAccounts.Get(resource).Do()

	// Then you can update the display name
	serviceAccount.DisplayName = newDisplayName
	serviceAccount, _ = iamService.Projects.ServiceAccounts.Update(
		resource, serviceAccount).Do()

	fmt.Println("Updated service account: " + serviceAccount.Email)
	return *serviceAccount
}

// [END iam_rename_service_account]

// [START iam_delete_service_account]
func deleteServiceAccount(email string) {
	iamService.Projects.ServiceAccounts.Delete(
		"projects/-/serviceAccounts/" + email).Do()
	fmt.Println("Deleted service account:", email)
}

// [END iam_delete_service_account]

// [START iam_create_key]
func createKey(serviceAccountEmail string) iam.ServiceAccountKey {
	resource := "projects/-/serviceAccounts/" + serviceAccountEmail
	request := iam.CreateServiceAccountKeyRequest{}
	key, err := iamService.Projects.ServiceAccounts.Keys.Create(
		resource, &request).Do()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Created key: ", key.Name)
	return *key
}

// [END iam_create_key]

// [START iam_list_keys]
func listKeys(serviceAccountEmail string) []*iam.ServiceAccountKey {
	resource := "projects/-/serviceAccounts/" + serviceAccountEmail
	response, _ := iamService.Projects.ServiceAccounts.Keys.List(
		resource).Do()
	for _, key := range response.Keys {
		fmt.Println("Key: ", key.Name)
	}
	return response.Keys
}

// [END iam_list_keys]

// [START iam_delete_key]
func deleteKey(fullKeyName string) {
	iamService.Projects.ServiceAccounts.Keys.Delete(fullKeyName).Do()
	fmt.Println("Deleted key: ", fullKeyName)
}

// [END iam_delete_key]
