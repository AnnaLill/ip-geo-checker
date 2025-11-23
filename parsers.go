package main

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
// Пример ответа: {"ip": "1.1.1.1", "city": "Hong Kong", ...}
type IPInfoResponse struct {
	City string `json:"city"`
}

// IPAPIResponse структура для ответа от ip-api.com (XML формат)
// Пример ответа: <query><city>Moscow</city></query>
type IPAPIResponse struct {
	XMLName xml.Name `xml:"query"`
	City    string   `xml:"city"`
}

// IPWhoIsResponse структура для ответа от ipwho.is
// Пример ответа: {"city": "Moscow", ...}
type IPWhoIsResponse struct {
	City string `json:"city"`
}

// DBIPResponse структура для ответа от api.db-ip.com
// Пример ответа: {"city": "Moscow", ...}
type DBIPResponse struct {
	City string `json:"city"`
}

// IPMeResponse структура для ответа от ip.me
// Этот API может возвращать разные форматы, попробуем JSON
type IPMeResponse struct {
	City string `json:"city"`
}

// GeoAPIResult результат запроса к одному API
type GeoAPIResult struct {
	City  string
	Error error
}

// fetchFromIPInfo запрашивает данные от ipinfo.io
// URL формат: http://ipinfo.io/{ip}/json
func fetchFromIPInfo(ip string) GeoAPIResult {
	url := fmt.Sprintf("http://ipinfo.io/%s/json", ip)

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		return GeoAPIResult{Error: fmt.Errorf("ipinfo.io request failed: %w", err)}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return GeoAPIResult{Error: fmt.Errorf("ipinfo.io returned status %d", resp.StatusCode)}
	}

	var result IPInfoResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return GeoAPIResult{Error: fmt.Errorf("ipinfo.io json decode failed: %w", err)}
	}

	// Проверяем, что город не пустой
	if result.City == "" {
		return GeoAPIResult{Error: fmt.Errorf("ipinfo.io: city field is empty")}
	}

	return GeoAPIResult{City: result.City}
}

// fetchFromIPAPI запрашивает данные от ip-api.com
// URL формат: http://ip-api.com/xml/{ip}
// Важно: этот API возвращает XML, а не JSON!
func fetchFromIPAPI(ip string) GeoAPIResult {
	url := fmt.Sprintf("http://ip-api.com/xml/%s", ip)

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		return GeoAPIResult{Error: fmt.Errorf("ip-api.com request failed: %w", err)}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return GeoAPIResult{Error: fmt.Errorf("ip-api.com returned status %d", resp.StatusCode)}
	}

	var result IPAPIResponse
	if err := xml.NewDecoder(resp.Body).Decode(&result); err != nil {
		return GeoAPIResult{Error: fmt.Errorf("ip-api.com xml decode failed: %w", err)}
	}

	// Проверяем, что город не пустой
	if result.City == "" {
		return GeoAPIResult{Error: fmt.Errorf("ip-api.com: city field is empty")}
	}

	return GeoAPIResult{City: result.City}
}

// fetchFromIPWhoIs запрашивает данные от ipwho.is
// URL формат: http://ipwho.is/{ip}
func fetchFromIPWhoIs(ip string) GeoAPIResult {
	url := fmt.Sprintf("http://ipwho.is/%s", ip)

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		return GeoAPIResult{Error: fmt.Errorf("ipwho.is request failed: %w", err)}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return GeoAPIResult{Error: fmt.Errorf("ipwho.is returned status %d", resp.StatusCode)}
	}

	var result IPWhoIsResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return GeoAPIResult{Error: fmt.Errorf("ipwho.is json decode failed: %w", err)}
	}

	// Проверяем, что город не пустой
	if result.City == "" {
		return GeoAPIResult{Error: fmt.Errorf("ipwho.is: city field is empty")}
	}

	return GeoAPIResult{City: result.City}
}

// fetchFromDBIP запрашивает данные от api.db-ip.com
// URL формат: https://api.db-ip.com/v2/free/{ip}
func fetchFromDBIP(ip string) GeoAPIResult {
	url := fmt.Sprintf("https://api.db-ip.com/v2/free/%s", ip)

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		return GeoAPIResult{Error: fmt.Errorf("db-ip.com request failed: %w", err)}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return GeoAPIResult{Error: fmt.Errorf("db-ip.com returned status %d", resp.StatusCode)}
	}

	var result DBIPResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return GeoAPIResult{Error: fmt.Errorf("db-ip.com json decode failed: %w", err)}
	}

	// Проверяем, что город не пустой
	if result.City == "" {
		return GeoAPIResult{Error: fmt.Errorf("db-ip.com: city field is empty")}
	}

	return GeoAPIResult{City: result.City}
}

// fetchFromIPMe запрашивает данные от ip.me
// URL формат: https://ip.me/ip/{ip}
// Этот API может возвращать простой текст или JSON, попробуем оба варианта
func fetchFromIPMe(ip string) GeoAPIResult {
	url := fmt.Sprintf("https://ip.me/ip/%s", ip)

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		return GeoAPIResult{Error: fmt.Errorf("ip.me request failed: %w", err)}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return GeoAPIResult{Error: fmt.Errorf("ip.me returned status %d", resp.StatusCode)}
	}

	// Читаем тело ответа
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return GeoAPIResult{Error: fmt.Errorf("ip.me read body failed: %w", err)}
	}

	bodyStr := strings.TrimSpace(string(bodyBytes))

	// Пробуем распарсить как JSON
	var jsonResult IPMeResponse
	if err := json.Unmarshal(bodyBytes, &jsonResult); err == nil && jsonResult.City != "" {
		return GeoAPIResult{City: jsonResult.City}
	}

	// Если не JSON, возможно это простой текст с IP или другой формат
	// Для ip.me может быть специфичный формат, попробуем найти город в тексте
	// Но по документации обычно это просто IP адрес, так что если нет JSON, возвращаем ошибку
	return GeoAPIResult{Error: fmt.Errorf("ip.me: unable to parse response, got: %s", bodyStr)}
}
