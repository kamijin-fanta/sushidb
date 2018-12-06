module github.com/kamijin-fanta/sushidb

require (
	github.com/0xAX/notificator v0.0.0-20171022182052-88d57ee9043b // indirect
	github.com/GeertJohan/go.incremental v0.0.0-20161212213043-1172aab96510 // indirect
	github.com/GeertJohan/go.rice v0.0.0-20170420135705-c02ca9a983da
	github.com/akavel/rsrc v0.0.0-20170831122431-f6a15ece2cfd // indirect
	github.com/codegangsta/envy v0.0.0-20141216192214-4b78388c8ce4 // indirect
	github.com/codegangsta/gin v0.0.0-20171026143024-cafe2ce98974 // indirect
	github.com/daaku/go.zipexe v0.0.0-20150329023125-a5fe2436ffcb // indirect
	github.com/gin-contrib/sse v0.0.0-20170109093832-22d885f9ecc7 // indirect
	github.com/gin-gonic/gin v1.3.0
	github.com/jessevdk/go-flags v1.4.0 // indirect
	github.com/joho/godotenv v1.3.0
	github.com/juju/errors v0.0.0-20181012004132-a4583d0a56ea // indirect
	github.com/kardianos/osext v0.0.0-20170510131534-ae77be60afb1 // indirect
	github.com/mattn/go-isatty v0.0.4 // indirect
	github.com/mattn/go-shellwords v1.0.3 // indirect
	github.com/pingcap/kvproto v0.0.0-20181105061835-1b5d69cd1d26
	github.com/pingcap/tidb v2.1.0-rc.2+incompatible
	github.com/pkg/errors v0.9.0 // indirect
	github.com/vmihailenco/msgpack v4.0.0+incompatible
	golang.org/x/net v0.0.0-20181029044818-c44066c5c816
	google.golang.org/grpc v1.16.0
	gopkg.in/go-playground/validator.v8 v8.18.2 // indirect
	gopkg.in/urfave/cli.v1 v1.20.0 // indirect
)

replace github.com/pkg/errors => github.com/pingcap/errors v0.9.0

replace github.com/pingcap/tidb => github.com/kamijin-fanta/tidb v0.0.0-20181206023524-be036b180cae
