package graphdb

import (
    "fmt"
    "strings"
    "os"
    "bufio"
)

type commandsDirector struct {
    commands map[string]Command
    allNodes map[string]*Node
    rootNode *Node
    lastPrepared string
    successfulCommands []string
    storeCommand bool
}

func NewCommandsDirector() *commandsDirector {
    dir := &commandsDirector{make(map[string]Command), make(map[string]*Node), nil, "", make([]string, 0), false}

    dir.RegisterCommand(&HelpCommand{"help", dir})
    dir.RegisterCommand(&AllCommand{"all", dir})

    dir.RegisterCommand(&GraphCommand{"g", dir})
    dir.RegisterCommand(&NewCommand{"new", dir})
    dir.RegisterCommand(&ReparentCommand{"reparent", dir})
    dir.RegisterCommand(&ConnectCommand{"connect", dir})

    dir.RegisterCommand(&SaveCommand{"save", dir})
    dir.RegisterCommand(&LoadCommand{"load", dir})

    return dir
}

func (dir *commandsDirector) RegisterCommand(cmd Command) {
    dir.commands[cmd.getName()] = cmd
}

func (dir *commandsDirector) Execute(cmdName string, params []string) bool {
    if cmd, exists := dir.commands[cmdName]; exists {
        if cmd.validateParams(params) {
            dir.storeCommand = false
            status := cmd.execute(params)
            if status && dir.storeCommand {
                dir.successfulCommands = append(dir.successfulCommands, dir.lastPrepared)
            }
            return status
        }
        fmt.Println(cmdName + " " + cmd.getHelp())
        return false
    }
    fmt.Println("command does not exist")
    return false
}

func (dir *commandsDirector) Prepare(str string) []string {
    dir.lastPrepared = str
    return strings.Fields(str)
}

func (dir *commandsDirector) HasCommand(cmdName string) bool {
    _, exists := dir.commands[cmdName]
    return exists
}

type Command interface {
    execute([]string) bool
    getName() string
    validateParams([]string) bool
    getHelp() string
}

type HelpCommand struct {
    name string
    dir *commandsDirector
}
func (cmd *HelpCommand) execute(params []string) bool {
    if len(params) == 1 {
        fmt.Println(params[0] + " " + cmd.dir.commands[params[0]].getHelp())
    } else {
        fmt.Println("help " + cmd.getHelp())
    }
    return true
}
func (cmd *HelpCommand) getName() string {
    return cmd.name
}
func (cmd *HelpCommand) getHelp() string {
    str := ""
    str += "<command>\n"
    str += "available commands:\n"
    for cmdName, _ := range cmd.dir.commands {
        str += "\t" + cmdName + "\n"
    }
    return str
}
func (cmd *HelpCommand) validateParams(params []string) bool {
    if len(params) == 1 {
        if _, ok := cmd.dir.commands[params[0]]; ok {
            return true
        }
    } else if len(params) == 0 {
        return true
    }
    return false
}

type GraphCommand struct {
    name string
    dir *commandsDirector
}
func (cmd *GraphCommand) execute(params []string) bool {
    var newG *Node
    cmd.dir.storeCommand = true
    if _, ok := cmd.dir.allNodes[params[0]]; ok {
        newG = cmd.dir.allNodes[params[0]]
    } else {
        newG = NewGraph(params[0])
    }
    cmd.dir.rootNode = newG
    cmd.dir.allNodes[params[0]] = newG
    return true
}
func (cmd *GraphCommand) getName() string {
    return cmd.name
}
func (cmd *GraphCommand) getHelp() string {
    str := "<name>\n\tCreate and/or activate the graph named <name>"
    return str
}
func (cmd *GraphCommand) validateParams(params []string) bool {
    if len(params) == 1 {
        return true
    }
    return false
}

type AllCommand struct {
    name string
    dir *commandsDirector
}
func (cmd *AllCommand) execute(params []string) bool {
    for _, node := range cmd.dir.allNodes {
        fmt.Println(node)
    }
    return true
}
func (cmd *AllCommand) getName() string {
    return cmd.name
}
func (cmd *AllCommand) getHelp() string {
    str := "\n\tprint all nodes"
    return str
}
func (cmd *AllCommand) validateParams(params []string) bool {
    if len(params) == 0 {
        return true
    }
    return false
}

type NewCommand struct {
    name string
    dir *commandsDirector
}
func (cmd *NewCommand) execute(params []string) bool {
    cmd.dir.storeCommand = true
    node := cmd.dir.rootNode.NewSubGraph(params[0])
    cmd.dir.allNodes[params[0]] = node
    return true
}
func (cmd *NewCommand) getName() string {
    return cmd.name
}
func (cmd *NewCommand) getHelp() string {
    str := "<name>\n\tcreate new node named <name>"
    return str
}
func (cmd *NewCommand) validateParams(params []string) bool {
    if len(params) == 1 {
        if _, ok := cmd.dir.allNodes[params[0]]; ok {
            fmt.Println("node '" + params[0] + "' already exists")
            return false
        }
        return true
    }
    return false
}


type ReparentCommand struct {
    name string
    dir *commandsDirector
}
func (cmd *ReparentCommand) execute(params []string) bool {
    cmd.dir.storeCommand = true
    what := cmd.dir.allNodes[params[0]]
    to := cmd.dir.allNodes[params[1]]
    oldParent := what.parent
    if oldParent != nil {
        for i, child := range oldParent.subnodes {
            if child == what {
                oldParent.subnodes = append(oldParent.subnodes[:i], oldParent.subnodes[i+1:]...)
                break
            }
        }
    }
    what.parent = to
    to.subnodes = append(to.subnodes, what)
    return true
}
func (cmd *ReparentCommand) getName() string {
    return cmd.name
}
func (cmd *ReparentCommand) getHelp() string {
    str := "<what> <to>\n\tmove <what> to <to>'s children"
    return str
}
func (cmd *ReparentCommand) validateParams(params []string) bool {
    if len(params) == 2 {
        if _, ok := cmd.dir.allNodes[params[0]]; !ok {
            return false
        }
        if _, ok := cmd.dir.allNodes[params[1]]; !ok {
            return false
        }
        return true
    }
    return false
}

type SaveCommand struct {
    name string
    dir *commandsDirector
}
func (cmd *SaveCommand) execute(params []string) bool {
    f, _ := os.Create(params[0])
    defer f.Close()
    for _, line := range cmd.dir.successfulCommands {
        f.WriteString(line + "\n")
    }
    f.Sync()

    return true
}
func (cmd *SaveCommand) getName() string {
    return cmd.name
}
func (cmd *SaveCommand) getHelp() string {
    str := "<filename>\n\tsave all successful commands to file"
    return str
}
func (cmd *SaveCommand) validateParams(params []string) bool {
    if len(params) != 1 {
        return false
    }
    return true
}

type LoadCommand struct {
    name string
    dir *commandsDirector
}
func (cmd *LoadCommand) execute(params []string) bool {
    file, _ := os.Open(params[0])
    defer file.Close()
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        command := scanner.Text()
        fmt.Println("> " + command)
        fields := cmd.dir.Prepare(command)
        status := cmd.dir.Execute(fields[0], fields[1:])
        if !status {
            fmt.Println("Not OK")
            return false
        }
    }
    return true
}
func (cmd *LoadCommand) getName() string {
    return cmd.name
}
func (cmd *LoadCommand) getHelp() string {
    str := "<filename>\n\tload commands from saved file and replay them"
    return str
}
func (cmd *LoadCommand) validateParams(params []string) bool {
    if len(params) == 1 {
        if _, err := os.Stat(params[0]); os.IsNotExist(err) {
            fmt.Println("File does not exist")
            return false
        }
        return true
    }
    fmt.Println("invalid number of parameters")
    return false
}

type DFSCommand struct {
    name string
    dir *commandsDirector
}
func (cmd *DFSCommand) execute(params []string) bool {
    start := cmd.dir.allNodes[params[0]]
    it := NewDFSIterator(start)
    var prevNode *Node
    for node := range it.Stream() {
        fmt.Println("\t* " + node.String())
        if node == prevNode {
            return false
        }
        prevNode = node
    }
    return true
}
func (cmd *DFSCommand) getName() string {
    return cmd.name
}
func (cmd *DFSCommand) getHelp() string {
    str := "<name>\n\tprint the graph starting at <name> by DFS order"
    return str
}
func (cmd *DFSCommand) validateParams(params []string) bool {
    if len(params) == 1 {
        if _, ok := cmd.dir.allNodes[params[0]]; !ok {
            return false
        }
        return true
    }
    return false
}

type ConnectCommand struct {
    name string
    dir *commandsDirector
}
func (cmd *ConnectCommand) execute(params []string) bool {
    cmd.dir.storeCommand = true
    left := cmd.dir.allNodes[params[0]]
    right := cmd.dir.allNodes[params[2]]
    if params[1] == "-" {
        left.neighbours = append(left.neighbours, right)
        right.neighbours = append(right.neighbours, left)
    }
    if params[1] == "<" {
        right.neighbours = append(right.neighbours, left)
    }
    if params[1] == ">" {
        left.neighbours = append(left.neighbours, right)
    }
    return true
}
func (cmd *ConnectCommand) getName() string {
    return cmd.name
}
func (cmd *ConnectCommand) getHelp() string {
    str := "<node> [<, >, -] <node>\n\tconnect two nodes\n\tthe second parameter tells the direction"
    return str
}
func (cmd *ConnectCommand) validateParams(params []string) bool {
    if len(params) == 3 {
        missing := ""
        if _, ok := cmd.dir.allNodes[params[0]]; !ok {
            missing = "first"
        }
        if _, ok := cmd.dir.allNodes[params[2]]; !ok {
            if len(missing) != 0 {
                missing += " and second"
            } else {
                missing = "second"
            }
        }
        if !(params[1] == ">" || params[1] == "<" || params[1] == "-") {
            fmt.Println("connection must be >, < or -")
        }
        if len(missing) > 6 {
            fmt.Println(missing + " parameters missing")
            return false
        } else if len(missing) > 0 {
            fmt.Println(missing + " parameter missing")
            return false
        } else {
            return true
        }
    }
    fmt.Println("invalid number of parameters")
    return false
}

type FooCommand struct {
    name string
    dir *commandsDirector
}
func (cmd *FooCommand) execute(params []string) bool {
    return true
}
func (cmd *FooCommand) getName() string {
    return cmd.name
}
func (cmd *FooCommand) getHelp() string {
    str := ""
    return str
}
func (cmd *FooCommand) validateParams(params []string) bool {
    return false
}
