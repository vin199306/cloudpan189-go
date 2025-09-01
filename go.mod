module github.com/tickstep/cloudpan189-go

go 1.19

require (
	github.com/GeertJohan/go.incremental v1.0.0
	github.com/json-iterator/go v1.1.10
	github.com/kardianos/osext v0.0.0-20170510131534-ae77be60afb1
	github.com/oleiade/lane v0.0.0-20160817071224-3053869314bb
	github.com/olekukonko/tablewriter v0.0.2-0.20190618033246-cc27d85e17ce
	github.com/peterh/liner v1.1.1-0.20190305032635-6f820f8f90ce
	github.com/tickstep/bolt v1.3.4
	github.com/tickstep/cloudpan189-api v0.1.0
	github.com/tickstep/library-go v0.1.2
	github.com/urfave/cli v1.21.1-0.20190817182405-23c83030263f
)

// 保持原版本 API 库，通过编译参数解决兼容性问题

require (
	github.com/boltdb/bolt v1.3.1 // indirect
	github.com/cpuguy83/go-md2man v1.0.10 // indirect
	github.com/denisbrodbeck/machineid v1.0.1 // indirect
	github.com/mattn/go-runewidth v0.0.9 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v0.0.0-20180701023420-4b7aa43c6742 // indirect
	github.com/russross/blackfriday v1.5.2 // indirect
	github.com/satori/go.uuid v1.2.0 // indirect
	golang.org/x/sys v0.0.0-20200202164722-d101bd2416d5 // indirect
)

//replace github.com/tickstep/bolt => /Users/tickstep/Documents/Workspace/go/projects/bolt
//replace github.com/tickstep/library-go => /Users/tickstep/Documents/Workspace/go/projects/library-go
//replace github.com/tickstep/cloudpan189-api => /Users/tickstep/Documents/Workspace/go/projects/cloudpan189-api
