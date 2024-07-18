package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type (
	Protocol struct {
		Name string `json:"name"`
		Port int    `json:"port"`
	}

	Domain struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}

	Platform struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}

	JmsGateway struct {
		Name      string     `json:"name"`
		Address   string     `json:"address"`
		Domain    Domain     `json:"domain"`
		Platform  Platform   `json:"platform"`
		Protocols []Protocol `json:"protocols"`
		OrgID     string     `json:"org_id"`
		OrgName   string     `json:"org_name"`
		IsActive  bool       `json:"is_active"`
	}
	JmsResponse struct {
		Count    int           `json:"count,omitempty"`
		Next     string        `json:"next"`
		Previous string        `json:"previous"`
		Results  []*JmsGateway `json:"results"`
	}

	JmsClient struct {
		Address string
		Token   string
		Client  *http.Client
		Logger  log.Logger
	}
)

func NewJmsClient(address, token string) *JmsClient {
	return &JmsClient{
		Address: address,
		Token:   token,
		Client:  &http.Client{},
		Logger:  log.Logger{},
	}
}

func (j *JmsClient) GatewayList() (ret []*ConnectionStatus, err error) {
	gatewayUrl := fmt.Sprintf("%s/api/v1/assets/hosts/?platform=Gateway&offset=0&limit=10000&display=1&draw=1", j.Address)
	req, err := http.NewRequest(http.MethodGet, gatewayUrl, nil)
	if err != nil {
		return
	}
	req.Header.Add("Authorization", "Token "+j.Token)
	req.Header.Add("X-JMS-ORG", "00000000-0000-0000-0000-000000000002")
	resp, err := j.Client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	jmsResp := &JmsResponse{}
	if err = json.Unmarshal(body, jmsResp); err != nil {
		return
	}

	for _, item := range jmsResp.Results {
		cs := &ConnectionStatus{
			IP:   item.Address,
			Name: item.Name,
			IsUp: false,
		}
		for _, protocol := range item.Protocols {
			if protocol.Name == "ssh" {
				cs.Port = protocol.Port
			}
		}
		ret = append(ret, cs)
	}

	return
}
