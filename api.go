package main

import (
	"compress/zlib"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"time"

	"github.com/fatih/color"

	"bytes"
	"compress/gzip"
	"net/url"
	"strings"

	"github.com/andybalholm/brotli"

	tls_client "github.com/bogdanfinn/tls-client"

	http "github.com/bogdanfinn/fhttp"
	httputil "github.com/bogdanfinn/fhttp/httputil"
)

//var client http.Client

var consumption = 0.0

func main() {
	port := flag.String("port", "8082", "A port number (default 8082)")
	flag.Parse()
	fmt.Println("Hosting a TLS API on port " + *port)
	fmt.Println("Forked and changed this API. If you want to donate to the real creator --> https://paypal.me/carcraftz")
	http.HandleFunc("/", handleReq)
	err := http.ListenAndServe(":"+string(*port), nil)
	if err != nil {
		log.Fatalln("Error starting the HTTP server:", err)
	}
}

func getCookieStr(targetUrl string, client tls_client.HttpClient) string {
	parsed, _ := url.Parse(targetUrl)
	cookie := client.GetCookies(parsed)
	if len(cookie) > 1 {
		cookieString := ""
		for _, c := range cookie {
			cookieString += c.Name + "=" + c.Value + "; "
		}
		cookieString = cookieString[:len(cookieString)-2]
		return cookieString
	}
	return ""
}

func handleReq(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	// Ensure page URL header is provided
	pageURL := r.Header.Get("Poptls-Url")
	if pageURL == "" {
		http.Error(w, "ERROR: No Page URL Provided", http.StatusBadRequest)
		return
	}
	// Remove header to ignore later
	r.Header.Del("Poptls-Url")

	// Ensure user agent header is provided
	userAgent := r.Header.Get("User-Agent")
	if userAgent == "" {
		http.Error(w, "ERROR: No User Agent Provided", http.StatusBadRequest)
		return
	}

	//Handle Proxy (http://host:port or http://user:pass@host:port)
	proxy := r.Header.Get("Poptls-Proxy")
	if proxy != "" {
		r.Header.Del("Poptls-Proxy")
	}

	timeoutraw := r.Header.Get("Poptls-Timeout")
	timeout, err := strconv.Atoi(timeoutraw)
	if err != nil {
		//default timeout of 6
		timeout = 6
	}
	if timeout > 60 {
		http.Error(w, "ERROR: Timeout cannot be longer than 60 seconds", http.StatusBadRequest)
		return
	}

	options := []tls_client.HttpClientOption{
		tls_client.WithTimeout(timeout),
		tls_client.WithClientProfile(tls_client.Firefox_106),
		tls_client.WithProxyUrl(proxy),
	}

	client, err := tls_client.NewHttpClient(tls_client.NewNoopLogger(), options...)
	if err != nil {
		log.Fatal(err)
	}

	// Forward query params
	var addedQuery string
	for k, v := range r.URL.Query() {
		addedQuery += "&" + k + "=" + v[0]
	}

	endpoint := pageURL
	if len(addedQuery) != 0 {
		endpoint = pageURL + "?" + addedQuery
		if strings.Contains(pageURL, "?") {
			endpoint = pageURL + addedQuery
		} else if addedQuery != "" {
			endpoint = pageURL + "?" + addedQuery[1:]
		}
	}
	req, err := http.NewRequest(r.Method, ""+endpoint, r.Body)
	if err != nil {
		panic(err)
	}
	//master header order, all your headers will be ordered based on this list and anything extra will be appended to the end
	//if your site has any custom headers, see the header order chrome uses and then add those headers to this list
	masterheaderorder := []string{
		"host",
		"connection",
		"cache-control",
		"device-memory",
		"viewport-width",
		"rtt",
		"downlink",
		"ect",
		"sec-ch-ua",
		"sec-ch-ua-mobile",
		"sec-ch-ua-full-version",
		"sec-ch-ua-arch",
		"sec-ch-ua-platform",
		"sec-ch-ua-platform-version",
		"sec-ch-ua-model",
		"upgrade-insecure-requests",
		"user-agent",
		"accept",
		"sec-fetch-site",
		"sec-fetch-mode",
		"sec-fetch-user",
		"sec-fetch-dest",
		"referer",
		"accept-encoding",
		"accept-language",
		"cookie",
	}
	headermap := make(map[string]string)
	//TODO: REDUCE TIME COMPLEXITY (This code is very bad)
	headerorderkey := []string{}
	for _, key := range masterheaderorder {
		for k, v := range r.Header {
			lowercasekey := strings.ToLower(k)
			if key == lowercasekey {
				headermap[k] = v[0]
				headerorderkey = append(headerorderkey, lowercasekey)
			}
		}

	}
	for k, v := range req.Header {
		if _, ok := headermap[k]; !ok {
			headermap[k] = v[0]
			headerorderkey = append(headerorderkey, strings.ToLower(k))
		}
	}

	//ordering the pseudo headers and our normal headers
	req.Header = http.Header{
		http.HeaderOrderKey:  headerorderkey,
		http.PHeaderOrderKey: {":method", ":authority", ":scheme", ":path"},
	}
	//set our Host header
	u, err := url.Parse(endpoint)
	if err != nil {
		panic(err)
	}
	//append our normal headers
	for k := range r.Header {
		if k != "Content-Length" && !strings.Contains(k, "Poptls") {
			v := r.Header.Get(k)
			req.Header.Set(k, v)
		}
	}
	req.Header.Set("Host", u.Host)
	resp, err := client.Do(req)

	if err != nil {
		fmt.Printf("[%s][%s][%s]\r\n", color.YellowString("%s", time.Now().Format("2006-01-02 15:04:05")), color.BlueString("%s", pageURL), color.RedString("Connection Failed"))
		hj, ok := w.(http.Hijacker)
		if !ok {
			panic(err)
		}
		conn, _, err := hj.Hijack()
		if err != nil {
			panic(err)
		}
		if err := conn.Close(); err != nil {
			panic(err)
		}

		requestBytes, _ := httputil.DumpRequest(req, true)
		kiloBytes := float64(len(requestBytes)) / 1000
		consumption += kiloBytes
		return
	}
	defer resp.Body.Close()

	//req.Close = true
	//forward response headers
	for k, v := range resp.Header {
		if k != "Content-Length" && k != "Content-Encoding" {
			for _, kv := range v {
				w.Header().Add(k, kv)
			}
		}
	}
	w.Header().Add("session-cookies", getCookieStr(pageURL, client))
	w.WriteHeader(resp.StatusCode)
	var status string
	if resp.StatusCode > 302 {
		status = color.RedString("%s", resp.Status)
	} else {
		status = color.GreenString("%s", resp.Status)
	}

	requestBytes, _ := httputil.DumpRequest(req, true)
	responseBytes, _ := httputil.DumpResponse(resp, true)

	kiloBytes := float64(len(requestBytes)+len(responseBytes)) / 1000
	consumption += kiloBytes

	fmt.Printf("[%s][%s][%s][%.2f kB][%.2f kB]\r\n", color.YellowString("%s", time.Now().Format("2006-01-02 15:04:05")), color.BlueString("%s", pageURL), status, kiloBytes, consumption)

	//forward decoded response body
	//encoding := resp.Header["Content-Encoding"]
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	finalres := string(body)
	/*if len(encoding) > 0 {
		if encoding[0] == "gzip" {
			unz, err := gUnzipData(body)
			if err != nil {
				panic(err)
			}
			finalres = string(unz)
		} else if encoding[0] == "deflate" {
			unz, err := enflateData(body)
			if err != nil {
				panic(err)
			}
			finalres = string(unz)
		} else if encoding[0] == "br" {
			unz, err := unBrotliData(body)
			if err != nil {
				panic(err)
			}
			finalres = string(unz)
		} else {
			fmt.Println("UNKNOWN ENCODING: " + encoding[0])
			finalres = string(body)
		}
	} else {
		finalres = string(body)
	}*/
	if _, err := fmt.Fprint(w, finalres); err != nil {
		log.Println("Error writing body:", err)
	}
}

func gUnzipData(data []byte) (resData []byte, err error) {
	gz, _ := gzip.NewReader(bytes.NewReader(data))
	defer gz.Close()
	respBody, err := ioutil.ReadAll(gz)
	return respBody, err
}
func enflateData(data []byte) (resData []byte, err error) {
	zr, _ := zlib.NewReader(bytes.NewReader(data))
	defer zr.Close()
	enflated, err := ioutil.ReadAll(zr)
	return enflated, err
}
func unBrotliData(data []byte) (resData []byte, err error) {
	br := brotli.NewReader(bytes.NewReader(data))
	respBody, err := ioutil.ReadAll(br)
	return respBody, err
}
