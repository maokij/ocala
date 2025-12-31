" Plug '<path-to-ocala>/share/ocala/misc'

augroup ocala_plugin
  au!
  au BufNewFile,BufRead *.oc setlocal filetype=ocala
  au FileType ocala call s:ocala_setup()
augroup END

function! s:ocala_setup()
  if !exists("b:did_ftplugin")
    let b:did_ftplugin = 1

    setlocal suffixesadd+=.oc
    setlocal comments=s1:/*,mb:*,ex:*/,://
    setlocal commentstring=//%s
  endif

  if !exists("b:did_indent")
    let b:did_indent = 1
    setlocal cindent
    let b:undo_indent = "setl cin<"
  endif

  if exists("g:syntax_on") && !exists("b:current_syntax")
    syn keyword ocalaBuiltin *patch* <reserved> alias align apply arch
    syn keyword ocalaBuiltin assert break break-if case compile-error
    syn keyword ocalaBuiltin const continue continue-if data debug-inspect
    syn keyword ocalaBuiltin defined? do else expand-loop exprtypeof
    syn keyword ocalaBuiltin fallthrough flat! goto goto-if if import
    syn keyword ocalaBuiltin incbin include link load-file loop macro
    syn keyword ocalaBuiltin make-counter module nameof nametypeof
    syn keyword ocalaBuiltin never-return once optimize pragma proc quote
    syn keyword ocalaBuiltin recur redo redo-if return return-if section
    syn keyword ocalaBuiltin sizeof struct tco warn when
    syn keyword ocalaOperand A B C D E H L AF BC DE HL AF-
    syn keyword ocalaOperand IX IY IXH IXL IYH IYL
    syn keyword ocalaOperand X Y S P SP PC __PC__ __PROC__
    syn keyword ocalaOperand NZ? NE? not-zero? !=?
    syn keyword ocalaOperand Z? EQ? zero? ==?
    syn keyword ocalaOperand NC? CC? not-carry? borrow? >=?
    syn keyword ocalaOperand C? CS? carry? not-borrow? <?
    syn keyword ocalaOperand PO? VC? odd? not-over?
    syn keyword ocalaOperand PE? VS? even? over?
    syn keyword ocalaOperand P PL? plus?
    syn keyword ocalaOperand M MI? minus?
    syn match   ocalaLineComment "\/\/.*"
    syn region  ocalaComment     start="/\*"  end="\*/"
    syn region  ocalaString      start=+"+ skip=+\\\\\|\\"+ end=+"+
    syn region  ocalaCharacter   start=+'+ skip=+\\\\\|\\'+ end=+'+

    hi def link ocalaLineComment Comment
    hi def link ocalaComment Comment
    hi def link ocalaString String
    hi def link ocalaCharacter Character
    hi def link ocalaBuiltin Statement
    hi def link ocalaOperand Constant
  endif
endfunction
