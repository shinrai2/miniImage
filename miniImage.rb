require 'ffi'
require 'os'
require 'pathname'

$DIR_PATH = Pathname.new(File.dirname(__FILE__)).realpath
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
    ffi_lib "#{ $DIR_PATH }/libs/#{ OS.bits }/miniImage.so" # 使用绝对路径，否则会报错
    attach_function :exportFromFile, [:uchar, :string], $GOINT
    attach_function :exportSave, [$GOINT, :string], :void
    attach_function :exportRelease, [:uchar, $GOINT], :void
    attach_function :exportIsGray, [$GOINT], :bool
    attach_function :exportToGray, [$GOINT], :void
    attach_function :exportToRgba, [$GOINT], :void
    attach_function :exportMoveBounds, [$GOINT, $GOINT, $GOINT, $GOINT, $GOINT, :uchar, :uchar, :uchar, :uchar], :void
    attach_function :exportNewBlank, [$GOINT, $GOINT, :uchar, :uchar, :uchar, :uchar], $GOINT
end

module BlankConstuctor
    def newBlank(width, height, r=255, g=255, b=255, a=255)
        return self.new(CallLibrary.exportNewBlank(width, height, r, g, b, a))
    end
end

module BasicConstuctor
    def loadFrom(path)
        if self == MiniImage::Image then
            symbl = 0x22
        elsif self == MiniImage::Font then
            symbl = 0x33
        end
        return self.new(CallLibrary.exportFromFile(symbl, path))
    end

    protected
    def new(key)
        super
    end
end

module MiniImage
    class Image
        extend BasicConstuctor
        extend BlankConstuctor

        def initialize(key)
            @keyOfMap = key

            ObjectSpace.define_finalizer(self, Finalizer.finalize(0x22, key)) # avoid mem leak
        end

        def save(path)
            CallLibrary.exportSave(@keyOfMap, path)
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

    class Font
        extend BasicConstuctor

        def initialize(key)
            @keyOfMap = key

            ObjectSpace.define_finalizer(self, Finalizer.finalize(0x33, key)) # avoid mem leak
        end
    end
    private
    class Finalizer
        def self.finalize(sym, key)
            proc {
                CallLibrary.exportRelease(sym, key)
                printf("GC 0x%x: %d\n", sym, key)
            }
        end
    end
end
