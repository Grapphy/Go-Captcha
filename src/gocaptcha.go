package gocaptcha

import (
	"encoding/json"
	"net/http"
	"net/url"
  "fmt"
	"time"
)

const CaptchaIn string = "https://2captcha.com/in.php"
const CaptchaRes string = "https://2captcha.com/res.php"

type ApiResponse struct {
	Status       int    `json:"status"`
	Request      string `json:"request"`
	ErrorMessage string `json:"error_message"`
}

type TwoCaptcha struct {
	Key     string
	Timeout int
}

func NewTwoCaptcha(key string) TwoCaptcha {
	return TwoCaptcha{Key: key, Timeout: 150}
}

func (t TwoCaptcha) SolveRecaptcha(siteKey, pageURL string) (string, error) {
	urlParams := url.Values{}
	urlParams.Set("key", t.Key)
	urlParams.Set("json", "1")
	urlParams.Set("method", "userrecaptcha")
	urlParams.Set("googlekey", siteKey)
	urlParams.Set("pageurl", pageURL)
	inURL := CaptchaIn + "?" + urlParams.Encode()

	res, err := http.Get(inURL)
	if err != nil {
		return "", err
	}

	apiResponse := ApiResponse{}
	json.NewDecoder(res.Body).Decode(&apiResponse)
	res.Body.Close()

	if apiResponse.Status == 0 {
		return "", fmt.Errorf(apiResponse.ErrorMessage)
	}

	urlParams = url.Values{}
	urlParams.Set("key", t.Key)
	urlParams.Set("json", "1")
	urlParams.Set("action", "get")
	urlParams.Set("id", apiResponse.Request)
	resURL := CaptchaRes + "?" + urlParams.Encode()
	var counter int = t.Timeout / 5

	for i := 0; i < counter; i++ {
		res, err := http.Get(resURL)
		if err != nil {
			return "", err
		}

		apiResponse = ApiResponse{}
		json.NewDecoder(res.Body).Decode(&apiResponse)
		res.Body.Close()

		if apiResponse.Status != 0 {
			return apiResponse.Request, nil
		}

		if apiResponse.ErrorMessage != "" {
			return "", fmt.Errorf(apiResponse.ErrorMessage)
		}

		time.Sleep(5 * time.Second)
	}

	return "", fmt.Errorf("Timeout!")
}
