Instructions to create and work with the workspace

Create the workspace folder
```bash
mkdie workspacenae
```
cd inside the workspace folder
```bash
cd workspace
```
Initiate module as:
```bash
go work init ./modulename
```
For example:
```bash
go work init ./walistner
```
The module folder will be created as empty folder, and the `go.work` file will looks something like:
```bash
go 1.18

use ./walistner
```
cd inside the module folder
```bash
cd modulename
```
Run the mod init command
```bash
go mod init
```
Run the mod tidy command
```bash
go mod tidy
```

To run any module inside the workspace use:
```bash
go run ./modulename
```
As:
```bash
go run ./walistner
```