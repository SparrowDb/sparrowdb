-- loads sparrowdb module
local s = require("sparrowdb")

image = s.getInputImage()

-- buf is global variable with image data
out = s.grayscale(image)

s.setOutputImage(out)