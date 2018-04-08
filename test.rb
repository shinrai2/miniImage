require 'ffi'
module Test
    extend FFI::Library
    ffi_lib 'miniImage.so'
    attach_function :rgb2GrayC, [:pointer, :pointer], :void
end

Test.rgb2GrayC("img/t.bmp", "output/t.bmp")