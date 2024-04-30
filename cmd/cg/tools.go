package main

import (
	"fmt"
	"net/url"

	"github.com/go-resty/resty/v2"
)

// CheckURL checks valid urlString
func CheckURL(urlString string) (*url.URL, error) {
    parsedURL, err := url.Parse(urlString)  
    if err != nil {  
        return nil, err
    }  
    
    if parsedURL.Scheme != "" && parsedURL.Host != "" {  
        return parsedURL, nil
    } else {  
        return nil, fmt.Errorf("url str has invalid URL schema or Host")
    }  
}

// Request request by url
func Request(url string, rawData any) error {
    _, err := CheckURL(url)
    if err != nil {
        return fmt.Errorf("[%s] is not correct url", url)
    }
    client := resty.New()
    _, err = client.R().SetResult(rawData).Get(url)
    if err != nil {
        return err
    }
    return nil
}
