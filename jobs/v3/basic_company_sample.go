package cjdsample

import (
	"fmt"
	"log"
	"os"
	"time"

	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	talent "google.golang.org/api/jobs/v3"
)

// [START get_parent]

/**
 *
 */
func GetParent() string {
	projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")
	return fmt.Sprintf("projects/%s", projectID)
}

// [END get_parent]

// [START create_service]

/**
 * Create service of Cloud Talent Solution
 */
func CreateCtsService() (*talent.Service, error) {
	// Authorize the client using Application Default Credentials.
	// See https://g.co/dv/identity/protocols/application-default-credentials
	ctx := context.Background()
	client, err := google.DefaultClient(ctx, talent.CloudPlatformScope)
	if err != nil {
		log.Fatalf("Failed to create client of Cloud Talent Solution, Err: %v", err)
	}
	// Create the jobs service client.
	ctsService, err := talent.New(client)
	if err != nil {
		log.Fatalf("Failed to create service of Cloud Talent Solution, Err: %v", err)
	}
	return ctsService, err
}

// [END create_service]

// [START basic_company]

/**
 * Construct a company
 */
func ConstructCompanyWithRequiredFields() *talent.Company {
	externalId := fmt.Sprintf("sample-company-%d", time.Now().UnixNano())
	return &talent.Company{
		ExternalId:  externalId,
		DisplayName: "Google Sample",
	}
}

// [END basic_company]

// [START create_company]

/**
 * Create a company
 */
func CreateCompany(service *talent.Service, companyToCreate *talent.Company) (*talent.Company, error) {
	createCompanyRquest := &talent.CreateCompanyRequest{
		Company: companyToCreate,
	}
	company, err := service.Projects.Companies.Create(GetParent(), createCompanyRquest).Do()
	if err != nil {
		log.Fatalf("Failed to create company %s, Err: %v", companyToCreate.DisplayName, err)
	}

	return company, err
}

// [END create_company]

// [START get_company]

/**
 * Get a company
 */
func GetCompany(service *talent.Service, name string) (*talent.Company, error) {
	company, err := service.Projects.Companies.Get(name).Do()
	if err != nil {
		log.Fatalf("Failed to get company %s, Err: %v", name, err)
	}

	return company, err
}

// [END get_company]

// [START update_company]

/**
 * Update company with all fields
 */
func UpdateCompany(service *talent.Service, name string, companyToUpdate *talent.Company) (*talent.Company, error) {
	updateCompanyRequest := &talent.UpdateCompanyRequest{
		Company: companyToUpdate,
	}
	company, err := service.Projects.Companies.Patch(name, updateCompanyRequest).Do()
	if err != nil {
		log.Fatalf("Failed to update company %s, Err: %v", name, err)
	}

	return company, err
}

// [END update_company]

// [START update_company_with_field_mask]

/**
 * Update company with field mask
 * mask: comma separated top-level fields of Company
 */
func UpdateCompanyWithMask(service *talent.Service, name string, mask string, companyToUpdate *talent.Company) (*talent.Company, error) {
	updateCompanyRequest := &talent.UpdateCompanyRequest{
		Company:    companyToUpdate,
		UpdateMask: mask,
	}
	company, err := service.Projects.Companies.Patch(name, updateCompanyRequest).Do()
	if err != nil {
		log.Fatalf("Failed to update company %s with mask %s, Err: %v", name, mask, err)
	}

	return company, err
}

// [END update_company_with_field_mask]

// [START delete_company]

/**
 * Delete a company
 */
func DeleteCompany(service *talent.Service, name string) (*talent.Empty, error) {
	empty, err := service.Projects.Companies.Delete(name).Do()
	if err != nil {
		log.Fatalf("Failed to delete company %s, Err: %v", name, err)
	}

	return empty, err
}

// [END delete_company

// [START list_companies]

/**
 * List companies in the project
 */
func ListCompanies(service *talent.Service) (*talent.ListCompaniesResponse, error) {
	resp, err := service.Projects.Companies.List(GetParent()).Do()
	if err != nil {
		log.Fatalf("Failed to list companies, Err: %v", err)
	}

	return resp, err
}

// [END list_companies]

// [START basic_company_sample_entry]
func BasicCompanySampleEntry() {
	service, _ := CreateCtsService()

	companyToCreate := ConstructCompanyWithRequiredFields()
	companyCreated, _ := CreateCompany(service, companyToCreate)
	fmt.Printf("CreateCompany: %s\n", companyCreated.DisplayName)

	name := companyCreated.Name
	//name := fmt.Sprintf("%s/companies/%d", GetParent(), rand.Uint64())
	companyGot, _ := GetCompany(service, name)
	fmt.Printf("GetCompany: %s\n", companyGot.DisplayName)

	companyToUpdate := companyCreated
	companyToUpdate.DisplayName = "Google Sample (updated)"
	companyUpdated, _ := UpdateCompany(service, name, companyToUpdate)
	fmt.Printf("UpdateCompany: %s\n", companyUpdated.DisplayName)

	companyUpdated.WebsiteUri = "http://googlesample.com"
	companyUpdated.DisplayName = "Google Sample (updated with mask)"
	companyUpdatedWithMask, _ := UpdateCompanyWithMask(service, name, "WebSiteUri,DisplayName", companyUpdated)
	fmt.Printf("UpdateCompanyWithMask: %s\n", companyUpdatedWithMask.DisplayName)

	empty, _ := DeleteCompany(service, name)
	fmt.Printf("DeleteCompany StatusCode: %d\n", empty.ServerResponse.HTTPStatusCode)

	resp, _ := ListCompanies(service)
	fmt.Printf("ListCompanies Request ID: %q\n", resp.Metadata.RequestId)

	for _, company := range resp.Companies {
		fmt.Printf("-- Company: %q\n", company.Name)
	}

}

// [END basic_company_sample_entry]
