require 'ffi'
module CallLibrary
    extend FFI::Library
    ffi_lib 'miniImage.so'
    attach_function :exportInitialize, [:string], :long_long
    attach_function :exportSave, [:long_long, :string], :void
    attach_function :exportRelease, [:long_long], :void
    attach_function :exportIsGray, [:long_long], :bool
    attach_function :exportToGray, [:long_long], :void
    attach_function :exportToRgba, [:long_long], :void
    attach_function :exportMoveBounds, [:long_long, :long_long, :long_long, :long_long, :long_long, :uchar, :uchar, :uchar, :uchar], :void
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
            CallLibrary.exportRelease(@keyOfMap)
        end

        def isGray()
            return CallLibrary.exportIsGray(@keyOfMap)
        end

        def toGray()
            CallLibrary.exportToGray(@keyOfMap)
        end

        def toRgba()
            CallLibrary.exportToRgba(@keyOfMap)
        end

        def moveBounds(left, top, right, bottom, r, g, b, a)
            CallLibrary.exportMoveBounds(@keyOfMap, left, top, right, bottom, r, g, b, a)
        end
    end
end
