// https://www.perplexity.ai/session/69c2e8ed-f4d5-4338-871a-912b3e88cf52

package treelinks

import (
	"fmt"
	"os"
)

// Methods should accept and return [Treelinkser] wherever
// possible; only actual interface implementations take as
// the method receiver a [*Treelinks].

// On input (reading from DB) we use UUIDs to link *Nork's
// into a tree, and on output (writing to the DB) we write
// out the UUIDs.

// We need to persist the order of Kid nodes to the DB, so 
//  - keep using fields Prev and Next, and
//  - make sure they get persisted to the DB, and
//  - when storing them in memory, continue to use a
//    linked list instead of storing kids in a slice.

// We'd like to use a slice for Kids. But it doesn't
// preserve prev/next info when persisted to the DB. 
// To persist prev/next info, we need a unique ID 
// for each node, but... if we use an ID assigned 
// by SQLite, or the node's abs/rel path - these 
// are error-prone. Instead we now use a UUID, so 
// it makes sense to also use the UUID for prev/next. 
// Therefore when we persist a node to the DB:
//  - We MUST persist the UUID for the parent node 
//  - We MUST for the prev & next nodes too 
//  - We MAY  for the first & last kids (kid0 & kidN)
//  - A tree may be reconstructed using only the first
//    three fields (parent,prev,next), but using five 
//    provides redundancy for error checking: when we
//    read a tree back in from the DB, we check paths
//    to match kids to parents independently of UUIDs
//
// Another problem with using slices was that if for
// example we want prev or next peer/sibling, we must
// somehow use the input ptr to access the containing
// slice... but if parent is nil, there is no access. 
//
// So it makes sense to match parents and kids (and
// prev and next) using UUIDs. So when a tree is redd
// (back) in from the DB, it should be straightforward
// to validate each set of (parent + kid(s)).
//
// 1-to-1 links can be stored dirctly in persisted
// items rather than in a línk table, as required
// for 1-to-N or N-to-1. Links for parent, prev,
// next, firstKid, and lastKid are all 1-to-1.
// Parent-to-kids is 1-to-N but this link is
// not strictly necessary, as discussed above.
//
// Therefore a persisted Nork has three or five fields:
// Prnt, Prev, Next ; Kid0, KidN.

// ==========================
//  BASIC PARENT/CHILD STUFF
// ==========================

func (p *Treelinks) Unlink() {
//    p.prnt, p.prev, p,next, p.kid0, p,kidN =
//        ZV(T), ZV(T), ZV(T), ZV(T), ZV(T) }
//    var theZV T = ZV[T]()
      p.SetParent(nil)
      p.SetPrev(nil)
      p.SetNext(nil)
      p.NilOutFirstKid()
      p.NilOutLastKid ()
}

func (p *Treelinks)   Parent() Treelinkser { return p.prnt }
func (p *Treelinks) FirstKid() Treelinkser { return p.kid0 }
func (p *Treelinks)  LastKid() Treelinkser { return p.kidN }
func (p *Treelinks) NilOutFirstKid() { p.kid0 = nil }
func (p *Treelinks) NilOutLastKid()  { p.kidN = nil }
func (p *Treelinks) SetParent(newP  Treelinkser) { p.prnt = newP.Links() }
func (p *Treelinks)    OuterOwner()   Treelinkser  { return p.outerOwner } 
func (p *Treelinks) SetOuterOwner(newOO Treelinkser) { p.outerOwner = newOO } 

func (p *Treelinks) HasKids() bool {
     if NotZV(p.kid0) && NotZV(p.kidN) { return true }
     if  IsZV(p.kid0) &&  IsZV(p.kidN) { return false }
     panic("IsZV:FirstKid != IsZV:LastKid")
}

// PrependAsFirstKid prepends a Kid to the linked list of Kids.
func (p *Treelinks) PrependAsFirstKid(pNew0 Treelinkser) {
      // This all should be placed behind a lock.
      if !p.HasKids() {
          addAsNewOnlyKid(p, pNew0.(*Treelinks) )
        } else {
          var pOld0 = p.kid0
          // Parent <=> new Kid0
          pNew0.SetParent(p)
          p.kid0 = pNew0.(*Treelinks)
          // newKid0 <=> oldKid0
          pNew0.SetPrev(nil)
          pNew0.SetNext(pOld0)
          pOld0.SetPrev(pNew0)
      }
}

// AppendAsLastKid appends a Kid to the linked list of Kids.
func (p *Treelinks) AppendAsLastKid(pNewN Treelinkser) {
      // This all should be placed behind a lock.
      if !p.HasKids() {
          addAsNewOnlyKid(p, pNewN.(*Treelinks))
        } else {
          var pOldN = p.kidN
          // Parent <=> new KidN
          pNewN.SetParent(p)
          p.kidN = pNewN.(*Treelinks)
          // newKidN <=> oldKidN
          pNewN.SetNext(nil)
          pNewN.SetPrev(pOldN)
          pOldN.SetNext(pNewN)
      }
}

func addAsNewOnlyKid(pParent, pKid *Treelinks) {
     pKid.prnt = pParent
     pParent.kid0, pParent.kidN = pKid, pKid
}

func (p *Treelinks) ValidateKids() error {
     // FIXME 
     return nil 
}

// =============
//  KID QUERIES
// =============

func (p *Treelinks) KidsAsSlice() []Treelinkser {
     var slofT = make([]Treelinkser, 0)
     if !p.HasKids() { return slofT }
     var nexT *Treelinks
     nexT = p.kid0 // MAYBE have to get its Treelinks ? 
     for NotZV(nexT) {
         slofT = append(slofT, nexT)
         nexT = nexT.Links().Next().(*Treelinks)
         }
     return slofT
}

// ======================
//    PREV & NEXT PEER
//  (i.e. SIBLING STUFF) 
// ======================

func (p *Treelinks) Prev() Treelinkser { return p.prev }
func (p *Treelinks) Next() Treelinkser { return p.next }

func (p *Treelinks) SetPrev(newP Treelinkser) { p.prev = newP.Links() }
func (p *Treelinks) SetNext(newN Treelinkser) { p.next = newN.Links() }

// =================
//  CHANGES TO KIDS 
// =================

// AddKid adds the supplied node as the last kid, and returns
// it (i.e. the new last kid), now linked into the tree.
func (p *Treelinks) AddKid(aKid Treelinkser) Treelinkser { // returns aKid
	// Set the level & parent now
//	aKid.SetLevel(p.Level() + 1)
	aKid.SetParent(p)
	if !p.HasKids() { 
		// ---------------------------------------
		//  No Kids yet: New Kid will be ONLY Kid
		// ---------------------------------------
		p.kid0, p.kidN = aKid.Links(), aKid.Links()
        aKid.Links().SetPrev(nil)
        aKid.Links().SetNext(nil)
	} else { 
		// ---------------------
		//  Add as new last kid
		// ---------------------
		// FIXME var last = p.kidN()
	}
	return aKid
}

// AddKids adds the supplied nodes as kids, after 
// any pre-existing kids, and returns the parent.
// The supplied nodes may NOT be already linked
// to each other as a linked list. 
func (p *Treelinks) AddKids(rKids []Treelinkser) Treelinkser { // returns p 
	// fmt.Printf("nord: ptrs? aKid<%T> p<%T> \n", aKid, p)
	for _, aKid := range rKids {	
	    if aKid.Prev() != nil || aKid.Next() != nil {
		fmt.Fprintf(os.Stdout, "FATAL in AddKids: Tag<< %+v >> kid<< %+v >>\n", p, aKid)
		panic("AddKids(K) can't cos K has siblings")
		}
	    if aKid.Parent() != nil && aKid.Parent() != p {
		fmt.Fprintf(os.Stdout, "FATAL in AddKids: Tag<< %+v >> kid<< %+v >>\n", p, aKid)
		panic("E.AddKids(K) can't cos K has non-P parent")
		}
	    // All clear! Go ahead and add the kid.
	    _ = p.AddKid(aKid)
	}
	// FIXME JEEZ
	// For each new, set parent, then append whole mess to kids
	/*
	var FK = p.firstKid
	var LK = p.lastKid
	// Set the level now
	aKid.setLevel(p.Level() + 1)
	// Is the new kid an only kid ?
	if FK == nil && LK == nil {
		p.firstKid, p.lastKid = aKid, aKid
		aKid.SetParent(p)
		aKid.SetPrev(nil)
		aKid.SetNext(nil)
		return aKid
	}
	if !(FK != nil && LK != nil) {
		panic("BAD KID LINKS")
	}
	// So, replace the last kid
	if LK != nil {
		if LK.Parent() != p {
			fmt.Fprintf(os.Stdout, "FATAL in AddKid: E<< %+v >> K<< %+v >>\n", p, aKid)
			panic("E.AddKid: E's last kid dusnt know E")
		}
		if LK.Next() != nil {
			fmt.Fprintf(os.Stdout, "FATAL in AddKid: E<< %+v >> K<< %+v >>\n", p, aKid)
			panic("E.AddKid: E's last kid has a next kid")
		}
		LK.SetNext(aKid) // LK.nextKid = aKid
		aKid.SetPrev(LK) // aKid.prevKid = LK
		p.lastKid = aKid
		aKid.SetParent(p)
		return aKid
	}

	}
	fmt.Fprintf(os.Stdout, "FATAL in AddKids: " +
		"E<< %+v >> K<< %+v >>\n", p, rKids)
	panic("AddKids: Chaos!")
	*/
	return p
}

func (p *Treelinks) InsertBefore(p2 Treelinkser) { }
func (p *Treelinks) InsertAfter (p2 Treelinkser) { }
func (p *Treelinks) ReplaceBy(p2 Treelinkser) Treelinkser { return Treelinkser(p) }
