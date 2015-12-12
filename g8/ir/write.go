package ir

import (
	"e8vm.io/e8vm/link8"
)

func writeBlock(f *link8.Func, b *Block) {
	for _, inst := range b.insts {
		f.AddInst(inst.inst)
		if inst.sym != nil {
			s := inst.sym
			f.AddLink(s.fill, &link8.PkgSym{s.pkg, s.sym})
		}
	}
}

func writeFunc(p *Pkg, f *Func) {
	lfunc := link8.NewFunc()

	for b := f.prologue; b != nil; b = b.next {
		writeBlock(lfunc, b)
	}

	p.lib.DefineFunc(f.name, lfunc)
}
