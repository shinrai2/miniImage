require_relative "miniImage"

timeDiff = Time.now.to_f # Record start time
# imgA = MiniImage::Image::loadFrom("img/t.bmp")
colorGray = MiniImage::Color.new(128, 128, 128, 255)
colorWhite = MiniImage::Color.new(255, 255, 255, 255)
imgA = MiniImage::Image::newBlank(200, 200, colorGray)
fontA = MiniImage::Font::loadFrom("fonts/timesi.ttf")
fontA.drawString(imgA, 100, 15, 100, "nice", colorWhite)
# imgA.toGray
imgA.save("output/t.bmp")
timeDiff = Time.now.to_f - timeDiff # Record finish time
printf("Test completed. Total %.4f ms.\n", timeDiff * 1000)
print("Exit :)\n")
