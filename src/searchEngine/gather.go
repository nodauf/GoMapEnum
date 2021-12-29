package searchengine

import (
	"GoMapEnum/src/utils"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// SEARCH_ENGINE contains url for search on Google and Bing
var SEARCH_ENGINE = map[string]string{"google": `https://www.google.com/search?q=site:linkedin.com/in+"%s"&start=%d&num=100`,
	"bing": `http://www.bing.com/search?q=site:linkedin.com/in+"%s"&first=%d`}

// REGEX_TITLE is the regex to extract title from search engine's results
var REGEX_TITLE = `<h[23](.*?")?>(.*?)<\/h[23]>`

// REGEX_LINKEDIN is the regex to extract field from the title
var REGEX_LINKEDIN = map[string]string{"google": `<h[23](.*?")?>(?P<FirstName>.*?) (?P<LastName>.*?) [-–] (?P<Title>.*?) [-–] (?P<Company>.*?)(\| LinkedIn)(.*?)<\/h[23]>`,
	"bing": `<h[23](.*?")?>(?P<FirstName>.*?) (?P<LastName>.*?) [-–] (?P<Title>.*?) [-–] (?P<Company>.*?)(\| LinkedIn)(.*?)<\/h[23]>`}

// Gather will search a company name and returned the list of people in specified format
func (options *Options) Gather() []string {
	var output []string
	var nbFoundEmail int
	var searchEngineToUse []string
	log = options.Log
	// Always insensitive case comparaison
	options.Company = strings.ToLower(options.Company)
	if options.SearchEngine != "" {
		searchEngineToUse = strings.Split(options.SearchEngine, ",")
	} else {
		searchEngineToUse = utils.GetKeysMap(SEARCH_ENGINE)
	}
	// For specififed search engine
	for _, searchEngine := range searchEngineToUse {
		formatUrl := SEARCH_ENGINE[searchEngine]
		// Reset the variables values
		var startSearch = 0

		// Compile the regex. May be different for each search engine
		reTitle := regexp.MustCompile(REGEX_TITLE)
		reData := regexp.MustCompile(REGEX_LINKEDIN[searchEngine])

		// As long as we have results, we continue
		for {
			// Reset the number to 0
			nbFoundEmail = 0
			log.Target = searchEngine
			log.Verbose("Searching on " + searchEngine + " about " + options.Company + " starting at result " + strconv.Itoa(startSearch))
			url := fmt.Sprintf(formatUrl, options.Company, startSearch)
			log.Debug("URL: " + url)
			// Get the results of the search
			body, statusCode, err := utils.GetBodyInWebsite(url, options.Proxy, nil)
			if err != nil {
				log.Error(err.Error())
				continue
			}
			// Too many requests. No results returned
			if statusCode == 429 {
				log.Error("Too many requests. No results returned. The IP may be blocked")
				break
			}
			// Extract all links of the body
			links := reTitle.FindAllString(body, -1)
			for _, link := range links {
				// Extract all the title
				result := utils.ReSubMatchMap(reData, link)
				// Remove <strong> tag. Bing add it for the searched keyword
				result = utils.SearchReplaceMap(result, "<strong>", "")
				result = utils.SearchReplaceMap(result, "</strong>", "")

				// Compare the company name with case insensitive
				companyName := strings.Trim(strings.ToLower(result["Company"]), " ")
				// Check for the company name (exact match or not)
				if (!options.ExactMatch && strings.Contains(companyName, options.Company)) || (options.ExactMatch && companyName == options.Company) {

					var email string
					email = options.Format
					log.Verbose(result["FirstName"] + " - " + result["LastName"] + " - " + result["Title"] + " - " + result["Company"])
					// output with the specified format
					email = strings.ReplaceAll(email, "{first}", result["FirstName"])
					email = strings.ReplaceAll(email, "{f}", result["FirstName"][0:1])
					email = strings.ReplaceAll(email, "{last}", result["LastName"])
					email = strings.ReplaceAll(email, "{l}", result["LastName"][0:1])
					email = strings.ToLower(email)
					log.Success(email)
					output = append(output, email)
					nbFoundEmail += 1
				}
			}
			// If not result break and go to the next searchEngine
			if nbFoundEmail == 0 {
				break
			}
			startSearch += len(links)
		}

	}
	output = utils.UniqueSliceString(output)
	log.Verbose("Total emails collected: " + strconv.Itoa(len(output)))
	return output
}
