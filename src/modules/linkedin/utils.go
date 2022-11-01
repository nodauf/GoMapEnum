package linkedin

import (
	"GoMapEnum/src/utils"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/mozillazg/go-unidecode"
)

// getCompanyInfo return a struct that contains detailed information on a company
func (options *Options) getCompanyInfo(compagnyID int) linkedinGetCompany {
	var company linkedinGetCompany

	linkedinURL := fmt.Sprintf(LINKEDIN_GET_COMPANY_INFO, compagnyID)
	header := make(map[string]string)
	header["csrf-token"] = "ajax:1337"
	header["x-restli-protocol-version"] = "2.0.0"
	header["cookie"] = "JSESSIONID='ajax:1337'; li_at=" + options.Cookie + ";"
	body, statusCode, err := utils.GetBodyInWebsite(linkedinURL, options.ProxyHTTP, header)
	if err != nil {
		if strings.Contains(err.Error(), "stopped after 10 redirects") {
			log.Error("The session cookie may be wrong")
		}
		log.Error(err.Error())
		return company
	}
	if statusCode != 200 {
		log.Error("Something went wrong. Status code " + strconv.Itoa(statusCode) + " != 200. Body: " + body)
		return company
	}
	err = json.Unmarshal([]byte(body), &company)
	if err != nil {
		log.Error("Fail to decode the json when requesting %s. Please run with debug flag for more information", linkedinURL)
		log.Debug(body)
	}
	return company
}

// getCompanies return a struct that contains all the companies of the research
func (options *Options) getCompanies() linkedinListCompany {
	var companies linkedinListCompany

	linkedinURL := fmt.Sprintf(LINKEDIN_LIST_COMPANY, url.QueryEscape(options.Company))
	header := make(map[string]string)
	header["csrf-token"] = "ajax:1337"
	header["x-restli-protocol-version"] = "2.0.0"
	header["cookie"] = "JSESSIONID='ajax:1337'; li_at=" + options.Cookie + ";"
	body, statusCode, err := utils.GetBodyInWebsite(linkedinURL, options.ProxyHTTP, header)
	if err != nil {
		if strings.Contains(err.Error(), "stopped after 10 redirects") {
			log.Error("The session cookie may be wrong")
		}
		log.Error(err.Error())
		return companies
	}
	if statusCode != 200 {
		log.Error("Something went wrong. Status code " + strconv.Itoa(statusCode) + " != 200. Body: " + body)
		return companies
	}
	err = json.Unmarshal([]byte(body), &companies)
	if err != nil {
		log.Error("Fail to decode the json when requesting %s. Please run with debug flag for more information", linkedinURL)
		log.Debug(body)
	}
	return companies
}

// getPeople return a list of people belonging to the company
func (options *Options) getPeople(companyID, start int) []string {
	var output []string
	linkedinURL := fmt.Sprintf(LINKEDIN_LIST_PEOPLE, companyID, start)
	header := make(map[string]string)
	header["csrf-token"] = "ajax:1337"
	header["x-restli-protocol-version"] = "2.0.0"
	header["cookie"] = "JSESSIONID='ajax:1337'; li_at=" + options.Cookie + ";"

	body, statusCode, err := utils.GetBodyInWebsite(linkedinURL, options.ProxyHTTP, header)
	if err != nil {
		log.Error(err.Error())
		return output
	}
	if statusCode != 200 {
		log.Error("Something went wrong. Status code " + strconv.Itoa(statusCode) + " != 200. Body: " + body)
		return output
	}
	var peopleStruct linkedinListPeople
	err = json.Unmarshal([]byte(body), &peopleStruct)
	if err != nil {
		log.Error("Fail to decode the json when requesting %s. Please run with debug flag for more information", linkedinURL)
		log.Debug(body)
	}
	numberPeople := 0
	// The people are in an element of the struct
	for _, element := range peopleStruct.Elements {
		// If the result is empty it is either not the right element for the people or there is no more people
		if element.Results == nil {
			continue
		}

		numberPeople = len(element.Results)
		log.Debug("Found " + strconv.Itoa(numberPeople) + " from " + strconv.Itoa(start) + " for " + options.Company)
		for _, people := range element.Results {
			// if it is an anonymous user, skip it
			if people.Title.Text == "LinkedIn Member" {
				continue
			}
			// Parse the name to output in the specified format
			name := strings.Split(people.Title.Text, " ")
			// If the name is composed of more than 2 words or the email should not be guessed, we skip it
			if len(name) == 2 && options.Email {
				var email string
				email = options.Format
				log.Verbose(name[0] + " - " + name[1])
				email = strings.ReplaceAll(email, "{first}", name[0])
				email = strings.ReplaceAll(email, "{f}", name[0][0:1])
				email = strings.ReplaceAll(email, "{last}", name[1])
				email = strings.ReplaceAll(email, "{l}", name[1][0:1])
				email = strings.ToLower(unidecode.Unidecode(email))
				log.Success(email + " - " + people.PrimarySubtitle.Text + " - " + people.SecondarySubtitle.Text)
				output = append(output, email)
			}
			if !options.Email {
				result := people.Title.Text + " - " + people.PrimarySubtitle.Text + " - " + people.SecondarySubtitle.Text
				log.Success(result)
				output = append(output, result)

			}

		}
	}
	// If we had people, it means we are not in last page
	if numberPeople > 0 {
		next := start + numberPeople
		output = append(output, options.getPeople(companyID, next)...)
	}
	return output
}
