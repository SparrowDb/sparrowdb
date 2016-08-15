-- loads sparrowdb module
local s = require("sparrowdb")

-- get image data
image = s.getInputImage()

-- convert image to grayscale
out = s.grayscale(image)

-- set new output for image
s.setOutputImage(out)