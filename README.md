# GoFileProcessor
Allows a user to upload json file with a person data, the data gets stored in mongo db. The structure of file can be found in testfiles/sampleuserdata.json. 
Allows to query with first name for all users (case sensitive). 


## Resources: 

### Go and MongoDB 
https://www.mongodb.com/languages/golang

### Using Go with MongoDB 
https://www.mongodb.com/blog/post/quick-start-golang-mongodb-starting-and-setup

### local mongodb server 
https://docs.mongodb.com/guides/server/install/



### Some misc: 

Exit with error anywhere in go: 
https://stackoverflow.com/questions/18963984/exit-with-error-code-in-go

Channels and OS signal. 
https://gobyexample.com/signals



### Tests
go test ./data/... -v