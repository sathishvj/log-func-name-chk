## Purpose
Check for standardized log statements in go code.  

Our standard logging required that the statements contain the go file name, line number, and the function it is logging from. 
Eg. 
email.go:117:initEmail(): Email config initialized.
main.go:117:main(): Started http server, Listening at 8080...

There are inbuilt ways we can get the file and line number into the log statements.  But there was no way I found to get the function automatically added.

This program scans a folder for .go files and then checks that a log statement is prefixed with a matching function name.  
As of now, it is very specific to my requirements. And these are the matchers.  

```
var funcBegin = regexp.MustCompile(`func.* (\w+)\(.*\{$`)
var funcEnd = regexp.MustCompile(`^\}.*$`)
var logPrefixRegex = regexp.MustCompile(`\s*logger\..*\(.*$`)
```

## Working
log-func-name-chk <dir>

* This will get all the log files in the directory ...
* then check that the log statements are of a certain desired format.

## Future
* allow configuration of begin, end, log statement
* allow configuration of exceptions
* allow interactivity and ability to choose which all statements are changed
* allow quite mode that does everything automatically
