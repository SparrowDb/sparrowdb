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