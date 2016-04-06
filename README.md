*This was a Go and SQLite web app written for my CS 457 Database Management Systems class opposed to doing it in PHP and MySQL. I may or may not come back to this and improve it slightly.*

#Setup

To build and run do the following commands:

##Build
```
go get github.com/mattn/go-sqlite3
go build createDB.go && createDB
go build airship.go
```

##Run
```
airship
```

If you ever wish to reset the database just rerun `createDB`

#Features

* View all tables
* Insert into a table
* Use dropdown queries on two different tables