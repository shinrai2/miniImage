require 'ffi'
module CallLibrary
    extend FFI::Library
    ffi_lib 'miniImage.so'
    # attach_function :rgb2GrayC, [:string, :string], :void
    # attach_function :verifyGrayC, [:string], :bool
    # attach_function :moveBoundsC, [:string, :string, :long_long, :long_long, :long_long, :long_long, :uchar, :uchar, :uchar, :uchar], :void
    attach_function :exportInitialize, [:string], :long_long
    attach_function :exportSave, [:long_long, :string], :void
end

module MiniImage
    class Image
        def initialize(path)
            @keyOfMap = CallLibrary.exportInitialize(path)
        end

        def save(path)
            CallLibrary.exportSave(@keyOfMap, path)
        end

        def release()
        end
    end
end
