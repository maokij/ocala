package core

var SpecialMap = map[*Keyword]SyntaxFn{
	Intern("__PC__"):          SyntaxFn((*Compiler).sCurloc),
	Intern("__ORG__"):         SyntaxFn((*Compiler).sCurorg),
	Intern("loaded-as-main?"): SyntaxFn((*Compiler).sLoadedAsMain),
	Intern("<reserved>"):      SyntaxFn((*Compiler).sReserved),
	Intern("__FILE__"):        SyntaxFn((*Compiler).sFilename),
}

var SyntaxMap = map[*Keyword]SyntaxFn{
	Intern("include"):             SyntaxFn((*Compiler).sInclude),
	Intern("load-file"):           SyntaxFn((*Compiler).sLoadFile),
	Intern("arch"):                SyntaxFn((*Compiler).sArch),
	Intern("align"):               SyntaxFn((*Compiler).sAlign),
	Intern("#.label"):             SyntaxFn((*Compiler).sLabel),
	Intern("#.tpl"):               SyntaxFn((*Compiler).sTpl),
	Intern("#.prog"):              SyntaxFn((*Compiler).sProg),
	Intern("#.block"):             SyntaxFn((*Compiler).sBlock),
	Intern("do"):                  SyntaxFn((*Compiler).sDo),
	Intern("apply"):               SyntaxFn((*Compiler).sApply),
	Intern("loop"):                SyntaxFn((*Compiler).sLoop),
	Intern("if"):                  SyntaxFn((*Compiler).sIf),
	Intern("case"):                SyntaxFn((*Compiler).sCase),
	Intern("when"):                SyntaxFn((*Compiler).sWhen),
	Intern("alias"):               SyntaxFn((*Compiler).sAlias),
	Intern("#.module"):            SyntaxFn((*Compiler).sModule),
	Intern("section"):             SyntaxFn((*Compiler).sSection),
	Intern("link"):                SyntaxFn((*Compiler).sLink),
	Intern("flat!"):               SyntaxFn((*Compiler).sFlatMode),
	Intern("pragma"):              SyntaxFn((*Compiler).sPragma),
	Intern("#.macro"):             SyntaxFn((*Compiler).sMacro),
	Intern("#.proc"):              SyntaxFn((*Compiler).sProc),
	Intern("#.callproc"):          SyntaxFn((*Compiler).sCallproc),
	Intern("fallthrough"):         SyntaxFn((*Compiler).sFallthrough),
	Intern("tco"):                 SyntaxFn((*Compiler).sTco),
	Intern("#.const"):             SyntaxFn((*Compiler).sConst),
	Intern("#.data"):              SyntaxFn((*Compiler).sData),
	Intern("#.datalist"):          SyntaxFn((*Compiler).sDataList),
	Intern("#.structdata"):        SyntaxFn((*Compiler).sStructData),
	Intern("#.struct"):            SyntaxFn((*Compiler).sStruct),
	Intern("#.array"):             SyntaxFn((*Compiler).sArray),
	Intern("#.BYTE"):              SyntaxFn((*Compiler).sByte),
	Intern("#.REP"):               SyntaxFn((*Compiler).sRep),
	Intern("#.INVALID"):           SyntaxFn((*Compiler).sInvalid),
	Intern("#.mem"):               SyntaxFn((*Compiler).sMem),
	Intern("#.valueof"):           SyntaxFn((*Compiler).sValueOf),
	Intern("#.field"):             SyntaxFn((*Compiler).sFieldOffset),
	Intern("#.with"):              SyntaxFn((*Compiler).sWith),
	Intern("compile-error"):       SyntaxFn((*Compiler).sCompileError),
	Intern("assert"):              SyntaxFn((*Compiler).sAssert),
	Intern("import"):              SyntaxFn((*Compiler).sImport),
	Intern("expand-loop"):         SyntaxFn((*Compiler).sExpandLoop),
	Intern("*patch*"):             SyntaxFn((*Compiler).sPatch),
	Intern("make-counter"):        SyntaxFn((*Compiler).sMakeCounter),
	Intern("debug-inspect"):       SyntaxFn((*Compiler).sDebugInspect),
	Intern("quote"):               SyntaxFn((*Compiler).sQuote),
	Intern("#.exprdata"):          SyntaxFn((*Compiler).sExprdata),
	Intern("#.invalid-expansion"): SyntaxFn((*Compiler).sInvalidExpansion),
	Intern("&&"):                  SyntaxFn((*Compiler).sAnd),
	Intern("||"):                  SyntaxFn((*Compiler).sOr),
	Intern("sizeof"):              SyntaxFn((*Compiler).sSizeof),
	Intern("nameof"):              SyntaxFn((*Compiler).sNameof),
	Intern("nametypeof"):          SyntaxFn((*Compiler).sNametypeof),
	Intern("exprtypeof"):          SyntaxFn((*Compiler).sExprtypeOf),
	Intern("defined?"):            SyntaxFn((*Compiler).sDefinedp),
}

var FunMap = map[*Keyword]SyntaxFn{
	Intern("#.make-id"):  SyntaxFn((*Compiler).fMakeId),
	Intern("*"):          SyntaxFn((*Compiler).fMul),
	Intern("/"):          SyntaxFn((*Compiler).fDiv),
	Intern("%"):          SyntaxFn((*Compiler).fMod),
	Intern("+"):          SyntaxFn((*Compiler).fAdd),
	Intern("-"):          SyntaxFn((*Compiler).fSub),
	Intern("<<"):         SyntaxFn((*Compiler).fLsl),
	Intern(">>"):         SyntaxFn((*Compiler).fAsr),
	Intern(">>>"):        SyntaxFn((*Compiler).fLsr),
	Intern("<"):          SyntaxFn((*Compiler).fLt),
	Intern("<="):         SyntaxFn((*Compiler).fLe),
	Intern(">"):          SyntaxFn((*Compiler).fGt),
	Intern(">="):         SyntaxFn((*Compiler).fGe),
	Intern("=="):         SyntaxFn((*Compiler).fEql),
	Intern("!="):         SyntaxFn((*Compiler).fNotEql),
	Intern("&"):          SyntaxFn((*Compiler).fAnd),
	Intern("|"):          SyntaxFn((*Compiler).fOr),
	Intern("^"):          SyntaxFn((*Compiler).fXor),
	Intern("~"):          SyntaxFn((*Compiler).fNot),
	Intern("!"):          SyntaxFn((*Compiler).fLogicalNot),
	Intern("byte"):       SyntaxFn((*Compiler).fByte),
	Intern("word"):       SyntaxFn((*Compiler).fWord),
	Intern("lobyte"):     SyntaxFn((*Compiler).fLobyte),
	Intern("hibyte"):     SyntaxFn((*Compiler).fHibyte),
	Intern("asword"):     SyntaxFn((*Compiler).fAsWord),
	Intern("unuse?"):     SyntaxFn((*Compiler).fUnusep),
	Intern("use?"):       SyntaxFn((*Compiler).fUsep),
	Intern("formtypeof"): SyntaxFn((*Compiler).fFormtypeof),
}
