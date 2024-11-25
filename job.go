package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
)

type sendDataResponse struct {
	Message string `json:"message"`
}

type Client struct {
	baseUrl        string
	authEndpoint   string
	writerEndpoint string
	clientId       string
	clientSecret   string
	headers        http.Header
	httpClient     *http.Client
}

type User struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	ClientId   string `json:"clientId"`
	CustomerId int    `json:"customerId"`
	Role       string `json:"role"`
}

type getAuthLoginResponse struct {
	Message string `json:"message"`
	User    User   `json:"user"`
	Token   string `json:"token"`
}

func new_client(baseUrl, authEndpoint, writerEndpoint, clientId, clientSecret string) *Client {
	headers := http.Header{}
	headers.Set("Content-Type", "application/json")

	return &Client{
		baseUrl:        baseUrl,
		headers:        headers,
		authEndpoint:   authEndpoint,
		writerEndpoint: writerEndpoint,
		clientId:       clientId,
		clientSecret:   clientSecret,
		httpClient:     &http.Client{},
	}
}

func (c *Client) auth_login() (getAuthLoginResponse, error) {
	empty := getAuthLoginResponse{}

	bytesObj := []byte(`{"clientId":"` + c.clientId + `", "clientSecret":"` + c.clientSecret + `"}`)
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
		return empty, errors.New("Falha no login. BODY: " + string(body))
	}

	response := getAuthLoginResponse{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return empty, err
	}

	return response, nil
}

func (c *Client) send_data(token string, records []map[string]interface{}, offset, customerID int, clientID string) error {
	payload, err := json.Marshal(map[string]interface{}{
		"records": records,
	})
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", c.baseUrl+c.writerEndpoint, bytes.NewBuffer(payload))
	if err != nil {
		return err
	}
	if err != nil {
		return err
	}

	req.Header = c.headers
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("CustomerID", strconv.Itoa(customerID))
	req.Header.Set("ClientID", clientID)

	if offset == 0 {
		req.Header.Set("DeleteData", "1")
	} else {
		req.Header.Set("DeleteData", "0")
	}

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
		return errors.New("Falha enviando os registros. STATUS: " + resp.Status)
	}

	response := sendDataResponse{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return err
	}

	if response.Message != "success" {
		return errors.New("Falha enviando os registros. BODY: " + string(body))
	}

	return nil
}

func RUN(db *database, baseUrl, authUrl, writerUrl, clientId, clientSecret string, batchSize int, initialDate string) func() {
	client := new_client(baseUrl, authUrl, writerUrl, clientId, clientSecret)

	return func() {
		logger.Info("Running job...")

		auth, err := client.auth_login()
		if err != nil {
			logger.Warn("auth_login: " + err.Error())
			return
		}

		customerID := auth.User.CustomerId
		clientID := auth.User.ClientId
		token := auth.Token

		if token == "" {
			logger.Warn("Failed to get token")
			return
		}

		offset := 0
		for {
			records, err := db.get_data(batchSize, offset, customerID, initialDate)
			if err != nil {
				logger.Info("Erro ao buscar registros: " + err.Error())
				break
			}

			if len(records) == 0 {
				logger.Info("Processo concluído.")
				break
			}

			err = client.send_data(token, records, offset, customerID, clientID)
			if err != nil {
				logger.Warn("Erro ao enviar registros para a API: " + err.Error())
				break
			}

			logger.Info("Enviando registros para a API", "Qtd:", len(records))

			offset += batchSize
		}
	}
}
