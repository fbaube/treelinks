package treelinks

import (
)

// Stringser is copied from
// [github.com/fbaube/stringutils.Stringser] to
// keep out dependencies. Infos is spelled that
// way cos of a name collision with the stdlib. 
type Stringser interface {
        Echo()  string
        Infos() string
        Debug() string
}

// Links (part of interface Treelinkser) gets promoted
// in embedding structs, so that it can be called on 
// an "outer owning" struct that embeds Treelinks. 
func (p *Treelinks) LinksPtr() *Treelinks { return p }

// Treelinkser is implemented by [*Treelinks] NOT by Treelinks. Therefore
//  - We never use a pointer to it (i.e. *Treelinkser makes no sense), and
//  - Every method on a Treelinks has write access to all the fields. 
//
// When an interface method returns a Treelinkser, it can 
// return an outer-owning embedding struct (that is not nil).
// But for now (2026.11) we will return a *Treelinks. 

type Treelinkser interface {

// ==========================
//  BASIC, SIMPLE PARENT/KID
//   INFORMATION AND QUERIES
//    (getters and setters 
//    have no side effects) 
// ==========================
	   Unlink()
	   Parent() Treelinkser 
	    Level() int 
	SetParent(Treelinkser) 
	SetLevel (newL int) 
	  HasKids() bool 
	 FirstKid() Treelinkser 
	  LastKid() Treelinkser 
	   OuterOwner() Treelinkser  
	SetOuterOwner(newOO Treelinkser) 
    LinksPtr() *Treelinks // had omitted this to 
//  avoid type-casting to *Treelinks all the time 
// =======================
//   PREV & NEXT PEERS
//   (i.e. SIBLINGS in
//  LINKED LIST of KIDS)
//  (getters and setters 
//  have no side effects) 
// =======================
	Prev() Treelinkser 
	Next() Treelinkser
    SetPrev(newPrev Treelinkser)
    SetNext(newNext Treelinkser)
// =========================
//  ADVANCED KID OPERATIONS
// =========================
	KidsAsSlice() []Treelinkser
    ValidateKids() error
// ==================================
//    COMPLEX KID OPERATIONS that
//  affect multiple fields and might
//  need to be placed behind a Mutex
// ==================================
	PrependAsFirstKid(Treelinkser) 
	 AppendAsLastKid (Treelinkser) 
     AppendAllKids ([]Treelinkser) 
     ReplaceBy(newTL Treelinkser) 
      BeforeItInsert(Treelinkser) 
       AfterItInsert(Treelinkser) 
// DELETE A KID ?
}

