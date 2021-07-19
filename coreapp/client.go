package coreapp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"time"
)

type CoreApp struct {
	Host           string    `json:"host"`
	User           string    `json:"user"`
	Password       string    `json:"password"`
	ClientID       string    `json:"client_id"`
	ClientSecret   string    `json:"client_secret"`
	AccessToken    string    `json:"access_token"`
	RefreshToken   string    `json:"refresh_token"`
	ExpiresIn      float64   `json:"expires_in"`
	ExpirationDate time.Time `json:"expiration_date"`
	Company        string    `json:"company"`
}

//Check if the access token is valid and if not, refresh it
func (c *CoreApp) IsValid() bool {
	if c.AccessToken == "" {
		return false
	}
	if time.Now().After(c.ExpirationDate) {
		c.RefreshTokenMethod()
		return false
	}
	return true
}

//Refresh Token
func (c *CoreApp) RefreshTokenMethod() {
	if c.RefreshToken == "" {
		log.Println("No refresh token found")
		return
	}
	log.Println("Refreshing token")
	resp, err := http.Post(c.Host+"/webservices/auth/token", "application/json", bytes.NewBuffer([]byte(`{"refresh_token": "`+c.RefreshToken+`", "client_id": "`+c.ClientID+`", "client_secret": "`+c.ClientSecret+`", "grant_type": "refresh_token"}`)))
	if err != nil {
		log.Println("Error refreshing token")
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error reading body")
		return
	}
	var tokenResponse map[string]interface{}
	err = json.Unmarshal(body, &tokenResponse)
	if err != nil {
		log.Println("Error unmarshalling body")
		return
	}
	c.AccessToken = tokenResponse["access_token"].(string)
	c.RefreshToken = tokenResponse["refresh_token"].(string)
	c.ExpiresIn = tokenResponse["expires_in"].(float64)
	c.ExpirationDate = time.Now().Add(time.Duration(c.ExpiresIn) * time.Second)
	log.Println("Token refreshed")
}

//Login into Coreapp, making a POST request with client_id, client_secret, password, username, and grant_type
func (c *CoreApp) Login() {

	//payload
	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)

	writer.WriteField("client_id", c.ClientID)
	writer.WriteField("client_secret", c.ClientSecret)
	writer.WriteField("grant_type", "password")
	writer.WriteField("password", c.Password)
	writer.WriteField("username", c.User)

	err := writer.Close()

	if err != nil {
		log.Println("Error writing payload")
		return
	}
	//request
	req, err := http.NewRequest("POST", c.Host+"/webservices/auth/token/", payload)

	//add headers
	req.Header.Add("Content-Type", writer.FormDataContentType())
	req.Header.Add("Accept", "application/json")

	if err != nil {
		log.Fatal(err)
	}
	//client
	client := &http.Client{}
	//response
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	//response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	//response body string
	response := string(body)

	fmt.Println(response)
	//json response
	var token map[string]interface{}
	json.Unmarshal([]byte(response), &token)
	//access_token
	c.AccessToken = token["access_token"].(string)
	//refresh_token
	c.RefreshToken = token["refresh_token"].(string)
	//expires_in
	c.ExpiresIn = token["expires_in"].(float64)
	//expiration_date
	c.ExpirationDate = time.Now().Add(time.Duration(c.ExpiresIn) * time.Second)

}

//SendSMS usign Coreapp endpoint "/webservices/custom_function/sendSMS/"
func (c *CoreApp) SendSMS(to []string, message string) {

	//Login if access token is empty
	if c.AccessToken == "" {
		c.Login()
	}

	//Check if access token is valid and if not, refresh it
	if !c.IsValid() {
		c.RefreshTokenMethod()
	}

	//data map
	data := make(map[string]interface{})
	data["to"] = to
	data["body"] = message
	data["record_source"] = "Monitor"

	//jsonValue
	jsonValue, _ := json.Marshal(data)
	//request
	req, err := http.NewRequest("POST", c.Host+"/webservices/custom_function/sendSMS/", bytes.NewBuffer(jsonValue))
	if err != nil {
		log.Fatal(err)
	}
	//add headers
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", "Bearer "+c.AccessToken)
	req.Header.Add("company", c.Company)
	//client
	client := &http.Client{}
	//response
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	//response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	//response body string
	response := string(body)
	fmt.Println(response)
}

//SendEmail usign Coreapp endpoint "/webservices/custom_function/sendMails/"
func (c *CoreApp) SendEmail(to []string, message string, subject string) {
	//Login if access token is empty
	if c.AccessToken == "" {
		c.Login()
	}
	//Check if access token is valid and if not, refresh it
	if !c.IsValid() {
		c.RefreshTokenMethod()
	}
	//data map
	data := make(map[string]interface{})
	data["to"] = to
	data["body"] = message
	data["record_source"] = "Monitor"
	data["subject"] = subject

	//jsonValue
	jsonValue, _ := json.Marshal(data)
	//request
	req, err := http.NewRequest("POST", c.Host+"/webservices/custom_function/sendMails/", bytes.NewBuffer(jsonValue))
	if err != nil {
		log.Fatal(err)
	}
	//Send access_token in header
	req.Header.Add("Authorization", "Bearer "+c.AccessToken)
	req.Header.Add("company", c.Company)
	//client
	client := &http.Client{}
	//response
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	//response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	//response body string
	response := string(body)
	fmt.Println(response)
}
