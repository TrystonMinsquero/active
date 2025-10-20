# Active

A command line tool that saves paths, and outputs the saved ones you asked it.
This is used in the command line to keep track of your active projects, making it easy to jump around in a jump list style.

Doesn't do much on it's own, but can use it with [a.sh](./a.sh) and `.zshrc` to make it easy and fast to keep track of recent projects. I'm sure there are better tools for this, especially with something like tmux. I had the idea and wanted to tinker with go a little bit more and was jumping between projects a lot and thought this would help.

## Usage

```
> a 	// go to most recent project
> a -2 	// go back to 2 recent projects
> a -l 	// list out last 10 projects (this is an L)
> a . 	// Save current directory as most recent project (can be any path)
> a -h 	// for more detailed usage
```

## Installation and setup
I don't think this tool is good enough to warrant putting it in a package manager,
so just build it with `go build` and add this to your `.zshrc` or `.bashrc` and replace the path for wherever you installed the project.
```sh
alias a="source <path_to_proj>/a.sh"
```

