# Figura

This project is a simple configuration server written in GO that centralizes config files for applications. Essentially what it does is allow applications to read their config files (in YML or JSON) from a central location and where ideally these config files are under source control (git for example). Config files can optionally be secured with a key stored within the file for a basic level of access control. 

All communication is done via json, REST and HTTP.


## Installation

The simplest way if you 'make' installed is to run 'make' which will install all the dependencies and install the apps. Otherwise After cloning the repository you need to install the few dependencies. Execute the following within the main directory.

```bash
$ go get github.com/riomhaire/figura
$ cd <gopath-root>/src/github.com/riomhaire/figura
$ go get ./...
```

The best way if you just want run is to build and install the apps:

```bash
 go install github.com/riomhaire/figura/frameworks/application/figura
 go install 
```

## Design

The design is fairly simple ... we have a simple http server which dishes up files based on the rest URI:

```
   http[s]://<host>:<port>/api/v1/configuration/<application-name>
```

and a 'GET' request where <application-name> will map to a <application-name>.yml or <<application-name>.yaml> or <application-name>.json. The directory (in the default implimentation) used to read from is controlled by the '-configs' command line option.

For example the 'lightauth2' application could have a configuration file called 'lightauth2.yml' or 'lightauth2.json'.

When a request is make the appropriate potential file names are looked up and the first match returned. The order of precidence is 'yml' then 'yaml' then 'json' then 'xml'.

Security is fairly simple in that the GET request Authorization field is checked against the contents of a key file based off the naming convention <application-name>.key - so lightauth2 key would be stored within 'lightauth2.key'. If no key file is present then anyone can access that configuration file.
 

## Execution

After building and installation there is little to config other than the directory where the files are and the port to dish stuff out on. The 'figura --help' command returns:

```bash
./figura --help
2017/12/26 18:18:14 [INFO] Initializing
Usage of ./figura:
  -configs string
    	Directory here configurations stored (default "configs")
  -port int
    	Port to use (default 3050)
  -profile
    	Enable profiling endpoint

```

So to run using config files within the '/etc/figura' and on port '8080' you would execute the command:

```bash
figura --port=8080 --configs=/etc/figura
```

## DevOps Endpoints

Apart from the main get config file endpoint there are two others within figura:

```
   http[s]://<host>:<port>/api/v1/configuration/statistics
   http[s]://<host>:<port>/api/v1/configuration/health
```

The former used the negroni statistics plugin and returns data like:

```
{
	"pid": 1587,
	"uptime": "2m53.303479465s",
	"uptime_sec": 173.303479465,
	"time": "2017-12-26 18:33:18.088128661 +0000 GMT m=+173.304615563",
	"unixtime": 1514313198,
	"status_code_count": {},
	"total_status_code_count": {
		"200": 2,
		"404": 1
	},
	"count": 0,
	"total_count": 3,
	"total_response_time": "549.238µs",
	"total_response_time_sec": 0.000549238,
	"average_response_time": "183.079µs",
	"average_response_time_sec": 0.000183079
}
```
The latter is a simple endpoint which returns:

```
{
	"status": "up"
}
```

if figura is up and running
