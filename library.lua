
TYPE_SIMPLE = "simple"
TYPE_XBUS = "xbus"

DIR_INPUT = "input"
DIR_OUTPUT = "output"

TERMINALS = {}

function create_terminal(name, index, type, direction, data)
    TERMINALS[name] = {type=type, index=index, direction=direction}
end