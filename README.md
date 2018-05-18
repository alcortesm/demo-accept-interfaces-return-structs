# Accept Interfaces, return structs


# Tips

Quick way to setup a mongo server on localhost:

```
docker run --name mymongo -d --rm -p 127.0.0.1:27017:27017 mongo:3.2.20-jessie
```

Quick way to open a mongo client to the server:

```
; docker run -it --link mymongo --rm mongo:3.2.20-jessie mongo --host mymongo test
MongoDB shell version: 3.2.20
connecting to: mymongo:27017/test
Welcome to the MongoDB shell.
For interactive help, type "help".
For more comprehensive documentation, see
	http://docs.mongodb.org/
Questions? Try the support group
	http://groups.google.com/group/mongodb-user
2018-05-18T17:34:44.573+0000 I STORAGE  [main] In File::open(), ::open for '/home/mongodb/.mongorc.js' failed with errno:2 No such file or directory
> db.abbreviations.find({})
{ "_id" : ObjectId("5aff0ea842b7c6fe754d7a3f"), "abbreviation" : "NATO", "meaning" : "North Atlantic Treaty Organization" }
```

# References

[Jack Lindamood, What “accept interfaces, return structs” means in Go](https://medium.com/@cep21/what-accept-interfaces-return-structs-means-in-go-2fe879e25ee8)
