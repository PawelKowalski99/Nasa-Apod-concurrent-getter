# Gogoapps recruitment task - Golang Developer

##Code author: Paweł Kowalski

## Usage
**Docker usage**
  * ```docker build -t gogoapps .```
  * ```docker run -p 8080:8080 --env-file .env -p (type_here_port_from_env_file):(type_here_port_from_env_file) gogoapps```
    
    
## Solution description

### Used external dependencies
* go-chi - small http router, which is fast and easy to use.
   * Additionally, I used middleware request throttle to limit and add to queue requesters
* go-chi/render - served me as an easy response renderer to json
* gjson - Served me as easy to use json key-value getter for query filtering
* storj.io/common - provided me a limiter struct to have concurrency

### Additional resources
[Storj.io article about concurrency and primitives](https://www.storj.io/blog/production-concurrency#ss)