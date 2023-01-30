module examples

go 1.19

require (
	github.com/fogleman/gg v1.3.0
	github.com/quickqr/gqr v0.2.0-beta
)

require (
	github.com/golang/freetype v0.0.0-20170609003504-e2365dfdc4a0 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	golang.org/x/image v0.3.0 // indirect
)

replace (
	github.com/quickqr/gqr => ../
	github.com/quickqr/gqr/export/image => ../export/image
)
