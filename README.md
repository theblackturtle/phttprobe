# phttprobe
Take a list of domain + port and probe for working http and https servers.

## Usage:
```
Usage of ./phttprobe:
  -c int
        set the concurrency level (default 20)
  -t int
        timeout (second) (default 10)
  -v    output errors to stderr
```

## Example:
##### `input.txt`
```
google.com,2087,2086,8880,2082,443,80,2052,2096,2083,8080,8443,2095,2053
yahoo.com,2087,2086,8880,2082,443,80,2052,2096,2083,8080,8443,2095,2053
sport.yahoo.com,2086,443,2096,2053,8080,2082,80,2083,8443,2052,2087,2095,8880
```
##### Command
```
cat input.txt | ./phttprobe -c 20 -t 10
```