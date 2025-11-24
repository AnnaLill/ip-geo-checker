package geo

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// IPInfoResponse структура для ответа от ipinfo.io
type IPInfoResponse struct {
	City string `json:"city"`
}

// IPAPIResponse структура для ответа от ip-api.com (XML формат)
type IPAPIResponse struct {
	XMLName xml.Name `xml:"query"`
	City    string   `xml:"city"`
}

// IPWhoIsResponse структура для ответа от ipwho.is
type IPWhoIsResponse struct {
	City string `json:"city"`
}

// DBIPResponse структура для ответа от api.db-ip.com
type DBIPResponse struct {
	City string `json:"city"`
}

// IPMeResponse структура для ответа от ip.me
type IPMeResponse struct {
	City string `json:"city"`
}

// fetchFromIPInfo запрашивает данные от ipinfo.io
func fetchFromIPInfo(ip string) Result {
	url := fmt.Sprintf("http://ipinfo.io/%s/json", ip)

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		return Result{APIName: "ipinfo.io", Error: fmt.Errorf("ipinfo.io request failed: %w", err)}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Result{APIName: "ipinfo.io", Error: fmt.Errorf("ipinfo.io returned status %d", resp.StatusCode)}
	}

	var result IPInfoResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return Result{APIName: "ipinfo.io", Error: fmt.Errorf("ipinfo.io json decode failed: %w", err)}
	}

	if result.City == "" {
		return Result{APIName: "ipinfo.io", Error: fmt.Errorf("ipinfo.io: city field is empty")}
	}

	return Result{APIName: "ipinfo.io", City: result.City}
}

// fetchFromIPAPI запрашивает данные от ip-api.com
func fetchFromIPAPI(ip string) Result {
	url := fmt.Sprintf("http://ip-api.com/xml/%s", ip)

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		return Result{APIName: "ip-api.com", Error: fmt.Errorf("ip-api.com request failed: %w", err)}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Result{APIName: "ip-api.com", Error: fmt.Errorf("ip-api.com returned status %d", resp.StatusCode)}
	}

	var result IPAPIResponse
	if err := xml.NewDecoder(resp.Body).Decode(&result); err != nil {
		return Result{APIName: "ip-api.com", Error: fmt.Errorf("ip-api.com xml decode failed: %w", err)}
	}

	if result.City == "" {
		return Result{APIName: "ip-api.com", Error: fmt.Errorf("ip-api.com: city field is empty")}
	}

	return Result{APIName: "ip-api.com", City: result.City}
}

// fetchFromIPWhoIs запрашивает данные от ipwho.is
func fetchFromIPWhoIs(ip string) Result {
	url := fmt.Sprintf("http://ipwho.is/%s", ip)

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		return Result{APIName: "ipwho.is", Error: fmt.Errorf("ipwho.is request failed: %w", err)}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Result{APIName: "ipwho.is", Error: fmt.Errorf("ipwho.is returned status %d", resp.StatusCode)}
	}

	var result IPWhoIsResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return Result{APIName: "ipwho.is", Error: fmt.Errorf("ipwho.is json decode failed: %w", err)}
	}

	if result.City == "" {
		return Result{APIName: "ipwho.is", Error: fmt.Errorf("ipwho.is: city field is empty")}
	}

	return Result{APIName: "ipwho.is", City: result.City}
}

// fetchFromDBIP запрашивает данные от api.db-ip.com
func fetchFromDBIP(ip string) Result {
	url := fmt.Sprintf("https://api.db-ip.com/v2/free/%s", ip)

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		return Result{APIName: "db-ip.com", Error: fmt.Errorf("db-ip.com request failed: %w", err)}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Result{APIName: "db-ip.com", Error: fmt.Errorf("db-ip.com returned status %d", resp.StatusCode)}
	}

	var result DBIPResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return Result{APIName: "db-ip.com", Error: fmt.Errorf("db-ip.com json decode failed: %w", err)}
	}

	if result.City == "" {
		return Result{APIName: "db-ip.com", Error: fmt.Errorf("db-ip.com: city field is empty")}
	}

	return Result{APIName: "db-ip.com", City: result.City}
}

// fetchFromIPMe запрашивает данные от ip.me
func fetchFromIPMe(ip string) Result {
	url := fmt.Sprintf("https://ip.me/ip/%s", ip)

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		return Result{APIName: "ip.me", Error: fmt.Errorf("ip.me request failed: %w", err)}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Result{APIName: "ip.me", Error: fmt.Errorf("ip.me returned status %d", resp.StatusCode)}
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return Result{APIName: "ip.me", Error: fmt.Errorf("ip.me read body failed: %w", err)}
	}

	bodyStr := strings.TrimSpace(string(bodyBytes))

	var jsonResult IPMeResponse
	if err := json.Unmarshal(bodyBytes, &jsonResult); err == nil && jsonResult.City != "" {
		return Result{APIName: "ip.me", City: jsonResult.City}
	}

	return Result{APIName: "ip.me", Error: fmt.Errorf("ip.me: unable to parse response, got: %s", bodyStr)}
}
