# phttprobe
Take a list of domain + port and probe for working http and https servers.

## Usage:
```
Usage of ./phttprobe:
  -c int
        set the concurrency level (default 20)
  -i int
        Index of ports (default 1)
  -t int
        timeout (second) (default 10)
  -v    output errors to stderr
```

## Example:
#### `input1.txt` (Format domain,ports...)
```
google.com,2087,2086,8880,2082,443,80,2052,2096,2083,8080,8443,2095,2053
yahoo.com,2087,2086,8880,2082,443,80,2052,2096,2083,8080,8443,2095,2053
sport.yahoo.com,2086,443,2096,2053,8080,2082,80,2083,8443,2052,2087,2095,8880
```
##### Command
```
cat input.txt | ./phttprobe -c 20 -t 10
```

#### `input2.txt` (Format domain,ip,ports...)
```
google.com,192.154.234.3,2087,2086,8880,2082,443,80,2052,2096,2083,8080,8443,2095,2053
yahoo.com,192.154.234.3,2087,2086,8880,2082,443,80,2052,2096,2083,8080,8443,2095,2053
sport.yahoo.com,192.154.234.3,2086,443,2096,2053,8080,2082,80,2083,8443,2052,2087,2095,8880
```
##### Command
```
cat input.txt | ./phttprobe -c 20 -t 10 -i 2
```