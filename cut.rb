require_relative "miniImage"

subpath = "sample_1999_38000"
timeDiff = Time.now.to_f # Record start time
colorWhite = MiniImage::Color.new(255, 255, 255, 255)
oriImg = MiniImage::Image::loadFrom("img/"+subpath+".png")

img0 = oriImg.moveBounds(0, 0, -512, -2048, colorWhite)
img0.save("output/"+subpath+"_0.png")
img1 = oriImg.moveBounds(-512, 0, 0, -2048, colorWhite)
img1.save("output/"+subpath+"_1.png")

img2 = oriImg.moveBounds(0, -512, -512, -1536, colorWhite)
img2.save("output/"+subpath+"_2.png")
img3 = oriImg.moveBounds(-512, -512, 0, -1536, colorWhite)
img3.save("output/"+subpath+"_3.png")

img4 = oriImg.moveBounds(0, -1024, -512, -1024, colorWhite)
img4.save("output/"+subpath+"_4.png")
img5 = oriImg.moveBounds(-512, -1024, 0, -1024, colorWhite)
img5.save("output/"+subpath+"_5.png")

img6 = oriImg.moveBounds(0, -1536, -512, -512, colorWhite)
img6.save("output/"+subpath+"_6.png")
img7 = oriImg.moveBounds(-512, -1536, 0, -512, colorWhite)
img7.save("output/"+subpath+"_7.png")

img8 = oriImg.moveBounds(0, -2048, -512, 0, colorWhite)
img8.save("output/"+subpath+"_8.png")
img9 = oriImg.moveBounds(-512, -2048, 0, 0, colorWhite)
img9.save("output/"+subpath+"_9.png")

timeDiff = Time.now.to_f - timeDiff # Record finish time
printf("Test completed. Total %.4f ms.\n", timeDiff * 1000)
print("Exit :)\n")