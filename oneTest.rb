require_relative "miniImage"

timeDiff = Time.now.to_f # Record start time
imgA = MiniImage::Image.new("img/t.bmp")
imgA.toGray
imgA.save("output/t.bmp")
imgA.release # Make sure the object is released, avoid memory leaks.
timeDiff = Time.now.to_f - timeDiff # Record finish time
printf("Test completed. Total %.4f ms.\n", timeDiff * 1000)
