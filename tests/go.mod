module github.com/anosh26/custardfin/tests

go 1.26.0

replace (
	github.com/anosh26/custardfin/apps/api => ../apps/api
	github.com/anosh26/custardfin/apps/edge => ../apps/edge
	github.com/anosh26/custardfin/libs/contracts => ../libs/contracts
)
