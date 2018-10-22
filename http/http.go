package http

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

//post with json
func PostWithJson(url string, body string) ([]byte, error) {
	return post(url, "application/json",body)
}

//post actual action
func post(url string,contentType string,bodyStr string) ([]byte,error) {

	//http包中实现了一个全局的DefaultClient
	resp, err := http.Post(url, contentType, strings.NewReader(bodyStr))

	if err != nil {
		fmt.Println("post error", err)
		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	return body, err
}
