# go-multilocker
This package aim to provide a funtionality to lock multiple resources at once using deadlock avoidance algorythms.

# Status
![](https://github.com/Arriven/go-multilocker/workflows/Go/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/Arriven/go-multilocker)](https://goreportcard.com/report/github.com/Arriven/go-multilocker)

# Usage
See [multilocker_test.go](https://github.com/Arriven/go-multilocker/blob/master/multilocker_test.go).

# Reason
For that rare cases where you need to acquire multiple resources at once and don't want to deal with all the scenarios of possible deadlocks and panics. I've just encountered such case in my code and decided to make a package for it. Unfortunatelly, [go-multilock](https://github.com/atedja/go-multilock) didn't work for me.
  
