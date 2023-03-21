package firespace

#IsBool:{
    in: _

    _valid: bool
    _valid: *in | false

    out: bool
    out: _valid
}

#IsStringArray:{
    in: _
    _valid: [...string] 
    _valid: *in | null

    _valid2: bool
   
    if _valid != null{
        _valid2: true
    }
   

    out: bool
    out: _valid2

}

notBool: "true"

_isBool: #IsBool & {in: notBool}
out: _isBool.out

