require 'ffi'
module Test
    extend FFI::Library
    ffi_lib 'miniImage.so'
    attach_function :rgb2GrayC, [:string, :string], :void
    attach_function :verifyGrayC, [:string], :bool
end

# Test.rgb2GrayC("img/t.bmp", "output/t.bmp")
puts Test.verifyGrayC("img/tg.bmp")
