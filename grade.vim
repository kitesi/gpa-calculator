" this is pretty bad, but vim regex is annoying to deal with and this works
" fine for now
" TODO: make numbers and floats bounded, as in they can't be part of a word

if exists('b:current_syntax')
 finish
endif

syntax region gradeParts start="^\s*>.*$" end="^\s*>" contains=gradePartWeight,gradePartData,gradePartName,integer,float,comment,assignment
syntax match gradePartName "\s*>.*$" contained 
syntax keyword gradePartWeight weight contained nextgroup=assignment
syntax keyword gradePartData data contained nextgroup=assignment

syntax match MetaHeader "^\s*\~ Meta" 

syntax match Assignment "\s*=\s*" 
syntax match Comment "#.*$"
syntax match Integer "\d\+" 
syntax match Float "\d\+\.?\d*" 

hi def link gradePartName Label
hi def link gradePartWeight Keyword
hi def link gradePartData Keyword
hi def link MetaHeader Structure
hi def link Float Float
hi def link Integer Number
hi def link Comment Comment

let b:current_syntax = 'grade'
