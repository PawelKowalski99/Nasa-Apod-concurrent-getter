# Gogoapps recruitment task - Golang Developer

##Code author: Paweł Kowalski

## Usage
**Docker usage**
  * ```docker build -t gogoapps .```
  * ```docker run --env-file .env -p (type_here_port_from_env_file):(type_here_port_from_env_file) gogoapps```
    
    
** Golang usage **
 * ```go run .```

## Solution description

### Used external dependencies
* go-chi - small http router, which is fast and easy to use.
   * Additionally, I used middleware request throttle to limit and add to queue requesters
* go-chi/render - served as an easy response renderer to json
* gjson - Served as a easy to use json key-value getter for query filtering
* storj.io/common - provides a limiter struct to have concurrency
* viper - provides easily configurable environment variables

### Additional resources
[Storj.io article about concurrency and primitives](https://www.storj.io/blog/production-concurrency#ss)
