# OpenFaas-server
## Run Server
`go build main.go`
## Usage
`POST localhost:8080/run`
```
{
	"fid": "hello",
	"src": "def handler(x):\n\tprint(x)",
	"params": "ldodafl",
	"lang": "python"
}
```
Response
```
{
  "fid":"hello",
  "result":"ldodafl\n"
}
```
