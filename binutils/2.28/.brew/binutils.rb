class Binutils < Formula
  desc "FSF/GNU ld, ar, readelf, etc. for native development"
  homepage "https://www.gnu.org/software/binutils/binutils.html"
  url "https://ftp.gnu.org/gnu/binutils/binutils-2.28.tar.gz"
  mirror "https://ftpmirror.gnu.org/binutils/binutils-2.28.tar.gz"
  sha256 "cd717966fc761d840d451dbd58d44e1e5b92949d2073d75b73fccb476d772fcf"

  # No --default-names option as it interferes with Homebrew builds.
  option "with-default-names", "Do not prepend 'g' to the binary" if OS.linux?
  option "without-gold", "Do not build the gold linker" if OS.linux?

  depends_on "zlib" => :recommended unless OS.mac?

  def install
    system "./configure", "--disable-debug",
                          "--disable-dependency-tracking",
                          ("--program-prefix=g" if build.without? "default-names"),
                          ("--with-sysroot=/" if OS.linux?),
                          "--enable-deterministic-archives",
                          "--prefix=#{prefix}",
                          "--infodir=#{info}",
                          "--mandir=#{man}",
                          "--disable-werror",
                          "--enable-interwork",
                          "--enable-multilib",
                          "--enable-64-bit-bfd",
                          ("--enable-gold" if build.with? "gold"),
                          ("--enable-plugins" if OS.linux?),
                          "--enable-targets=all"
    system "make"
    system "make", "install"
    bin.install_symlink "ld.gold" => "gold" if build.with? "gold"

    # Reduce the size of the bottle.
    system "strip", *Dir[bin/"*", lib/"*.a"] unless OS.mac?
  end

  test do
    # Better to check with?("default-names"), but that doesn't work.
    nm = build.with?("default-names") ? "nm" : "gnm"
    assert_match "main", shell_output("#{bin}/#{nm} #{bin}/#{nm}")
  end
end
