module github.com/niccolot/GoDex/commands

go 1.22.5

require (
    github.com/niccolot/GoDex/pokeapi v0.0.0
    github.com/niccolot/GoDex/types v0.0.0
)

replace github.com/niccolot/GoDex/pokeapi v0.0.0 => ../pokeapi
replace github.com/niccolot/GoDex/types v0.0.0 => ../types
