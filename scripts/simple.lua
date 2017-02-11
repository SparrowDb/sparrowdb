-- If image name contains gray, use grayscale effect
if string.match(imageCtx:name(), "gray") then
    imageCtx:grayscale()
end

-- If image name contains blue, use gaussian blur effect
if string.match(imageCtx:name(), "blur") then
    imageCtx:gaussianBlur(3.0)
end