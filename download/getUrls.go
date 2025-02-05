package download

import (
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"
)

func getGlbUrlsFromJsonUrl(jsonUrl string) []string {
	jsonByte := getJsonByte(jsonUrl)
	baseUrl := getBaseUrlFromJsonUrl(jsonUrl)
	glbUrls := getGlbUrlsFromJsonByte(baseUrl, jsonByte)
	return glbUrls
}

func getJsonByte(jsonUrl string) []byte {
	// 下載 JSON 檔案
	resp, err := http.Get(jsonUrl)
	if err != nil {
		log.Fatalf("Failed to download JSON file: %v\n", err)
	}
	defer resp.Body.Close()
	jsonByte, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed to read JSON file: %v\n", err)
	}
	return jsonByte
}

func getBaseUrlFromJsonUrl(jsonUrl string) string {
	// 取得 JSON 檔案的基礎 URL
	urlSplit := strings.Split(jsonUrl, "/")
	baseUrl := strings.Join(urlSplit[:len(urlSplit)-1], "/")
	return baseUrl
}

func getGlbUrlsFromJsonByte(baseUrl string, jsonByte []byte) []string {
	text := string(jsonByte)
	glbUrls := []string{}
	// 使用正規表達式找出所有的 glb 檔案連結
	re := regexp.MustCompile(`"uri":\s*"([^"]+\.glb)"`)
	matches := re.FindAllStringSubmatch(text, -1)
	for _, match := range matches {
		url := baseUrl + "/" + match[1]
		glbUrls = append(glbUrls, url)
	}
	return glbUrls
}
