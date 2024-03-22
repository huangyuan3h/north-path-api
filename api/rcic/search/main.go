package main

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"

	"api.north-path.site/utils/errors"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type RCICRequest struct {
	RCIC string `json:"rcic" validate:"required"`
}

const target_url = "https://secure-archive.college-ic.ca/search/do/lang/en"

func Handler(request events.APIGatewayV2HTTPRequest) (events.APIGatewayProxyResponse, error) {

	var rcicReq RCICRequest
	err := json.Unmarshal([]byte(request.Body), &rcicReq)

	if err != nil {
		return errors.New(errors.JSONParseError, http.StatusBadRequest).GatewayResponse()
	}

	data := url.Values{
		"lang":                     {"en"},
		"search_iccrc_number":      {"R534829"},
		"start":                    {"0"},
		"search_last_name":         {""},
		"search_first_name":        {""},
		"search_membership_status": {""},
		"search_company_name":      {""},
		"search_agent":             {""},
		"search_country":           {""},
		"search_province":          {""},
		"search_postal_code":       {""},
		"search_city":              {""},
		"search_street":            {""},
		"query":                    {""},
		"letter":                   {""},
	}

	encoded := "form_fields=" + url.QueryEscape(data.Encode())

	req, err := http.NewRequest("POST", target_url, strings.NewReader(encoded))

	if err != nil {
		// 创建 HTTP 请求失败
		return errors.New(errors.JSONParseError, http.StatusBadRequest).GatewayResponse()
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "en-GB,en;q=0.9,en-US;q=0.6,ja;q=0.5")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Origin", "https://secure-archive.college-ic.ca")
	req.Header.Set("Referer", "https://secure-archive.college-ic.ca/search-new/EN")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0.0.0 Safari/537.36")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("sec-ch-ua", "\"Chromium\";v=\"122\", \"Not(A:Brand\";v=\"24\", \"Google Chrome\";v=\"122\"")
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", "\"macOS\"")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		// 发送请求失败
		return errors.New(errors.JSONParseError, http.StatusBadRequest).GatewayResponse()
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {

		// 读取响应结果失败
		return errors.New(errors.JSONParseError, http.StatusBadRequest).GatewayResponse()
	}

	return events.APIGatewayProxyResponse{
		Body:       string(body),
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(Handler)
}
