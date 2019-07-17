# Generating tests.pb.go

```
# Note: Change /usr/local/google/home/deklerk/workspace/googleapis to wherever
# you've installed https://github.com/googleapis/googleapis.
# Note: Run whilst cd-ed in this directory.
protoc --go_out=. -I /usr/local/google/home/deklerk/workspace/googleapis -I . *.proto
```