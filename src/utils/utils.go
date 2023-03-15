package utils

import (
	templateResources "GoMapEnum/src/template"
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/md5"
	"crypto/rc4"
	"crypto/tls"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"regexp"
	"strings"
	"text/template"
	"time"

	"github.com/rs/dnscache"
	"golang.org/x/net/proxy"
)

func init() {
}

var Transporter *http.Transport

func init() {

	rand.Seed(time.Now().UnixNano())

	/*
		'High performance' http transport for golang
		increases MaxIdleConns and conns per host since we expect
		to be talking to a lot of other hosts all the time
		Also adds a basic in-process dns cache to help
		in docker environments since the standard alpine build appears
		to have no in container dns cache
	*/
	r := &dnscache.Resolver{}
	Transporter = &http.Transport{

		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		DialContext: func(ctx context.Context, network string, addr string) (conn net.Conn, err error) {
			separator := strings.LastIndex(addr, ":")
			ips, err := r.LookupHost(ctx, addr[:separator])
			if err != nil {
				return nil, err
			}
			for _, ip := range ips {
				conn, err = net.Dial(network, ip+addr[separator:])
				if err == nil {
					break
				}
			}
			return
		},
		MaxIdleConns:    1024,
		MaxConnsPerHost: 100,
		IdleConnTimeout: 10 * time.Second,
	}
	go func() {
		clearUnused := true
		t := time.NewTicker(5 * time.Minute)
		defer t.Stop()
		for range t.C {
			r.Refresh(clearUnused)
		}
	}()
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
func GetBodyInWebsite(url string, proxy func(*http.Request) (*url.URL, error), headers map[string]string, bodyRequest io.Reader) (string, int, error) {
	// Get random user agent
	userAgent := GetUserAgent()
	var req *http.Request
	if bodyRequest == nil {
		req, _ = http.NewRequest("GET", url, nil)
	} else {
		req, _ = http.NewRequest("POST", url, bodyRequest)
	}
	req.Header.Add("User-Agent", userAgent)
	// Add the headers to the request
	for headerName, headerValue := range headers {
		req.Header.Add(headerName, headerValue)
	}
	Transporter.Proxy = proxy
	client := &http.Client{
		Transport: Transporter,
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", -1, err
	}
	body, _ := io.ReadAll(resp.Body)
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
func SearchReplaceMap(mapToReplace map[string]string, oldString, newString string) map[string]string {
	var newMap = make(map[string]string)
	for key, value := range mapToReplace {
		value = strings.ReplaceAll(value, oldString, newString)
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

// StringInSlice returns true if the string is present in the slice, otherwise false
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

// DataToHTML returns a datatables html page. It takes as input the rows, the columns and the title.
func DataToHTML(rows [][]string, columns []string, title string) bytes.Buffer {

	var tpl bytes.Buffer
	customFunctions := template.FuncMap{
		"replace": func(input, from, to string) string { return strings.ReplaceAll(input, from, to) },
	}

	t, _ := template.New("datatables.tpl").Funcs(customFunctions).Parse(templateResources.GetTemplateDatatables())

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

// SearchInStruct iterate over the a struct to search for a field name (column parameter) and return the value separeted by a new line
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

// OpenConnectionWoProxy open a connection to the given url with a proxy or not if not given.
func OpenConnectionWoProxy(target, port string, timeout int, proxyTCP proxy.Dialer) (net.Conn, error) {
	var conn net.Conn
	var err error
	if proxyTCP != nil {
		// Use the proxy dialer
		conn, err = proxyTCP.Dial("tcp", net.JoinHostPort(target, port))
	} else {
		// If not proxy is given
		defaultDialer := &net.Dialer{Timeout: time.Duration(timeout * int(time.Second))}
		conn, err = defaultDialer.Dial("tcp", net.JoinHostPort(target, port))
	}

	// Check the error
	if err != nil || conn == nil {
		var errStr string
		if err != nil {
			errStr = err.Error()
		} else {
			errStr = "No error but connection is nil"
		}
		return conn, fmt.Errorf(errStr)
	}
	return conn, nil
}

// GetHmacMd5 returns the hmac of the given string with the given key
func GetHmacMd5(data, key []byte) []byte {
	mac := hmac.New(md5.New, key)
	mac.Write(data)
	return mac.Sum(nil)
}

// RC4Decrypt decrypts the given data with the given key with RC4 algorithm
func RC4Decrypt(encData, key []byte) ([]byte, error) {
	dst := make([]byte, len(encData))
	rc4cipher, err := rc4.NewCipher(key)

	if err != nil {
		return nil, err
	}

	rc4cipher.XORKeyStream(dst, encData)
	return dst, nil
}
