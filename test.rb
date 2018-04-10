require 'ffi'
module Test
    extend FFI::Library
    ffi_lib 'miniImage.so'
    attach_function :rgb2GrayC, [:string, :string], :void
    attach_function :verifyGrayC, [:string], :bool
    attach_function :moveBoundsC, [:string, :string, :long_long, :long_long, :long_long, :long_long, :uchar, :uchar, :uchar, :uchar], :void
end

# Test.rgb2GrayC("img/t.bmp", "output/t.bmp")
# puts Test.verifyGrayC("img/tg.bmp")
Test.moveBoundsC("img/t.bmp", "output/t.bmp", 50, 50, 50, 50, 0, 0, 0, 255)
