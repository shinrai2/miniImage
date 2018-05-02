require 'ffi'
require 'os'

if OS.windows? == false then
    print "Don\'t support non Windows operating system."
    exit 0
elsif OS.bits == 64 then
    $GOINT = :long_long
else
    $GOINT = :int
end

module CallLibrary
    extend FFI::Library
    ffi_lib "libs/#{ OS.bits }/miniImage.so"
    attach_function :exportInitialize, [:string], $GOINT
    attach_function :exportSave, [$GOINT, :string], :void
    attach_function :exportRelease, [$GOINT], :void
    attach_function :exportIsGray, [$GOINT], :bool
    attach_function :exportToGray, [$GOINT], :void
    attach_function :exportToRgba, [$GOINT], :void
    attach_function :exportMoveBounds, [$GOINT, $GOINT, $GOINT, $GOINT, $GOINT, :uchar, :uchar, :uchar, :uchar], :void
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
