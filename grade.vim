if exists('b:current_syntax')
 finish
endif

syntax region gradeParts start="^\s*>.*$" end="^\s*>" contains=gradePartWeight,gradePartData,gradePartName,integer,float,comment,assignment
syntax match gradePartName "\s*>.*$" contained 
syntax keyword gradePartWeight weight contained nextgroup=assignment
syntax keyword gradePartData data contained nextgroup=assignment

" syntax region meta start="^\s*\~" end="^\s*>" contains=comment,assignment,integer,float,metaVariable,metaHeader
" syntax match metaVariable "\k+" contained nextgroup=assignment
" syntax match metaHeader "^\s*\~ Meta" contained 
"
syntax match metaHeader "^\s*\~ Meta" 

syntax match assignment "\s*=\s*" 
syntax match comment "^#.*$"
syntax match integer "\d\+" 
syntax match float "\d\+\.?\d*" 

hi def link gradePartName Label
hi def link gradePartWeight Keyword
hi def link gradePartData Keyword

hi def link metaHeader Structure
" hi def link metaVariable Keyword

" hi def link assignment Statement
hi def link float Float
hi def link integer Number
" hi def link weight Identifier
" hi def link meta Define
hi def link comment Comment

let b:current_syntax = 'grade'
