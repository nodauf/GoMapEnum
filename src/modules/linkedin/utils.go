package linkedin

import (
	"GoMapEnum/src/utils"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	unidecode "github.com/mozillazg/go-unidecode"
)

// getCompanyInfo return a struct that contains detailed information on a company
func (options *Options) getCompanyInfo() (linkedinGetCompany, error) {
	var company linkedinGetCompany

	linkedinURL := fmt.Sprintf(LINKEDIN_GET_COMPANY_INFO, options.CompanyID)
	header := make(map[string]string)
	header["csrf-token"] = "ajax:1337"
	header["x-restli-protocol-version"] = "2.0.0"
	header["cookie"] = "JSESSIONID=ajax:1337; li_at=" + options.Cookie + ";"
	options.Log.Debug("request to %s", linkedinURL)
	body, statusCode, err := utils.GetBodyInWebsite(linkedinURL, options.ProxyHTTP, header, nil)
	if err != nil {
		return company, err
	}
	if statusCode != 200 {
		log.Debug(body)
		return company, fmt.Errorf("something went wrong. Status code %d != 200. Please run with debug flag for more information", statusCode)
	}
	err = json.Unmarshal([]byte(body), &company)
	if err != nil {
		log.Debug(body)
		return company, fmt.Errorf("Fail to decode the json when requesting %s. Please run with debug flag for more information", linkedinURL)
	}
	return company, nil
}

// getCompanies return a struct that contains all the companies of the research
func (options *Options) getCompanies() (linkedinListCompany, error) {
	var companies linkedinListCompany

	linkedinURL := fmt.Sprintf(LINKEDIN_LIST_COMPANY, url.QueryEscape(options.Company))
	header := make(map[string]string)
	header["csrf-token"] = "ajax:1337"
	header["x-restli-protocol-version"] = "2.0.0"
	header["cookie"] = "JSESSIONID=ajax:1337; li_at=" + options.Cookie + ";"
	options.Log.Debug("request to %s", linkedinURL)
	body, statusCode, err := utils.GetBodyInWebsite(linkedinURL, options.ProxyHTTP, header, nil)
	if err != nil {
		return companies, err
	}
	if statusCode != 200 {
		options.Log.Debug(body)
		return companies, fmt.Errorf("something went wrong. Status code %d != 200. Please run with debug flag for more information", statusCode)
	}
	err = json.Unmarshal([]byte(body), &companies)
	if err != nil {
		log.Debug(body)
		return companies, fmt.Errorf("fail to decode the json when requesting %s. Please run with debug flag for more information: %s", linkedinURL, err.Error())
	}
	return companies, nil
}

// getPeople return a list of people belonging to the company
func (options *Options) getPeople(companyID, start int) ([]string, error) {
	var output []string
	linkedinURL := fmt.Sprintf(LINKEDIN_LIST_PEOPLE, companyID, start)
	header := make(map[string]string)
	header["csrf-token"] = "ajax:1337"
	header["x-restli-protocol-version"] = "2.0.0"
	header["cookie"] = "JSESSIONID=ajax:1337; li_at=" + options.Cookie + ";"
	options.Log.Debug("request to %s", linkedinURL)
	body, statusCode, err := utils.GetBodyInWebsite(linkedinURL, options.ProxyHTTP, header, nil)
	if err != nil {
		return output, err
	}
	if statusCode != 200 {
		options.Log.Debug(body)
		return output, fmt.Errorf("something went wrong. Status code  %d != 200. Please run with debug flag for more information", statusCode)
	}

	//options.Log.Debug(body)
	var peopleStruct linkedinListPeople
	err = json.Unmarshal([]byte(body), &peopleStruct)
	if err != nil {
		options.Log.Debug(body)
		return output, fmt.Errorf("fail to decode the json when requesting %s. Please run with debug flag for more information", linkedinURL)
	}
	//fmt.Printf("%+v", peopleStruct)
	numberPeople := 0
	// The people are in an element of the struct
	for _, element := range peopleStruct.Elements {
		// If the result is empty it is either not the right element for the people or there is no more people
		if element.Items == nil {
			continue
		}

		numberPeople = len(element.Items)
		options.Log.Debug("Found " + strconv.Itoa(numberPeople) + " from " + strconv.Itoa(start) + " for " + options.Company)
		for _, people := range element.Items {
			// if it is an anonymous user, skip it
			if people.ItemUnion.EntityResult.Title.Text == "LinkedIn Member" {
				continue
			}
			// Parse the name to output in the specified format
			name := strings.Split(people.ItemUnion.EntityResult.Title.Text, " ")
			// If the name is composed of more than 2 words or the email should not be guessed, we skip it
			if len(name) == 2 && options.Email {
				var email string
				email = options.Format
				options.Log.Verbose(name[0] + " - " + name[1])
				email = strings.ReplaceAll(email, "{first}", name[0])
				email = strings.ReplaceAll(email, "{f}", name[0][0:1])
				email = strings.ReplaceAll(email, "{last}", name[1])
				email = strings.ReplaceAll(email, "{l}", name[1][0:1])
				email = strings.ToLower(unidecode.Unidecode(email))
				options.Log.Success(email + " - " + people.ItemUnion.EntityResult.PrimarySubtitle.Text + " - " + people.ItemUnion.EntityResult.SecondarySubtitle.Text)
				output = append(output, email)
			}
			if !options.Email {
				result := people.ItemUnion.EntityResult.Title.Text + " - " + people.ItemUnion.EntityResult.PrimarySubtitle.Text + " - " + people.ItemUnion.EntityResult.SecondarySubtitle.Text
				options.Log.Success(result)
				output = append(output, result)

			}

		}
	}
	// If we had people, it means we are not in last page
	if numberPeople > 0 {
		next := start + numberPeople
		peoples, err := options.getPeople(companyID, next)
		if err != nil {
			return output, err
		}
		output = append(output, peoples...)
	}
	return output, nil
}
