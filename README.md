# btc-billionaire

This project contains API services that accept BTC :wink:

PS: I didn't use any databases, just stored on the cache to be simple so if the program is re-run, data is in the memory will be lost.
## Steps to start

-  Clone the repo
  ```
  git clone https://github.com/salihkemaloglu/btc-billionaire.git
   ```
 
## Steps to run

1. Go to main directory:
   ```
   cd btc-billionaire/
   ```
2. Download `go` dependencies:
   ```
   go mod download
   ```
3. Run :
   ```
   go run cmd/main.go
   ```
 ## API endpoints
 - Save Records Sample Request
     ```
       curl -X POST \
        http://localhost:8080/wallet \
        -H 'cache-control: no-cache' \
        -H 'content-type: application/json' \
        -H 'postman-token: 91e2d545-e208-2b40-3715-640dc427a2cb' \
        -d '{
          "datetime": "2021-09-16T22:30:05+07:00",
          "amount": 150.52
      }'
     ```
    Response
   ```
   "OK"
   ```
  - History of Records Sample Request
     ```
      curl -X POST \
        http://localhost:8080/wallet/history \
        -H 'cache-control: no-cache' \
        -H 'content-type: application/json' \
        -H 'postman-token: 125639a9-927b-efd8-0f35-31293506238d' \
        -d '{
      "startDatetime": "2021-09-16T18:20:05+07:00",
      "endDatetime": "2021-09-16T23:25:05+07:00"
      }'
     ```
    Response
    ```
    [
        {
            "datetime": "2021-09-16T18:00:00+07:00",
            "amount": 1000
        },
        {
            "datetime": "2021-09-16T19:00:00+07:00",
            "amount": 1150.52
        }
    ]

## To improve
![architecture](https://freepngimg.com/thumb/street_fighter/35134-8-street-fighter-ii-image.png)