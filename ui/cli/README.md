## SParts Command Line Interface

### Overview

The SParts Command Line interface (cli) is the command line tool to interface with the ledger. For example you can list out the suppliers, create new software parts and associate artifacts with existing parts. It was designed to have a similar look, feel  and behavior as the git command line interface. Enter the following to get a list of available commands

```
sparts --help
```

Examples of some of the more common commands include:

```
  sparts <command> --help			// obtain detailed help description for
 								  // <command>'s 
  										
  sparts config --list				// list the sparts local working direct
  
  sparts ping					   // checks if ledger is up
  
  sparts supplier --list		    // list the current suppliers

  sparts add <file1> <file2>		// Add artifacts to the staging areas
  
  sparts status 				   // Display artifacts stagedstatus for posting to ledger
  
  sparts part --create 				// create a new part for the default supplier
```
### To Build

You will need to install the following third party components:

```
go get github.com/mattn/go-sqlite3 [MIT]  // might need to install gcc on Windows
go get github.com/nu7hatch/gouuid  [MIT]
```

Edit build.sh

```
# Set to directory within your $PATH (otherwise build in local directory)
BIN_DIR="."
```

Run build.sh