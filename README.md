# Mocking Utility

This is a (profoundly) alpha version of some software I have written to
generate mocks. I know there are a few existing solutions that have come out of
FAANG, but I don't massively like magic when it comes to mocking.

This utility generates stubbed out implementations of interfaces, which you can
then override as you see fit.

General usage is as follows;

```
go run main.go -filename <source.go> > mocks.go
```

Don't use this. It works well enough I think but like I said, it's pretty
solidly alpha stage right now.
