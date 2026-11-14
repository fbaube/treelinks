package treelinks

import (
)

// Stringser is copied from
// [github.com/fbaube/stringutils.Stringser] to keep out dependencies.
type Stringser interface {
        Echo()  string
        Infos() string
        Debug() string
}

// OBS.? The interfaces have to be divided into two sets:
//  - One for Nork, that works with a basic tree.
//  - One for Cnty, that works for such a tree, but
//    with additional func's defined for a filesystem.
//  - Cnty can redefine-and-override Nork versions. 

// Treelinkser is satisfied by [*Treelinks] 
// and NOT by Treelinks. Therefore
//  - We never use a pointer (i.e. *Treelinkser), and
//  - Every method on a Treelinks has write access to
//    all the fields. 
//
// Func's are divided into three groups:
//  - Basic parent/child stuff
//     - Unlink()
//     - Links() TL
//     - UUID() UUID 
//     - Parent() TL 
//     - SetParent(p2 TL)
//     - OuterOwner() TL
//     - HasKids() bool 
//     - ValidateKids() error 
//  - Kid queries
//     - FirstKid() TL
//     - LastKid()  TL
//     - KidsAsSlice() []TL 
//  - Prev & Next peer (i.e. sibling node) 
//     - Prev() *TL
//     - Next() *TL
//     - SetPrev(newPrev *TL)
//     - SetNext(newNext *TL)
//  - Kid changes
//     - AddKids(rKids []T) T // returns p, the Parent
//     - AddKid (aKid   T) T // returns aKid
//     - SetFirstKid(p2 T)  // does what ?
//     - SetLastKid (p2 T) // does what ?
//     - func's to delete kid(s) ? 

// Treelinkser specifies the general case for any/all struct types.
// type Treelinkser interface { Links() *Treelinks }

// Links implements interface Treelinkser on struct type [Treelinks].
// It seems silly/obvious, but it will work out (I guess). 
func (p *Treelinks) Links() *Treelinks { return p }

type Treelinkser interface {

// ==========================
//  BASIC PARENT/CHILD STUFF
// ==========================
    Unlink()
    Links()  *Treelinks
//  UUID()    uuid.UUID
//  TableName()  string // three letters 
	SetParent(Treelinkser)
	Parent()  Treelinkser 
	HasKids() bool
	FirstKid()  Treelinkser
	LastKid()   Treelinkser
    NilOutFirstKid()
    NilOutLastKid()
	PrependAsFirstKid(Treelinkser) 
	 AppendAsLastKid (Treelinkser) 
	KidsAsSlice() []Treelinkser
    ValidateKids() error
// ======================
//    PREV & NEXT PEER
//  (i.e. SIBLING STUFF)
// ======================
	Prev() Treelinkser 
	Next() Treelinkser
    SetPrev(p2 Treelinkser)
    SetNext(p2 Treelinkser)

// =================
//  CHANGES TO KIDS     
// =================

// AddKid returns the just-added kid, 
// who is added at the end of (i.e.
// appended to) the list of Kids.
// AddKid(Treelinkser) Treelinkser

// AddKids returns the method target
// - the parent of all the kids 
// AddKids([]Treelinkser) Treelinkser
   ReplaceBy(Treelinkser) Treelinkser // returns Whom
   InsertBefore(Treelinkser)
   InsertAfter (Treelinkser)
// DELETE A KID ?
}

/*

type Norker interface {
     	Stringser
        Treelinkser 
	// =====================
	//  PATHS (Rel. & Abs.) 
	// =====================
	// RelPath is rel.filepath for a file/dir, and for a DOM
	// node, it is meaningless, unless it is a [RootNorker],
	// for which it is the rel.path to the containing document. 
	RelPath() string
	SetRelPath(string) 
	// AbsPath is abs.filepath for a file/dir, and for a DOM
	// node, the (abs.)path of the node w.r.t. the document
	// root, except for a [RootNorker], for which it is
	// the abs.path to the containing document. 
	AbsPath() string
	SetAbsPath(string)
	// ======================
	//  PATH CHARACTERISTICS
	// ======================
	// Level is zero-based (i.e. root Nork's is 0) 
	Level() int
	// SetLevel(int) // should be self-calculated by node 
	// Root should always return the root
	//  (at arena index 0, if one is used) 
	Root() Treelinkser // RootNork
	IsRoot() bool
	IsDir() bool
	IsDirlike() bool
	// ============
	//  MISCELLANY
	// ============
	LinePrefixString() string
	PrintTree(io.Writer) error
}

// Cntyer is satisfied by [*Cnty] NOT by Cnty.
type Cntyer interface {
	// =========================
	//  FILEPATHS (Rel. & Abs.) 
	// =========================
	// RelFP is rel.filepath for a file/dir.
	RelFP() string
	SetRelFP(string) 
	// AbsFP is abs.filepath for a file/dir.
	AbsFP() string
	SetAbsFP(string)
	// Kids stuff that only makes sense for Kids in a linked list. 
	SetFirstKid(Treelinkser)
	SetLastKid (Treelinkser)
	SetPrevPeer(Treelinkser)
	SetNextPeer(Treelinkser)
}

*/

