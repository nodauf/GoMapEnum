package linkedin

import (
	"strconv"
	"strings"
)

// URL to search for companies
var LINKEDIN_LIST_COMPANY = "https://www.linkedin.com/voyager/api/voyagerSearchDashTypeahead?decorationId=com.linkedin.voyager.dash.deco.search.typeahead.GlobalTypeaheadCollection-27&q=globalTypeahead&query=%s"

// URL to search for people
var LINKEDIN_LIST_PEOPLE = "https://www.linkedin.com/voyager/api/search/dash/clusters?decorationId=com.linkedin.voyager.dash.deco.search.SearchClusterCollection-126&origin=COMPANY_PAGE_CANNED_SEARCH&q=all&query=(flagshipSearchIntent:SEARCH_SRP,queryParameters:(currentCompany:List(%d),resultType:List(PEOPLE)),includeFiltersInResponse:false)&start=%d"

// Gather return a list of all users belongins to the specified company (or multiple companies if the search return more than one)
func (options *Options) Gather() []string {
	var output []string
	log = options.Log
	// Always insensitive case compare
	options.Company = strings.ToLower(options.Company)
	// Get all the companies from the option
	companies := options.getCompany()
	for _, company := range companies.Elements {
		// Extract the company name from the struct
		companyLinkedinName := strings.ToLower(company.EntityLockupView.Title.Text)
		// If the ID is not empty and check for the company name (exact match or not)
		if company.EntityLockupView.TrackingUrn != "" && (!options.ExactMatch && strings.Contains(companyLinkedinName, options.Company) || options.ExactMatch && companyLinkedinName == options.Company) {
			log.Debug("Company name: " + companyLinkedinName)
			companyID, _ := strconv.Atoi(strings.Split(company.EntityLockupView.TrackingUrn, ":")[3])
			// Get the people of the company, starting from 0
			output = options.getPeople(companyID, 0)
		}
	}

	return output
}
