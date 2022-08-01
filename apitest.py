import requests
import json

s = requests.Session()

def makeReq(url, headers={}, payload={}):
    if(payload != {}):  #if payload setted do post, if not do get
        res = s.post(url, headers=headers, data=payload)
        
        #append session cookies to session:
        sescookies = res.headers["session-cookies"].split('; ')
        for x in range(len(sescookies)):
            domain = url.split('://')[1]
            if '/' in domain:
                domain = domain.split('/')[0]
            s.cookies.set(sescookies[x].split('=')[0], sescookies[x].split('=')[1], domain=domain)
    else:
        res = s.get(url, headers=headers)
        
        #append session cookies to session:
        sescookies = res.headers["session-cookies"].split('; ')
        for x in range(len(sescookies)):
            domain = url.split('://')[1]
            if '/' in domain:
                domain = domain.split('/')[0]
            s.cookies.set(sescookies[x].split('=')[0], sescookies[x].split('=')[1], domain=domain)
    return res

headersGet = {
    'Poptls-Url': 'https://en.zalando.de/men-home/?_rfl=de',
    'accept': 'text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9',
    'accept-encoding': 'gzip, deflate, br',
    'accept-language': 'de,en-US;q=0.9,en;q=0.8',
	'cache-control': 'max-age=0',
    'sec-ch-ua': '".Not/A)Brand";v="99", "Google Chrome";v="103", "Chromium";v="103"',
    'sec-ch-ua-mobile': '?0',
    'sec-ch-ua-platform': '"Windows"',
    'sec-fetch-dest': 'document',
    'sec-fetch-mode': 'navigate',
    'sec-fetch-site': 'none',
    'sec-fetch-user': '?1',
    'upgrade-insecure-requests': '1',
    'user-agent': 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/103.0.0.0 Safari/537.36'
}

headersPost = {
    'Poptls-Url': 'https://httpbin.org/post',
    'accept': 'text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9',
    'accept-encoding': 'gzip, deflate, br',
    'accept-language': 'de,en-US;q=0.9,en;q=0.8',
    'content-type': 'application/json',
	'cache-control': 'max-age=0',
    'sec-ch-ua': '".Not/A)Brand";v="99", "Google Chrome";v="103", "Chromium";v="103"',
    'sec-ch-ua-mobile': '?0',
    'sec-ch-ua-platform': '"Windows"',
    'sec-fetch-dest': 'document',
    'sec-fetch-mode': 'navigate',
    'sec-fetch-site': 'none',
    'sec-fetch-user': '?1',
    'upgrade-insecure-requests': '1',
    'user-agent': 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/103.0.0.0 Safari/537.36'
}

postData = {
    "user": "test",
    "password": "test"
}

postReq = makeReq('http://localhost:8082', headers=headersPost, payload=json.dumps(postData))
getReq = makeReq('http://localhost:8082', headers=headersGet)
input('Press enter to close...')