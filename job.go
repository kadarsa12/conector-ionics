package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

type Client struct {
	baseUrl        string
	authEndpoint   string
	writerEndpoint string
	username       string
	password       string
	headers        http.Header
	httpClient     *http.Client
}

type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Role     string `json:"role"`
}

type getAuthLoginResponse struct {
	Message string `json:"message"`
	User    User   `json:"user"`
	Token   string `json:"token"`
}

// type sendInvoicesResponse struct {
// 	Message string `json:"message"`
// 	Status  string `json:"status"`
// }

func new_client(baseUrl, authEndpoint, writerEndpoint, username, password string) *Client {
	headers := http.Header{}
	headers.Set("Content-Type", "application/json")

	return &Client{
		baseUrl:        baseUrl,
		headers:        headers,
		authEndpoint:   authEndpoint,
		writerEndpoint: writerEndpoint,
		username:       username,
		password:       password,
		httpClient:     &http.Client{},
	}
}

func (c *Client) auth_login() (getAuthLoginResponse, error) {
	empty := getAuthLoginResponse{}

	bytesObj := []byte(`{"username":"` + c.username + `", "password":"` + c.password + `"}`)
	bodyObj := bytes.NewBuffer(bytesObj)

	req, err := http.NewRequest("POST", c.baseUrl+c.authEndpoint, bodyObj)
	if err != nil {
		return empty, err
	}
	req.Header = c.headers

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return empty, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return empty, err
	}

	if resp.StatusCode != 200 {
		return empty, errors.New("Failed to auth login. BODY: " + string(body))
	}

	response := getAuthLoginResponse{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return empty, err
	}

	return response, nil
}

func (c *Client) send_data(token string, qtySendData int) error {
	invoicesBody, err := json.Marshal(map[string]any{"invoices": invoices, "inserted": inserted, "not_found": notFound, "removed": disapproved})
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", c.baseUrl+c.writerEndpoint, bytes.NewBuffer(invoicesBody))
	if err != nil {
		return err
	}
	req.Header = c.headers

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return errors.New("Failed to send invoices. BODY: " + string(body) + " STATUS: " + resp.Status)
	}

	response := sendInvoicesResponse{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return err
	}

	if response.Status != "success" {
		return errors.New("Failed to send invoices. BODY: " + string(body))
	}

	return nil
}

func RUN(db *database, baseUrl, authUrl, writerUrl, username, password string, qtySendData int) func() {
	client := new_client(baseUrl, authUrl, writerUrl, username, password)

	return func() {
		logger.Info("Running job...")

		auth, err := client.auth_login()
		if err != nil {
			logger.Warn("auth_login: " + err.Error())
			return
		}

		token := auth.Token
		if token == "" {
			logger.Warn("Failed to get token")
			return
		}

		data = db.get_data(qtySendData)

		client.headers.Set("Authorization", "Bearer "+token)
		data, err := client.send_data(qtySendData)
		if err != nil {
			logger.Warn("send_data: " + err.Error())
			return
		}

		// err = db.insert_status(data.Approved, "L", data.CompanyID)
		// if err != nil {
		// 	logger.Warn("insert_status error: " + err.Error())
		// }

		// err = db.remove_status(data.Disapproved, data.CompanyID)
		// if err != nil {
		// 	logger.Warn("remove_status error: " + err.Error())
		// }

		// invoices, err := db.get_invoices(data.New)
		// if err != nil {
		// 	logger.Warn("get_invoices error: " + err.Error())
		// 	return
		// }

		// notFound := []string{}
		// for _, invoice := range data.New {
		// 	found := false
		// 	for _, i := range invoices {
		// 		if i.NfeKey != nil && *i.NfeKey == invoice {
		// 			found = true
		// 			break
		// 		}
		// 	}
		// 	if !found {
		// 		notFound = append(notFound, invoice)
		// 	}
		// }

		// err = client.send_invoices(invoices, data.Approved, notFound, data.Disapproved)
		// if err != nil {
		// 	logger.Warn(err.Error())
		// 	return
		// }

		logger.Info("Sent new invoices data to FG...")
	}
}
