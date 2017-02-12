<img src="http://golang.org/doc/gopher/frontpage.png" alt="Golang logo" align="right"/>

[![GoDoc](https://godoc.org/github.com/SparrowDb/sparrowdb?status.svg)](https://godoc.org/github.com/SparrowDb/sparrowdb)
[![Build Status](https://travis-ci.org/SparrowDb/sparrowdb.svg?branch=master)](https://travis-ci.org/SparrowDb/sparrowdb)
[![Go Report Card](https://goreportcard.com/badge/github.com/SparrowDb/sparrowdb)](https://goreportcard.com/report/github.com/SparrowDb/sparrowdb)

Whats is SparrowDB?
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
	
	curl -X PUT http://127.0.0.1:8081/api/database_name

Show databases:

    curl -X GET http://127.0.0.1:8081/api/_all


Sending an image to database:

	curl -i -X PUT -H "Content-Type: multipart/form-data"  \
        -F "uploadfile=@image.jpg" \
        http://127.0.0.1:8081/api/database_name/image_key


Querying an image:

	curl -X GET http://127.0.0.1:8081/api/database_name/image_key


Accessing image from browser:
	
	http://localhost:8081/g/database_name/image_key


Token
====================

If is set in database configuration file, generate_token = true, SparrowDB will generate a token for each image uploaded. The token’s value is randomly assigned by and stored in database. The token effectively eliminates attacks aimed at guessing valid URLs for photos.

Accessing image from browser with token:
	
	http://localhost:8081/g/database_name/image_key/token_value


Image Processing
====================

SparrowDB uses [bild](https://github.com/anthonynsimon/bild) to allow image processing using [LUA](https://github.com/yuin/gopher-lua) script.

All SparrowDB scripts must be in 'scripts' folder.

Example of script with image effect:

```lua
-- If image name contains gray, use grayscale effect
if string.match(imageCtx:name(), "gray") then
    imageCtx:grayscale()
end

-- If image name contains blue, use gaussian blur effect
if string.match(imageCtx:name(), "blur") then
    imageCtx:gaussianBlur(3.0)
end
```

Example of script with pixel iteration:

```lua
-- get image bounds
b = imageCtx:bounds()

-- create an editable image
p = sparrowRGBA.new(imageCtx)

-- iterate over pixels
for i = 0, b['width'] do
    for j = 0, b['height'] do
        -- get current pixel color: red, green, blue, alpha
        v = p:getPixel(i, j)

        -- set current pixel color: red, green, blue, alpha
        p:setPixel(i, j, v['red'], v['green'], v['blue'], v['alpha'])
    end
end

-- set processed image as outputsss
imageCtx:setOutput(p)
```


License
====================
This software is under MIT license.
