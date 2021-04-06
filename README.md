![Multi-Splay Tree](https://imgur.com/XKvuwDp.png)

# Multi-Splay Tree

This is an implementation of a multi-splay tree, written in Go.

## Multi-what tree, now?

To spare you the time of reading the paper just to have an undertanding of
what the data structure might be, this section explains the main idea behind it.

The multi-splay tree is the splay-tree counterpart of [Tango Trees][tango],
which are a special class of binary search trees which are proven to be
O(log log n)-competitive on both searches and updates.
It makes use of red-black trees to achieve this bound, mainly by breaking
the tree into many smaller trees, which are governed by another (implicit)
red-black tree, called the reference tree.

The paper for Tango Trees is available at [Erik Demaine's website][tango-paper],
and is worth a read (it's reasonably short, too).

[tango]: https://en.wikipedia.org/wiki/Tango_tree
[tango-paper]: https://erikdemaine.org/papers/Tango_SICOMP/paper.pdf

Both the tango tree and multi-splay trees are efforts on the direction of answering 
the [dynamic optimality conjecture][dynamic-optimality], introduced by
Sleator and Tarjan on the original splay-tree paper. Both achieve a very good bound
of O(log log n)-competitiveness. While it keeps this bound, the multi-splay tree 
improves on the tango tree by providing a better amortized cost for operations on
the tree (O(log n) instead of the O(log n log log n) achieved by tango trees). It 
is also an open problem whether the multi-splay tree is dinamically optimal, while
the tango tree is known not to be.

[dynamic-optimality]: https://en.wikipedia.org/wiki/Optimal_binary_search_tree#Dynamic_optimality

The general construction of the multi-splay tree is represented on the following
diagram, which shows both the reference tree (which is an implicit red-black tree) 
as wells as the actual multi-splay tree. The splay trees on the MST correspond to
preferred paths (paths of preferred children) on the reference tree:

![Diagram 1](https://imgur.com/sMABdIx.png)

## Is this practical?

No.

The constant factors associated with maintaining a multi-splay tree are _huge_, and
having them on a practical setting is not viable. It is a very interesting piece of
research, though, and very nice to play with.

## Why?

I originally wrote it in C for a college assignment (and eventually gave up on using
it because it timed out on the automated tests, unfortunatelly).

I rewrote it in Go to get used to the language.

## TODO List
- Implement deletions (the original assignment didn't require them 'w').
- Write some tests.
- Cleanup (?). I'm new to Go, since I mostly translated the code directly from C, 
  there might be some rough edges in it.

## References

[WDS06] C. C. Wang, J. Derryberry, and D. D. Sleator.
        O(log log n)-competitive dynamic binary search trees.
        SODA, pp. 374–383, 2006.
        Available at https://www.cs.cmu.edu/~chengwen/paper/MST.pdf
        Extended version available at https://pdfs.semanticscholar.org/8006/044b0a69b9d1828711ce909f3201f49c7b06.pdf.

[ST85] D. D. Sleator and R. E. Tarjan.
       Self-Adjusting Binary Search Trees.
       Journal of the Association for Computing Machinery, Vol. 32, No. 3,
       pp. 652-686, 1985. 
       Available at https://www.cs.cmu.edu/~sleator/papers/self-adjusting.pdf.
