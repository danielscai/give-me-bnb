package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

type AllProxies struct {
	All     []string  `json:"all"`
	History []History `json:"history"`
	Name    string    `json:"name"`
	Now     string    `json:"now"`
	Type    string    `json:"type"`
	UDP     bool      `json:"udp"`
}

type Proxy struct {
	History []History `json:"history"`
	Name    string    `json:"name"`
	Type    string    `json:"type"`
	UDP     bool      `json:"udp"`
}

type History struct {
	Time  time.Time `json:"time"`
	Delay int       `json:"delay"`
}

const (
	PROXY_BASE string = "http://127.0.0.1:9090/proxies/"
)

func clash(r run_single) error {

	client := &http.Client{}
	proxy_name := "✈️ 手动切换"
	escapeUrl := url.PathEscape(proxy_name)
	reqest_url := PROXY_BASE + escapeUrl
	// fmt.Print(reqest_url)
	resp, err := http.Get(reqest_url)
	if err != nil {
		fmt.Println("error")
		return err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Print("error reading body ")
		return err
	}

	var allProxy AllProxies
	err = json.Unmarshal(body, &allProxy)
	if err != nil {
		fmt.Println("error:", err)
		return err
	}
	// fmt.Print(allProxy.All)
	for _, proxy_name := range allProxy.All {
		fmt.Print(proxy_name)
		encoded := url.PathEscape(proxy_name)
		rurl := PROXY_BASE + encoded
		// fmt.Print(rurl)
		resp, err := http.Get(rurl)
		if err != nil {
			fmt.Println("error:", err)
			return err
		}
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("error:", err)
			return err
		}
		// fmt.Print(string(body))
		var proxy Proxy
		err = json.Unmarshal(body, &proxy)
		if err != nil {
			fmt.Println("error:", err)
			return err
		}
		// for _, history := range proxy.History {
		// 	fmt.Print(history.Delay, "\n")
		// }
		delay := proxy.History[0].Delay
		// fmt.Print(delay, "\n")
		if delay >= 1000 {
			fmt.Print("skip 1000 delay ", "\n")
			continue
		}
		request_data := map[string]string{"name": proxy_name}

		json_request, err := json.Marshal(request_data)
		if err != nil {
			panic(err)
		}
		req, err := http.NewRequest(http.MethodPut, reqest_url, bytes.NewBuffer(json_request))
		// fmt.Print(rurl, string(json_request))
		if err != nil {
			return err
		}
		resp, err = client.Do(req)
		if err != nil {
			return err
		}
		if resp.StatusCode == 204 {
			fmt.Print("switch proxy success. do get bnb... ")
			err = r()
			if err != nil {
				fmt.Print(err, "\n")
				// return err
				continue
			}

		} else {
			fmt.Print(resp.StatusCode)
		}
	}
	return nil
}
