# TLS-Fingerprint-API

A server that proxies requests and uses my fork of cclient & fhttp (fork of net/http) to prevent your requests from being fingerprinted. Built on open source software, this repo is a simple yet effective solution to companies violating your privacy. It uses cclient to spoof tls fingerprints, and fhttp to enable mimicry of chrome http/2 connection settings, header order, pseudo header order, and enable push.

## Support

Forked and changed this API. If you want to donate to the real creator --> https://paypal.me/carcraftz

## How to use:

*Important* note: I modified the api and the libs it is using so that you get the exact tls fingerprint as the latest Chrome browser (You can check that on this website: https://tls.peet.ws/api/all). Make sure to always specify your User Agent as a Chrome User Agent.

*Important* note: If you're using this in a language other than go, then use this repo. If you're using this in go, then take a look into the "godirect.test" file. In this file should be all you need. Rename it to "godirect.go", put it in an other directory, run this command: "go mod init SOMENAME" and install the necessary packages.

Deploy this server somewhere. Localhost is preferrable to reduce latency. The go source code is given if you want to build it yourself on any platform (windows, macos, linux). I attached a prebuilt windows exe.

Modify your code to make requests to the server INSTEAD of the endpoint you want to request. Ex: If running on localhost, make requests to http://127.0.0.1:8082. Make sure to also remove any code that uses a proxy in the request.

Add the request header "poptls-url", and set it equal to the endpoint you want to request. For example, if you want to request https://httpbin.org/get, you would add the header "poptls-url" = "https://httpbin.org/get"

You will find the session cookies in the response-header "session-cookies".

Optional: Add the request header "poptls-proxy" and set it equal to the URL for the proxy you want to use (format: http://user:pass@host:port or http://host:port). This will make the server use your proxy for the request.

Optional: Add the request header "poptls-allowredirect" and set it to true or false to enable/disable redirects. Redirects are enabled by default.

Optional: Add the request header "poptls-timeout" and set it to an integer (in seconds) to specify the max timeout to wait for a request.



## Run on a different Port:

By default the program runs on port 8082. You can specify another port by passing a flag --port=PORTNUM

## Examples:

### GO

Take a look into the "godirect.test" file.

### Python

Take a look into the "apitest.py" file.

### Node.js

To call this in node.js, lets say with node-fetch, you could do

````
fetch("http://localhost:8082",{
headers:{
"poptls-url":"https://httpbin.org/get",
"poptls-proxy":"https://user:pass@ip:port", //optional
"poptls-allowredirect:"true" //optional (TRUE by default)
}
})```
````
