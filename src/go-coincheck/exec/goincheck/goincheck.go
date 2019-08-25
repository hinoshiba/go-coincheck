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
	return nil
}

func get_rate(base_url string) error {
	url := base_url + "/api/rate/btc_jpy"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	c := new(http.Client)
	res, err := c.Do(req)
	if err != nil {
		return err
	}

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	fmt.Printf("rate : %s\n", string(b))
	return nil
}

func get_balance(base_url string) error {
	url := base_url + "api/accounts/balance"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	t := time.Now().UnixNano() / (int64(time.Millisecond) / int64(time.Nanosecond))
	an := strconv.FormatInt(t, 10)
	body := ""

	sig := genhmac(an + url + body, SECRET_KEY)

	req.Header.Set("ACCESS-KEY", API_KEY)
	req.Header.Set("ACCESS-NONCE", an)
	req.Header.Set("ACCESS-SIGNATURE", sig)
	req.Header.Add("content-type", "application/json")
	req.Header.Add("cache-control", "no-cache")


	c := new(http.Client)
	res, err := c.Do(req)
	if err != nil {
		return err
	}

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	fmt.Printf("balance : %s\n", string(b))
	return nil
}

func get_history(base_url string) error {
	url := base_url + "api/exchange/orders/transactions_pagination"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	t := time.Now().UnixNano() / (int64(time.Millisecond) / int64(time.Nanosecond))
	an := strconv.FormatInt(t, 10)
	body := ""

	sig := genhmac(an + url + body, SECRET_KEY)

	req.Header.Set("ACCESS-KEY", API_KEY)
	req.Header.Set("ACCESS-NONCE", an)
	req.Header.Set("ACCESS-SIGNATURE", sig)
	req.Header.Add("content-type", "application/json")
	req.Header.Add("cache-control", "no-cache")


	c := new(http.Client)
	res, err := c.Do(req)
	if err != nil {
		return err
	}

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	fmt.Printf("history : %s\n", string(b))
	return nil
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
