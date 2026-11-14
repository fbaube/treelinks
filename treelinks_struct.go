package treelinks

import(
)

// As laid out at
// https://www.perplexity.ai/session/69c2e8ed-f4d5-4338-871a-912b3e88cf52
// 
// It's best to keep Treelinks simple. If we simplify this 
// package's methods to apply only to Treelinks itself, so
// that we don't try to get clever by promoting its methods 
// to emedding super-structs, then we have two choices about 
// how to get a pointer to a Treelinks's embedding super-struct:
//  - an explicit owner reference in the Treelinks, or
//  - a sort-of-global  UUID-to-node lookup.
// 
// The latter seems dangerous, cos it does things globally; the
// former keeps an embedding super-struct tightly linked to its 
// own Treelinks but does not make the pair reachable from other 
// code. So Perplexity suggests the following code:

// Treelinks stores a set of links for a tree node embedded in a
// DB-persistable object. However... every DB-persistable object
// needs a UUID but may not need a set of links. And one may even
// have multiple sets of links. So a UUID for a DB-persistable
// object must be embedded below the level of links, so let's
// assign a UUID to every Nork. 
// 
// So there can be a Treelinks for child nodes, for incoming links,
// for any other object that itself has a Treelinks.
// Each target of a link has a UUID, and therefore so does every
// instance of Treelinks have a UUID. ANd therefore it should be
// stored by value, not as a pointer.
//
// When persisting
// to disk, the fields ( prnt prev next ) are sufficient 
// to reconstruct a tree that has all kid nodes in their
// correct order.
//
// There needs to be a distinguished "empty value" 
// to mark a non-value, and for a pointer this is 
// nil. (Be careful tho, because a nil instance of
// an interface is not nil!)
//
// When a tree structure is built that extends upward,
// with the root at the bottom, maybe some of these field
// names should be aliased to avoid confusion, since "up"
// and "down" are reversed, and so are "parent and "kid". 
// .
type Treelinks struct {
     // UUIDv7 includes the DB table, so that [OuterOwner] 
     // (the embedding super-struct) can be fetched by UUID
     // directly from the DB, rather than going thru some 
     // kind of in-memory router / map / UUID registry
     // that figures out which table.
//      UUIDv7 uuid.UUID
     outerOwner Treelinkser // not: ~Treelinks, which is a constraint 
     prnt *Treelinks   // left: optional: shown at top of L-H list
     prev *Treelinks   // previous sibling, or nil if is the first 
     next *Treelinks   // next sibling,  or nil if it is the last 
     kid0 *Treelinks   // first kid in linked list of kids 
     kidN *Treelinks   //  last kid in linked list of kids 
     usrs []*Treelinks // left: incoming: symlinks(FSO)/referrers(CMS)
//   tags []*Treelinks // topp: incoming: tags/facets
//   vers []*Treelinks // botm: versions: is a list, not a tree
     level int 
}

// ===========================
//   THESE NEXT THREE func's
//  BELONG IN (say) miscutils
// ===========================

// ZV returns the zero value for any T.
// https://www.perplexity.ai/search/717b5c69-194e-4599-88d6-2f00ee7a236a 
func ZV[T comparable]() T {
     // zed is not initialized to any specific valur,
     // so it gets initialized to the "zero value" !
     var zed T
     return zed
}

// IsZV reports whether arg is the zero value for T.
// https://www.perplexity.ai/search/717b5c69-194e-4599-88d6-2f00ee7a236a 
// NOTE: "==" works byte-for-byte for UUIDs !!
// See https://www.perplexity.ai/search/717b5c69-194e-4599-88d6-2f00ee7a236a 
// .
func IsZV[T comparable](arg T) bool {
    var zed T
    return arg == zed
}

// NotZV reports whether arg is the zero value for T.
// https://www.perplexity.ai/search/717b5c69-194e-4599-88d6-2f00ee7a236a 
// NOTE: "==" works byte-for-byte for UUIDs !!
// See https://www.perplexity.ai/search/717b5c69-194e-4599-88d6-2f00ee7a236a 
// .
func NotZV[T comparable](arg T) bool {
    var zed T
    return arg != zed
}

