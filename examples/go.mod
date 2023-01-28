module examples

go 1.19

require (
	github.com/quickqr/gqr v0.1.0-beta
	github.com/quickqr/gqr/export/image v0.0.1-beta
)

require (
	github.com/fogleman/gg v1.3.0 // indirect
	github.com/golang/freetype v0.0.0-20170609003504-e2365dfdc4a0 // indirect
	github.com/pkg/errors v0.9.1 // indirect
)

replace (
	github.com/quickqr/gqr => ../
	github.com/quickqr/gqr/export/image => ../export/image
)
