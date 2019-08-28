package main

import (
	"os"
	"io/ioutil"
	"fmt"
	"net/http"
	"flag"
	"time"
	"strconv"
	"encoding/hex"
	"crypto/sha256"
	"crypto/hmac"
)

var (
	API_KEY string
	SECRET_KEY string
)

func die(s string, msg ...interface{}) {
	fmt.Fprintf(os.Stderr, s + "\n" , msg...)
	os.Exit(1)
}

func goincheck() error {
	var base_url string = "https://coincheck.com/"

	if err := get_rate(base_url); err != nil {
		return err
	}
	if err := get_history(base_url); err != nil {
		return err
	}
	if err := get_balance(base_url); err != nil {
		return err
	}
	return nil
}

func get_rate(base_url string) error {
	url := base_url + "api/rate/btc_jpy"
	ret, err := request("GET", url, nil)
	if err != nil {
		return err
	}
	fmt.Printf("rate : %s\n", string(ret))
	return nil
}

func get_balance(base_url string) error {
	url := base_url + "api/accounts/balance"
	ret, err := request("GET", url, nil)
	if err != nil {
		return err
	}
	fmt.Printf("balance : %s\n", string(ret))
	return nil
}

func get_history(base_url string) error {
	url := base_url + "api/exchange/orders/transactions_pagination"
	ret, err := request("GET", url, nil)
	if err != nil {
		return err
	}
	fmt.Printf("hist : %s\n", string(ret))
	return nil
}

func request(method string, url string, body []byte) ([]byte, error) {
	if body == nil {
		body = []byte{}
	}
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}
	t := time.Now().UnixNano() / (int64(time.Millisecond) / int64(time.Nanosecond))
	an := strconv.FormatInt(t, 10)
	sig := genhmac(an + url + string(body), SECRET_KEY)

	req.Header.Set("ACCESS-KEY", API_KEY)
	req.Header.Set("ACCESS-NONCE", an)
	req.Header.Set("ACCESS-SIGNATURE", sig)
	req.Header.Add("content-type", "application/json")
	req.Header.Add("cache-control", "no-cache")

	c := new(http.Client)
	res, err := c.Do(req)
	if err != nil {
		return nil, err
	}

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func genhmac(msg string, key string) string {
	mac := hmac.New(sha256.New, []byte(key))
	mac.Write([]byte(msg))
	return hex.EncodeToString(mac.Sum(nil))
}

func init() {
	var api_key string
	var secret_key string
	flag.StringVar(&api_key, "a", "", "your API KEY.")
	flag.StringVar(&secret_key, "s", "", "your SECRET KEY.")
	flag.Parse()

	if flag.NArg() < 0 {
		die("usage : goincheck -k <api key> -s <secret key>")
	}
	if api_key == "" {
		die("empty api key")
	}
	if secret_key == "" {
		die("empty secret_key")
	}
	API_KEY = api_key
	SECRET_KEY = secret_key
}

func main() {
	if err := goincheck(); err != nil {
		die("%s", err)
	}
}
