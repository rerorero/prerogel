module github.com/rerorero/prerogel/examples/maximum

go 1.12

replace github.com/rerorero/prerogel => ../../

// TODO: remove
replace github.com/AsynkronIT/protoactor-go => ../../../../AsynkronIT/protoactor-go

require (
	github.com/AsynkronIT/protoactor-go v0.0.0-20190619151400-46c416c5ae7f
	github.com/gogo/protobuf v1.2.1
	github.com/pkg/errors v0.8.1
	github.com/rerorero/prerogel v0.0.0-20190621071045-a79bc9d3320f
)
