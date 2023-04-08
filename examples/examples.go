package examples

import _ "embed"

//go:embed function.em
var Function string

//go:embed main_and_helper.em
var MainAndHelper string
