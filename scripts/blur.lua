-- loads sparrowdb module
local s = require("sparrowdb")

-- get image data
image = s.getInputImage()

-- apply blur effect
out = s.gaussianBlur(image, 3.0)

-- set new output for image
s.setOutputImage(out)