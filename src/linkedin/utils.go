package linkedin

import (
	"GoMapEnum/src/utils"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

// getCompany return a struct that contains all the company of the research
func (options *Options) getCompany() linkedinListCompany {
	var companies linkedinListCompany

	url := fmt.Sprintf(LINKEDIN_LIST_COMPANY, url.QueryEscape(options.Company))
	header := make(map[string]string)
	header["csrf-token"] = "ajax:1337"
	header["x-restli-protocol-version"] = "2.0.0"
	header["cookie"] = "JSESSIONID='ajax:1337'; li_at=" + options.Cookie + ";"
	body, statusCode, err := utils.GetBodyInWebsite(url, options.Proxy, header)
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
	json.Unmarshal([]byte(body), &companies)
	return companies
}

// getPeople return a list of people belonging to the company
func (options *Options) getPeople(companyID, start int) []string {
	var output []string
	url := fmt.Sprintf(LINKEDIN_LIST_PEOPLE, companyID, start)
	header := make(map[string]string)
	header["csrf-token"] = "ajax:1337"
	header["x-restli-protocol-version"] = "2.0.0"
	header["cookie"] = "JSESSIONID='ajax:1337'; li_at=" + options.Cookie + ";"

	body, statusCode, err := utils.GetBodyInWebsite(url, options.Proxy, header)
	if err != nil {
		log.Error(err.Error())
		return output
	}
	if statusCode != 200 {
		log.Error("Something went wrong. Status code " + strconv.Itoa(statusCode) + " != 200. Body: " + body)
		return output
	}
	var peopleStruct linkedinListPeople
	json.Unmarshal([]byte(body), &peopleStruct)
	numberPeople := 0
	// The people are in an element of the struct
	for _, element := range peopleStruct.Elements {
		// If the result if empty it is either not the right element for the people or there is no more people
		if element.Results == nil {
			continue
		}

		numberPeople = len(element.Results)
		for _, people := range element.Results {
			// if it is an anonymous user, skip it
			if people.Title.Text == "LinkedIn Member" {
				continue
			}
			// Parse the name to output in the specified format
			name := strings.Split(people.Title.Text, " ")
			// If the name is composed of more than 2 words, we skip it
			if len(name) == 2 {
				var email string
				email = options.Format
				log.Verbose(name[0] + " - " + name[1])
				email = strings.ReplaceAll(email, "{first}", name[0])
				email = strings.ReplaceAll(email, "{f}", name[0][0:1])
				email = strings.ReplaceAll(email, "{last}", name[1])
				email = strings.ReplaceAll(email, "{l}", name[1][0:1])
				email = strings.ToLower(email)
				log.Success(email + " - " + people.PrimarySubtitle.Text + "-" + people.SecondarySubtitle.Text)
				output = append(output, email)
			}

		}
	}
	// If we had people, it means we are not in last page
	if numberPeople > 0 {
		next := start + numberPeople
		output = append(options.getPeople(companyID, next), output...)
	}
	return output
}
