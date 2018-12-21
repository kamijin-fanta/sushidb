module github.com/kamijin-fanta/sushidb

require (
	github.com/GeertJohan/go.rice v0.0.0-20170420135705-c02ca9a983da
	github.com/daaku/go.zipexe v0.0.0-20150329023125-a5fe2436ffcb // indirect
	github.com/gin-contrib/sse v0.0.0-20170109093832-22d885f9ecc7 // indirect
	github.com/gin-gonic/gin v1.3.0
	github.com/joho/godotenv v1.3.0
	github.com/json-iterator/go v1.1.5 // indirect
	github.com/kardianos/osext v0.0.0-20170510131534-ae77be60afb1 // indirect
	github.com/mattn/go-isatty v0.0.4 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.1 // indirect
	github.com/oliveagle/jsonpath v0.0.0-20180606110733-2e52cf6e6852
	github.com/pingcap/pd v2.1.0-rc.4+incompatible
	github.com/pingcap/tidb v2.1.0-rc.2+incompatible
	github.com/pkg/errors v0.9.0 // indirect
	github.com/stretchr/testify v1.2.2
	github.com/vmihailenco/msgpack v4.0.0+incompatible
	gopkg.in/go-playground/assert.v1 v1.2.1 // indirect
	gopkg.in/go-playground/validator.v8 v8.18.2 // indirect
	gopkg.in/stretchr/testify.v1 v1.2.2
	k8s.io/client-go v10.0.0+incompatible
)

replace github.com/pkg/errors => github.com/pingcap/errors v0.9.0

replace github.com/pingcap/tidb => github.com/kamijin-fanta/tidb v0.0.0-20181206023524-be036b180cae
