package utils

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"html/template"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"regexp"
	"strings"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// GetStringOrFile return the content of the file if it is a file otherwise return the string
func GetStringOrFile(arg string) string {
	var file []byte
	var err error
	if file, err = ioutil.ReadFile(arg); os.IsNotExist(err) {
		return arg
	}
	// Remove last \n or \r
	if file[len(file)-1] == byte(10) || file[len(file)-1] == byte(13) {
		file = file[:len(file)-1]
	}
	return string(file)
}

// RandomString return a string of length n
func RandomString(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	s := make([]rune, n)
	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}
	return string(s)
}

// StructToMap return a url.Values from a struct
func StructToMap(i interface{}) (values url.Values) {
	values = url.Values{}
	iVal := reflect.ValueOf(i).Elem()
	typ := iVal.Type()
	for i := 0; i < iVal.NumField(); i++ {
		values.Set(typ.Field(i).Tag.Get("form"), fmt.Sprint(iVal.Field(i)))
	}
	return
}

// NewUUID generate an UUID
func NewUUID() (string, error) {
	var uuid = make([]byte, 16)
	_, err := rand.Read(uuid)
	if err != nil {
		return "", err
	}

	uuid[6] = (uuid[6] & 0x0f) | 0x40 // Version 4
	uuid[8] = (uuid[8] & 0x3f) | 0x80 // Variant is 10
	uuidString := fmt.Sprintf("%x-%x-%x-%x-%x",
		uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:])
	return uuidString, nil

}

// GetUserAgent return an agent among popular user agent
func GetUserAgent() string {
	var userAgents = []string{"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.131 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.107 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:90.0) Gecko/20100101 Firefox/90.0",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.164 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.131 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.107 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/14.1.2 Safari/605.1.15",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/14.1.1 Safari/605.1.15",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:91.0) Gecko/20100101 Firefox/91.0"}
	return userAgents[rand.Intn(len(userAgents))]
}

// ReSubMatchMap will applied a regex with named capture and return a map
// Credits: https://stackoverflow.com/a/46202939/7245054
func ReSubMatchMap(r *regexp.Regexp, str string) map[string]string {
	matches := r.FindAllStringSubmatch(str, -1)

	subMatchMap := make(map[string]string)
	for _, match := range matches {
		for i, name := range r.SubexpNames() {
			if i != 0 && len(match) >= i {
				subMatchMap[name] = match[i]
			}
		}
	}

	return subMatchMap
}

// GetBodyInWebsite return the body of the website
func GetBodyInWebsite(url string, proxy func(*http.Request) (*url.URL, error), headers map[string]string) (string, int, error) {
	// Get random user agent
	userAgent := GetUserAgent()
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("User-Agent", userAgent)
	// Add the headers to the request
	for headerName, headerValue := range headers {
		req.Header.Add(headerName, headerValue)
	}

	client := &http.Client{

		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			Proxy:           proxy,
		},
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", -1, err
	}
	body, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	return string(body), resp.StatusCode, nil
}

// UniqueSliceString returns a slice without any duplicae
func UniqueSliceString(list []string) []string {
	check := make(map[string]bool)
	var uinque []string

	for _, val := range list {
		if _, ok := check[val]; !ok {
			check[val] = true
			uinque = append(uinque, val)
		}
	}

	return uinque
}

// SearchReplaceMap will search and replace a string by another in all the map
func SearchReplaceMap(mapToReplace map[string]string, old, new string) map[string]string {
	var newMap = make(map[string]string)
	for key, value := range mapToReplace {
		value = strings.ReplaceAll(value, old, new)
		newMap[key] = value
	}
	return newMap
}

// GetKeysMap return all the keys of the map. According to https://programmerah.com/how-to-get-all-the-keys-of-map-by-golang-1723/ this is the most efficient way to do it
func GetKeysMap(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

// IndexInSlice returns the index of the element in the given slice. If the element is not found it returns -1
func IndexInSlice(slice []string, value string) int {
	for p, v := range slice {
		if v == value {
			return p
		}
	}
	return -1
}

func StringInSlice(sliceStr []string, str string) bool {
	for _, value := range sliceStr {
		if value == str {
			return true
		}
	}
	return false
}

/*func DataToHTML(dataTable interface{}, title string) bytes.Buffer {
	var columns []string
	v := reflect.ValueOf(dataTable).Type().Elem()
	for i := 0; i < v.NumField(); i++ {
		columns = append(columns, v.Field(i).Name)
	}
	var tpl bytes.Buffer
	t, _ := template.ParseFiles("template/datatables.tpl")
	//f, _ := os.Create("users.html")

	data := struct {
		DataTable interface{}
		Columns   []string
		Title     string
	}{
		DataTable: dataTable,
		Columns:   columns,
		Title:     title,
	}
	t.Execute(&tpl, data)
	return tpl
}*/
func DataToHTML(rows [][]string, columns []string, title string) bytes.Buffer {

	var tpl bytes.Buffer
	customFunctions := template.FuncMap{
		"replace": func(input, from, to string) string { return strings.Replace(input, from, to, -1) },
	}

	t, _ := template.New("datatables.tpl").Funcs(customFunctions).ParseFiles("template/datatables.tpl")
	//t, err := template.ParseFiles("template/datatables.tpl")
	//f, _ := os.Create("users.html")

	data := struct {
		Rows    [][]string
		Columns []string
		Title   string
	}{
		Rows:    rows,
		Columns: columns,
		Title:   title,
	}
	t.Execute(&tpl, data)
	return tpl
}

func SearchInStruct(item reflect.Value, column string) string {
	var element string

	switch item.FieldByName(column).Type().Kind() {
	case reflect.Slice:
		var dataSlice string
		for j := 0; j < item.FieldByName(column).Len(); j++ {

			dataSlice += item.FieldByName(column).Index(j).String() + "\n"
		}
		element = dataSlice

	case reflect.String:
		element = item.FieldByName(column).String()
	default:
		fmt.Println(item.FieldByName(column).Type().Kind())
	}

	return element
}
