# Data Enricher & Dispatcher

## Test assignment
This test task demonstrates the work of the service to receive a list of users from the API A service and send data to API B under certain conditions.
Condition: 
1. If the user's email ends with ".biz" - send to API B.
2. If the POST request is unsuccessful - make 3 attempts with an interval.

The full description of the task is contained in "description.txt"

Parameters can be set in the .env file.
If parameters are not specified or .env file is missing, values ​​will be set by default
### Example of .env file
```
GET_URL="https://jsonplaceholder.typicode.com/users"
POST_URL="https://webhook.site"
RETRY_DELAY=1s
RETRY_CNT=3
POST_TIMEOUT=110ms
```
### To work with the program, you can use commands from Makefile.
run app
```
make run
```
run Test
```
make test
```
### To cancel requests on timeout, a context is used. It can also be set in the POST_TIMEOUT .env
The context is applied to each request separately.
### Logging of query execution results is performed in the console, as well as in the log file "logs/app.log"