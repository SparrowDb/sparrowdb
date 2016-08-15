Scripts
====================
SparrowDB uses [bild](https://github.com/anthonynsimon/bild) to allow image processing using [LUA](github.com/yuin/gopher-lua) script.

All SparrowDB scripts must be in 'scripts' folder.

Script Example
====================
Example of script that converts image to grayscale.


```lua
-- loads sparrowdb module
local s = require("sparrowdb")

-- get image data
image = s.getInputImage()

-- convert image to grayscale
out = s.grayscale(image)

-- set new output for image
s.setOutputImage(out)
```

