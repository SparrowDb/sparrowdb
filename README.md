﻿Whats is SparrowDB?
====================
SparrowDB is an image database that works like an append-only object store. Sparrow has tools that allow image processing and HTTP server to access images.


Sparrow Object Store
====================
Sparrow consists of three files – the actual Sparrow store file containing the images data, plus an index file and a bloom filter file.

There is a corresponding data definition record followed by the image bytes for each image in the storage file. The index file provides the offset of the data definition in the storage file.


Features
====================
1. Built-in HTTP API so you don't have to write any server side code to get up and running.
2. Optimizations for image storing.
3. Websocket server to provide real time information about the server.


Getting started
====================
This short guide will walk you through getting a basic server up and running, and demonstrate some simple reads and writes.



Using Sparrow
====================
Creating a database:
	
	curl -X POST -d '{"type":"create_database", "params":{"name":"database_name"}}' http://127.0.0.1:8081/query

Show databases:

    curl -X POST -d '{"type":"show_databases"}' http://127.0.0.1:8081/query


Sending an image to database:

	curl -i -X POST -H "Content-Type: multipart/form-data"  \
        -F "uploadfile=@image.jpg" \
        -F "dbname=database_name" \
        -F "key=image_name" \
        http://127.0.0.1:8081/upload


Accessing image from browser:
	
	http://localhost:8081/g/database_name/image_key

License
====================
This software is under MIT license.
