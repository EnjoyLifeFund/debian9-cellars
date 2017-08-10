/* Default linker script, for normal executables */
/* Example Linker Script for linking NS CR16 elf32 files.
   Copyright (C) 2014-2017 Free Software Foundation, Inc.
   Copying and distribution of this script, with or without modification,
   are permitted in any medium without royalty provided the copyright
   notice and this notice are preserved.  */
OUTPUT_FORMAT("elf32-cr16")
OUTPUT_ARCH(cr16)
ENTRY(_start)
/* Define memory regions.  */
MEMORY
{
        rom         : ORIGIN = 0x2,         LENGTH = 3M
        ram         : ORIGIN = 4M,          LENGTH = 10M
}
/*  Many sections come in three flavours.  There is the 'real' section,
    like ".data".  Then there are the per-procedure or per-variable
    sections, generated by -ffunction-sections and -fdata-sections in GCC,
    and useful for --gc-sections, which for a variable "foo" might be
    ".data.foo".  Then there are the linkonce sections, for which the linker
    eliminates duplicates, which are named like ".gnu.linkonce.d.foo".
    The exact correspondences are:
    Section	Linkonce section
    .text	.gnu.linkonce.t.foo
    .rdata	.gnu.linkonce.r.foo
    .data	.gnu.linkonce.d.foo
    .bss	.gnu.linkonce.b.foo
    .debug_info	.gnu.linkonce.wi.foo  */
SECTIONS
{
  .init :
  {
    __INIT_START = .;
    KEEP (*(.init))
    __INIT_END = .;
  } > rom
  .fini :
  {
    __FINI_START = .;
    KEEP (*(.fini))
    __FINI_END = .;
  } > rom
  .jcr :
  {
    KEEP (*(.jcr))
  } > rom
  .text :
  {
    __TEXT_START = .;
    *(.text) *(.text.*) *(.gnu.linkonce.t.*)
    __TEXT_END = .;
  } > rom
  .rdata :
  {
    __RDATA_START = .;
    *(.rdata_4) *(.rdata_2) *(.rdata_1) *(.rdata.*) *(.gnu.linkonce.r.*) *(.rodata*)
    __RDATA_END = .;
  } > rom
  .ctor ALIGN(4) :
  {
    __CTOR_START = .;
    /* The compiler uses crtbegin.o to find the start
       of the constructors, so we make sure it is
       first.  Because this is a wildcard, it
       doesn't matter if the user does not
       actually link against crtbegin.o; the
       linker won't look for a file to match a
       wildcard.  The wildcard also means that it
       doesn't matter which directory crtbegin.o
       is in.  */
    KEEP (*crtbegin*.o(.ctors))
    /* We don't want to include the .ctor section from
       the crtend.o file until after the sorted ctors.
       The .ctor section from the crtend file contains the
       end of ctors marker and it must be last */
    KEEP (*(EXCLUDE_FILE (*crtend*.o) .ctors))
    KEEP (*(SORT(.ctors.*)))
    KEEP (*(.ctors))
    __CTOR_END = .;
  } > rom
  .dtor ALIGN(4) :
  {
    __DTOR_START = .;
    KEEP (*crtbegin*.o(.dtors))
    KEEP (*(EXCLUDE_FILE (*crtend*.o) .dtors))
    KEEP (*(SORT(.dtors.*)))
    KEEP (*(.dtors))
    __DTOR_END = .;
  } > rom
  .data :
  {
    __DATA_START = .;
    *(.data_4) *(.data_2) *(.data_1) *(.data) *(.data.*) *(.gnu.linkonce.d.*)
    __DATA_END = .;
  } > ram AT > rom
  .bss (NOLOAD) :
  {
    __BSS_START = .;
    *(.bss_4) *(.bss_2) *(.bss_1) *(.bss) *(COMMON) *(.bss.*) *(.gnu.linkonce.b.*)
    __BSS_END = .;
  } > ram
/* You may change the sizes of the following sections to fit the actual
   size your program requires.
   The heap and stack are aligned to the bus width, as a speed optimization
   for accessing data located there.  */
  .heap (NOLOAD) :
  {
    . = ALIGN(4);
    __HEAP_START = .;
    . += 0x2000; __HEAP_MAX = .;
  } > ram
  .stack (NOLOAD) :
  {
    . = ALIGN(4);
    . += 0x6000;
    __STACK_START = .;
  } > ram
  .istack (NOLOAD) :
  {
    . = ALIGN(4);
    . += 0x100;
    __ISTACK_START = .;
  } > ram
  .comment        0 : { *(.comment) }
  /* DWARF debug sections.
     Symbols in the DWARF debugging sections are relative to the beginning
     of the section so we begin them at 0.  */
  /* DWARF 1 */
  .debug          0 : { *(.debug) }
  .line           0 : { *(.line) }
  /* GNU DWARF 1 extensions */
  .debug_srcinfo  0 : { *(.debug_srcinfo) }
  .debug_sfnames  0 : { *(.debug_sfnames) }
  /* DWARF 1.1 and DWARF 2 */
  .debug_aranges  0 : { *(.debug_aranges) }
  .debug_pubnames 0 : { *(.debug_pubnames) }
  /* DWARF 2 */
  .debug_info     0 : { *(.debug_info .gnu.linkonce.wi.*) }
  .debug_abbrev   0 : { *(.debug_abbrev) }
  .debug_line     0 : { *(.debug_line .debug_line.* .debug_line_end ) }
  .debug_frame    0 : { *(.debug_frame) }
  .debug_str      0 : { *(.debug_str) }
  .debug_loc      0 : { *(.debug_loc) }
  .debug_macinfo  0 : { *(.debug_macinfo) }
  /* SGI/MIPS DWARF 2 extensions */
  .debug_weaknames 0 : { *(.debug_weaknames) }
  .debug_funcnames 0 : { *(.debug_funcnames) }
  .debug_typenames 0 : { *(.debug_typenames) }
  .debug_varnames  0 : { *(.debug_varnames) }
  /* DWARF 3 */
  .debug_pubtypes 0 : { *(.debug_pubtypes) }
  .debug_ranges   0 : { *(.debug_ranges) }
  /* DWARF Extension.  */
  .debug_macro    0 : { *(.debug_macro) }
  .debug_addr     0 : { *(.debug_addr) }
}
__DATA_IMAGE_START = LOADADDR(.data);