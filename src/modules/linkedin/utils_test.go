package linkedin

import (
	"GoMapEnum/src/logger"
	"flag"
	"strings"
	"testing"
)

var options Options

func init() {
	flag.StringVar(&options.Cookie, "cookie", "", "Session cookie named li_at")

}

// TestGetCompanyInfo test the function getCompanyInfo to check the LinkedIn API is not being updated
func TestGetCompanyInfo(t *testing.T) {
	log := logger.New("test", "linkedin", "Linkedin")
	log.SetLevel(logger.DebugLevel)
	options.Log = log
	options.CompanyID = 11452158
	name := "Contoso"
	websiteURL := "http://www.contoso.org"
	description := "Contoso Ltd. (also known as Contoso and Contoso University) is a fictional company used by Microsoft as an example company and domain."

	if options.Cookie == "" {
		t.Error("Argument cookie is not set. Cannot execute test on linkedin module")
		return
	}
	companyInfo, err := options.getCompanyInfo()
	if err != nil {
		if strings.Contains(err.Error(), "stopped after 10 redirects") {
			t.Errorf("The session cookie may be wrong")
		} else {
			t.Errorf(err.Error())
		}
		return
	}
	if companyInfo.BasicCompanyInfo.MiniCompany.Name != name {
		t.Errorf("Expected name %s not matched %s", name, companyInfo.BasicCompanyInfo.MiniCompany.Name)
	}

	if companyInfo.WebsiteURL != websiteURL {
		t.Errorf("Expected websiteURL %s not matched %s", websiteURL, companyInfo.WebsiteURL)
	}

	if companyInfo.Description != description {
		t.Errorf("Expected description %s not matched %s", description, companyInfo.Description)
	}

}

// TestGetCompanies test the function getCompanies to check the LinkedIn API is not being updated
func TestGetCompanies(t *testing.T) {
	log := logger.New("test", "linkedin", "Linkedin")
	log.SetLevel(logger.DebugLevel)
	options.Log = log
	options.Company = "contos"
	companyName := "contoso"
	companyID := "urn:li:company:11452158"

	if options.Cookie == "" {
		t.Error("Argument cookie is not set. Cannot execute test on linkedin module")
		return
	}

	companies, err := options.getCompanies()
	if err != nil {
		if strings.Contains(err.Error(), "stopped after 10 redirects") {
			t.Errorf("The session cookie may be wrong")
		} else {
			t.Errorf(err.Error())
		}
		return
	}
	var found = false
	for _, company := range companies.Elements {
		companyLinkedinName := strings.ToLower(company.EntityLockupView.Title.Text)

		if companyLinkedinName == companyName && company.EntityLockupView.TrackingUrn == companyID {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("%s not found in the list of %d companies with the search %s", companyName, len(companies.Elements), options.Company)
	}
}

// TestGetPeople test the function getPeople to check the LinkedIn API is not being updated
func TestGetPeople(t *testing.T) {
	log := logger.New("test", "linkedin", "Linkedin")
	log.SetLevel(logger.DebugLevel)
	options.Log = log
	options.Company = "contoso"
	options.Format = "{first}.{last}@contoso.com"
	options.Email = true
	peopleEmail := "kurt.shintaku@contoso.com"
	companyID := 11452158

	if options.Cookie == "" {
		t.Error("Argument cookie is not set. Cannot execute test on linkedin module")
		return
	}

	peoples, err := options.getPeople(companyID, 0)
	if err != nil {
		if strings.Contains(err.Error(), "stopped after 10 redirects") {
			t.Errorf("The session cookie may be wrong")
		} else {
			t.Errorf(err.Error())
		}
		return
	}
	var found = false
	for _, people := range peoples {
		if people == peopleEmail {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("%s not found in the list of %d peoples with the search %d", peopleEmail, len(peoples), companyID)
	}
}
