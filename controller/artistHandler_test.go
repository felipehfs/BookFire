// Package controller tests the controller
package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/bookfire/model"
)

const (
	PORT = 8081
)

var artistsRoute = fmt.Sprintf("http://localhost:%d/artists/", PORT)

// getToken provide the access to database
func getToken() ([]byte, error) {
	url := fmt.Sprintf("http://localhost:%v/login/", PORT)
	form := []byte(`{ "login": "admin", "password": "admin" }`)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(form))
	if err != nil {
		return nil, err
	}
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	fmt.Println(string(body))
	return body, nil
}

// TestCreateArtist tests if the insertion works well
func TestCreateArtist(t *testing.T) {
	var request = struct {
		Name        string
		Email       string
		Description string
	}{
		Name:  "Belchior",
		Email: "belchior@ig.com.br",
		Description: `Apenas um latino americano, 
			sem conta no banco, sem parentes importantes e vindo do interior`,
	}

	jsonData, err := json.Marshal(request)
	if err != nil {
		t.Error(err)
	}

	client := &http.Client{}

	req, err := http.NewRequest("POST", artistsRoute, bytes.NewBuffer(jsonData))
	token, err := getToken()
	if err != nil {
		t.Error(err)
	}
	req.Header.Add("Authorization", string(token))
	if err != nil {
		t.Error(err)
	}

	resp, err := client.Do(req)
	if err != nil {
		t.Error(err)
	}

	if resp.StatusCode != http.StatusCreated {
		t.Errorf("Status code expected: %v, got %v", http.StatusCreated, resp.StatusCode)
	}
	artists, err := getArtists()

	if err != nil {
		t.Error(err)
	}

	last := len(artists) - 1
	equals := artists[last].Name == request.Name && artists[last].Email == request.Email
	if !equals {
		t.Errorf("Values are different: expected %s when %s", request.Name, artists[0].Name)
	}
}

// checkRequestOK makes the test to fail if status code is not ok
func checkRequestOK(response *http.Response, t *testing.T) {
	if response.StatusCode != http.StatusOK {
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			t.Fatal(err)
		}
		t.Errorf("Status code different received: %s %s", response.Status, string(body))
	}
}

// TestReadArtists tests the data retrieve
func TestReadArtists(t *testing.T) {
	request, err := http.NewRequest("GET", artistsRoute, nil)
	token, err := getToken()
	if err != nil {
		t.Error(err)
	}
	request.Header.Add("Authorization", string(token))
	if err != nil {
		t.Error(err)
	}
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		t.Error(err)
	}

	checkRequestOK(response, t)
}

// getArtists retrieves all data from the server
// if exists
func getArtists() ([]model.Artist, error) {
	req, err := http.NewRequest("GET", artistsRoute, nil)
	token, err := getToken()
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", string(token))
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	var artists []model.Artist
	body, err := ioutil.ReadAll(resp.Body)

	if err := json.Unmarshal(body, &artists); err != nil {
		return nil, err
	}

	return artists, nil
}

// TestFindArtist searches the data through the get method
func TestFindArtists(t *testing.T) {
	searched, err := getArtists()
	if err != nil {
		t.Error(err)
	}

	if len(searched) < 1 {
		t.Error("Empty list")
	}

	tt := []struct {
		id       string
		name     string
		expected int
	}{
		{searched[0].ID.Hex(), "busca que funciona", 200},
		{"5b608966767bb623cf13c303", "busca que falha", 404},
	}

	token, err := getToken()
	if err != nil {
		t.Error(err)
	}

	auth := string(token)

	for _, tr := range tt {
		req, err := http.NewRequest("GET", artistsRoute+tr.id, nil)
		if err != nil {
			t.Error(err)
		}
		req.Header.Add("Authorization", auth)
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			t.Error(err)
		}

		if resp.StatusCode != tr.expected {
			t.Errorf("Status code expected %d received %d", tr.expected, resp.StatusCode)
		}
	}
}

// TestUpdateArtist tests the change of a artists
func TestUpdateArtist(t *testing.T) {
	artists, err := getArtists()
	if err != nil {
		t.Error(err)
	}

	last := len(artists) - 1
	if last < 0 {
		t.Errorf("empty list")
	}

	updated := artists[last]
	updated.Name = "Alter"
	updated.Email = "alterado@gmail.com"

	newURL := artistsRoute + updated.ID.Hex()

	body, err := json.Marshal(updated)
	if err != nil {
		t.Error(err)
	}
	req, err := http.NewRequest("PUT", newURL, bytes.NewBuffer(body))
	token, err := getToken()

	if err != nil {
		t.Error(err)
	}
	req.Header.Add("Authorization", string(token))
	if err != nil {
		t.Error(err)
	}
	client := &http.Client{}
	res, err := client.Do(req)

	if err != nil {
		t.Error(err)
	}

	if res.StatusCode != http.StatusOK {
		t.Errorf("Status expected %v got %v", http.StatusOK, res.StatusCode)
	}
}

// TestDeleteArtist tests the delete method
func TestDeleteArtist(t *testing.T) {
	artists, err := getArtists()
	if err != nil {
		t.Error(err)
	}
	if len(artists) < 0 {
		t.Error("empty list")
	}
	last := artists[len(artists)-1]
	tt := []struct {
		id           string
		expectedCode int
	}{
		{last.ID.Hex(), 200},
		{"5b608966767bb623cf13c303", 500},
	}

	token, err := getToken()
	if err != nil {
		t.Error(err)
	}
	for _, tr := range tt {
		req, err := http.NewRequest("DELETE", artistsRoute+tr.id, nil)
		req.Header.Add("Authorization", string(token))
		if err != nil {
			t.Error(err)
		}
		client := &http.Client{}
		res, err := client.Do(req)
		if err != nil {
			t.Error(err)
		}

		if res.StatusCode != tr.expectedCode {
			t.Errorf("Error: expected %v got %v\n", tr.expectedCode, res.StatusCode)
		}
	}
}

// TestUpdateManyCases test some possibilities of crash in update
func TestUpdateManyCase(t *testing.T) {
	artists, err := getArtists()
	if err != nil {
		t.Error(err)
	}
	if len(artists) < 0 {
		t.Error("Empty list")
	}
	last := artists[len(artists)-1]

	tt := []struct {
		name           string
		id             string
		body           string
		statusExpected int
	}{
		{"Update which exists",
			last.ID.Hex(),
			`{"name": "data updated", "email": "dataupdated@gmail.com"}`,
			200},
		{"Update which id not exists",
			"5b6057e8767bb423cf13c305",
			`{"name": "Belchior", "email": "belchior@gmail.com"}`,
			500},
		{
			"Update which id returns error",
			"5b6057e8767bb623cf13c305",
			`{ "_id": 5b6057e8767bb623cf13c303, 
				"name": "Belchior error", "email": "belchior@gmail.com",
				"description": "fjfjs"}`,
			500},
		{
			"Update which values the same",
			"5b47efb6391bc8b26eb2deac",
			`{"name": "Ken thompson", "email": "kenthompson@gmail.com", "description": "O cara que criou o grande C"}`,
			200},
	}

	token, err := getToken()
	if err != nil {
		t.Error(err)
	}

	for _, tr := range tt {
		req, err := http.NewRequest("PUT", artistsRoute+tr.id, bytes.NewBuffer([]byte(tr.body)))
		req.Header.Add("Authorization", string(token))
		if err != nil {
			t.Error(err)
		}

		client := &http.Client{}
		resp, err := client.Do(req)

		if err != nil {
			t.Error(err)
		}

		if resp.StatusCode != tr.statusExpected {
			t.Errorf("%s: Error: expected %d got %d\n", tr.name, tr.statusExpected, resp.StatusCode)
		}
	}
}
