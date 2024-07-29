module github.com/niccolot/GoDex

go 1.22.5

require (
    github.com/niccolot/GoDex/types v0.0.0
    github.com/niccolot/GoDex/commands v0.0.0
    github.com/niccolot/GoDex/utils v0.0.0
)

replace (
    github.com/niccolot/GoDex/types v0.0.0 => ./types
    github.com/niccolot/GoDex/commands v0.0.0 => ./commands
    github.com/niccolot/GoDex/utils v0.0.0 => ./utils
)

