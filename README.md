# goparent
a go based infant data tracking system for my graduate program web engineering class


## requirements

1. running instance of rethinkdb with tables:
    - users
    - sleep
    - feeding
    - waste
2. gui front end [here](https://github.com/sasimpson/goparentgui)
3. go service running

to run do this:

    go get -u 
    go run 

or this:  

    go get -u 
    go build
    ./main


## ToDo:

- [ ] handle multiple kids
- [ ] add in some statistical analysis of data