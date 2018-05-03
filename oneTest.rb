require_relative "miniImage"

timeDiff = Time.now.to_f # Record start time
# imgA = MiniImage::Image::loadFrom("img/t.bmp")
imgA = MiniImage::Image::newBlank(100, 100)
fontA = MiniImage::Font::loadFrom("fonts/times.ttf")
imgA.toGray
imgA.save("output/t.bmp")
timeDiff = Time.now.to_f - timeDiff # Record finish time
printf("Test completed. Total %.4f ms.\n", timeDiff * 1000)
print("Exit :)\n")
