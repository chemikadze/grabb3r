GRABB3R
=======
Coding excercise solution export tool.

Usage
-----
```
go get github.com/chemikadze/cmd/grabb3r
export LEETCODE_USER=myuser
export LEETCODE_PASSWORD=MyPaSsWoRd
grabb3r
ls ./dest/file/
```

Internals
---------
`SolutionSource` abstracts exercise solution sites with listing and 
solution retrieval methods. Right now only LeetCode is provided. 

`SolutionDestination` provides interface for export destination.
At this moment, only mock (console) and file are supported.
