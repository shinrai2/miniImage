require 'ffi'
require 'os'
require 'pathname'

$DIR_PATH = Pathname.new(File.dirname(__FILE__)).realpath
$SYM_IMAGE = 0x22
$SYM_FONT = 0x33
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
    attach_function :exportToGray, [$GOINT, :bool], $GOINT
    attach_function :exportToRgba, [$GOINT, :bool], $GOINT
    attach_function :exportMoveBounds, [$GOINT, $GOINT, $GOINT, $GOINT, $GOINT, :uchar, :uchar, :uchar, :uchar, :bool], $GOINT
    attach_function :exportNewBlank, [$GOINT, $GOINT, :uchar, :uchar, :uchar, :uchar], $GOINT
    attach_function :exportDrawString, [$GOINT, $GOINT, :double, $GOINT, $GOINT, :string, :uchar, :uchar, :uchar, :uchar], :void
    attach_function :exportConcat, [$GOINT, $GOINT, :bool], $GOINT
end

module BlankConstuctor
    def newBlank(width, height, color)
        return self.new(CallLibrary.exportNewBlank(width, height, *color.args))
    end
end

module BasicConstuctor
    def loadFrom(path)
        if self == MiniImage::Image then
            symbl = $SYM_IMAGE
        elsif self == MiniImage::Font then
            symbl = $SYM_FONT
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

            ObjectSpace.define_finalizer(self, Finalizer.finalize($SYM_IMAGE, key)) # avoid mem leak
        end

        def save(path)
            CallLibrary.exportSave(@keyOfMap, path)
        end

        def isGray()
            return CallLibrary.exportIsGray(@keyOfMap)
        end

        def toGray!()
            CallLibrary.exportToGray(@keyOfMap, true)
        end

        def toRgba!()
            CallLibrary.exportToRgba(@keyOfMap, true)
        end

        def moveBounds!(left, top, right, bottom, color)
            CallLibrary.exportMoveBounds(@keyOfMap, left, top, right, bottom, *color.args, true)
        end

        def toGray()
            n = CallLibrary.exportToGray(@keyOfMap, false)
            return MiniImage::Image.send(:new, n)
        end

        def toRgba()
            n = CallLibrary.exportToRgba(@keyOfMap, false)
            return MiniImage::Image.send(:new, n)
        end

        def moveBounds(left, top, right, bottom, color)
            n = CallLibrary.exportMoveBounds(@keyOfMap, left, top, right, bottom, *color.args, false)
            return MiniImage::Image.send(:new, n)
        end

        def self.concat(img1, img2, direct)
            n = CallLibrary.exportConcat(img1.keyOfMap, img2.keyOfMap, direct)
            return MiniImage::Image.send(:new, n)
        end

        attr_reader :keyOfMap
    end

    class Font
        extend BasicConstuctor

        def initialize(key)
            @keyOfMap = key

            ObjectSpace.define_finalizer(self, Finalizer.finalize($SYM_FONT, key)) # avoid mem leak
        end

        def drawString(img, fontSize, x, y, content, color)
            CallLibrary.exportDrawString(@keyOfMap, img.keyOfMap, fontSize, x, y, content, *color.args)
        end

        attr_reader :keyOfMap
    end

    class Color
        def initialize(r, g, b, a)
            @r, @g, @b, @a = r, g, b, a
        end

        def args()
            return [@r, @g, @b, @a]
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
