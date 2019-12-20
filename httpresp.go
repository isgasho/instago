package instago

// This file sends HTTP requests and gets HTTP responses.

import (
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
)

var userAgent = "Instagram 10.26.0 (iPhone8,1; iOS 10_2; en_US; en-US; " +
	"scale=2.00; gamut=normal; 750x1334) AppleWebKit/420+"

// SetUserAgent let you set User-Agent header in HTTP requests.
func SetUserAgent(s string) {
	userAgent = s
}

// Send HTTP request and get http response without login and with gis info. Used
// in get all post codes without login.
func getHTTPResponseNoLoginWithGis(url, gis string) (b []byte, err error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}

	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("X-Instagram-GIS", gis)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		err = errors.New(url +
			"\nresp.StatusCode: " + strconv.Itoa(resp.StatusCode))
		return
	}

	return ioutil.ReadAll(resp.Body)
}

// Send HTTP request and get http response without login.
func getHTTPResponseNoLogin(url string) (b []byte, err error) {
	resp, err := http.Get(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		err = errors.New(url +
			"\nresp.StatusCode: " + strconv.Itoa(resp.StatusCode))
		return
	}

	return ioutil.ReadAll(resp.Body)
}

// Send HTTP request and get http response on behalf of a specific Instagram
// user. After login to Instagram, you can get the cookies of *ds_user_id*,
// *sessionid*, *csrftoken* in Chrome Developer Tools.
// See https://stackoverflow.com/a/44773079
// or
// https://github.com/hoschiCZ/instastories-backup#obtain-cookies
func (m *IGApiManager) getHTTPResponse(url, method string) (b []byte, err error) {
	if method != "POST" {
		method = "GET"
	}

	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return
	}

	for name, value := range m.cookies {
		req.AddCookie(&http.Cookie{Name: name, Value: value})
	}

	req.Header.Set("User-Agent", userAgent)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		err = errors.New(url +
			"\nresp.StatusCode: " + strconv.Itoa(resp.StatusCode))
		return
	}

	return ioutil.ReadAll(resp.Body)
}
