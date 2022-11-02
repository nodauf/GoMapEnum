package linkedin

import (
	"strconv"
	"strings"
)

// LINKEDIN_LIST_COMPANY is the URL to search for companies
var LINKEDIN_LIST_COMPANY = "https://www.linkedin.com/voyager/api/voyagerSearchDashTypeahead?decorationId=com.linkedin.voyager.dash.deco.search.typeahead.GlobalTypeaheadCollection-27&q=globalTypeahead&query=%s"

// LINKEDIN_LIST_PEOPLE is the URL to search for people
var LINKEDIN_LIST_PEOPLE = "https://www.linkedin.com/voyager/api/search/dash/clusters?decorationId=com.linkedin.voyager.dash.deco.search.SearchClusterCollection-126&origin=COMPANY_PAGE_CANNED_SEARCH&q=all&query=(flagshipSearchIntent:SEARCH_SRP,queryParameters:(currentCompany:List(%d),resultType:List(PEOPLE)),includeFiltersInResponse:false)&start=%d"

// LINKEDIN_GET_COMPANY_INFO is the URL to get information about the company
var LINKEDIN_GET_COMPANY_INFO = "https://www.linkedin.com/voyager/api/entities/companies/%d"

// Gather return a list of all users belongings to the specified company (or multiple companies if the search return more than one)
func (options *Options) Gather() string {
	var output []string
	log = options.Log
	if options.CompanyID != 0 {
		companyInfo := options.getCompanyInfo(int(options.CompanyID))
		options.Log.Verbose("Company name: %s Website: %s Description: %s", companyInfo.BasicCompanyInfo.MiniCompany.Name, companyInfo.WebsiteURL, companyInfo.Description)

		output = options.getPeople(int(options.CompanyID), 0)
		log.Debug("Found " + strconv.Itoa(len(output)) + " peoples for " + companyInfo.BasicCompanyInfo.MiniCompany.Name)
	} else {
		// Always insensitive case compare
		options.Company = strings.ToLower(options.Company)
		// Get all the companies from the option
		companies := options.getCompanies()
		log.Debug("Found " + strconv.Itoa(len(companies.Elements)) + " companies matching " + options.Company)
		for _, company := range companies.Elements {
			// Extract the company name from the struct
			companyLinkedinName := strings.ToLower(company.EntityLockupView.Title.Text)
			log.Debug("Checking for " + companyLinkedinName)
			// If the ID is not empty and check for the company name (exact match or not)
			if company.EntityLockupView.TrackingUrn != "" && (!options.ExactMatch && strings.Contains(companyLinkedinName, options.Company) || options.ExactMatch && companyLinkedinName == options.Company) {
				log.Verbose("Company name: " + companyLinkedinName + " match")
				companyID, _ := strconv.Atoi(strings.Split(company.EntityLockupView.TrackingUrn, ":")[3])
				companyInfo := options.getCompanyInfo(companyID)
				options.Log.Debug("Company name: %s Website: %s Description: %s", companyInfo.BasicCompanyInfo.MiniCompany.Name, companyInfo.WebsiteURL, companyInfo.Description)

				// Get the people of the company, starting from 0
				output = append(output, options.getPeople(companyID, 0)...)
				log.Debug("Found " + strconv.Itoa(len(output)) + " peoples for " + options.Company)
			}
		}
	}

	return strings.Join(output, "\n")
}
