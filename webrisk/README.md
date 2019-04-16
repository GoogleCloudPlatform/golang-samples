# Reference Implementation for the Usage of Google Cloud WebRisk APIs (Beta)

The `webrisk` Go package can be used with the
[Google Cloud WebRisk APIs (Beta)](https://cloud.google.com/web-risk/)
to access the Google Cloud WebRisk lists of unsafe web resources. Inside the
`cmd` sub-directory, you can find two programs: `wrlookup` and `wrserver`. The
`wrserver` program creates a proxy local server to check URLs and a URL
redirector to redirect users to a warning page for unsafe URLs. The `wrlookup`
program is a command line service that can also be used to check URLs.

This **README.md** is a quickstart guide on how to build, deploy, and use the
`[WebRisk](https://godoc.org/cloud.google.com/go/webrisk/apiv1beta1)` Go package. It can be used out-of-the-box. The GoDoc and API
documentation provide more details on fine tuning the parameters if desired.


# How to Build

To download and install from the source, run the following command:

```
go get github.com/GoogleCloudPlatform/golang-samples/webrisk
```

The programs below execute from your `$GOPATH/bin` folder.
Add that to your `$PATH` for convenience:

```
export PATH=$PATH:$GOPATH/bin
```

The program expects an API key as a parameter, export it with the following
command for later use:

```
export APIKEY=Your Api Key
```

# Proxy Server

The `wrserver` server binary runs a WebRisk API lookup proxy that allows
users to check URLs via a simple JSON API.

1.	Once the Go environment is setup, run the following command with your API key:

	```
	go get github.com/GoogleCloudPlatform/golang-samples/webrisk/cmd/wrserver
	wrserver -apikey $APIKEY
	```

	With the default settings this will start a local server at **127.0.0.1:8080**.

2.  The server also uses an URL redirector (listening on `/r`) to show an interstitial for anything marked unsafe.  
If the URL is safe, the client is automatically redirected to the target. Else, an interstitial warning page is shown as recommended by Web Risk.  
Try these URLs:

	```
	127.0.0.1:8080/r?url=http://testsafebrowsing.appspot.com/apiv4/ANY_PLATFORM/MALWARE/URL/
	127.0.0.1:8080/r?url=http://testsafebrowsing.appspot.com/apiv4/ANY_PLATFORM/SOCIAL_ENGINEERING/URL/
	127.0.0.1:8080/r?url=http://testsafebrowsing.appspot.com/apiv4/ANY_PLATFORM/UNWANTED_SOFTWARE/URL/
	127.0.0.1:8080/r?url=http://www.google.com/
	```

3.	The server also has a lightweight implementation of the API v4 threatMatches endpoint.  
To use the local proxy server to check a URL, send a POST request to `127.0.0.1:8080/v1beta1/urs:search` with the following JSON body:

	```json
	{
          "uri":"http://testsafebrowsing.appspot.com/apiv4/ANY_PLATFORM/MALWARE/URL/",
          "threatTypes":[
            "MALWARE"
          ]
        }
	```

# Command-Line Lookup

The `wrlookup` command-line binary is another example of how the Go Safe
Browsing library can be used to protect users from unsafe URLs. This
command-line tool filters unsafe URLs piped via STDIN. Example usage:

```
$ go get github.com/GoogleCloudPlatform/golang-samples/webrisk/cmd/wrlookup
$ echo "http://testsafebrowsing.appspot.com/apiv4/ANY_PLATFORM/MALWARE/URL/" | wrlookup -apikey=$APIKEY
```


# WebRisk System Test
To perform an end-to-end test on the package with the WebRisk backend,
run the following command:

```
go test github.com/GoogleCloudPlatform/golang-samples/webrisk -v -run TestWebriskClient -apikey $APIKEY
```
