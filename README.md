# DNC List Search

Simple golang script to run a binary search through a dnc list file.
Currently a work in progress as it includes the need to upload files and maintain them

## Installation

### Homebrew

With brew installed, you can run the following command to install the binary along with adding `dnclistsearch` to your working path
```bash
brew install keyglee/keyglee/dnclistsearch
```
or
```bash
brew tap keyglee/keyglee
brew install dnclistsearch
```

### Manual

Close or download the repository to get the source files and run
`go build` in the same directory to build the

in the folder, there needs to be a `dnc.txt` file that contains a dnc list

```bash
head dnc.txt
>> 201,0000xxx
>> 201,0000xxx
>> 201,0020xxx
>> 201,0087xxx
>> 201,0201xxx
>> 201,0207xxx
>> 201,0248xxx
>> 201,0249xxx
>> 201,0249xxx
>> 201,0277xxx
# the xxx's are supposed to be numbers but they are omitted
```


## Usage

Example

```bash
dnclistsearch (201)0277-xxxx (201)0201xxxx 2010208xxxx
# (201)0277-xxxx Found in the file
# (201)0201xxxx Found in the file
# 2010208xxxx Not Found in the file
```

## Command Line Options

| Flag         | Description                        | Default       |
| ------------ | ---------------------------------- | ------------- |
| `-speedrun`  | Runs program with execution timer  | false         |
| `-log`       | Set log level (error, info, debug) | "error"       |
| `-pretty`    | Enable pretty print output         | false         |
| `-delimiter` | Response delimiter                 | "none"        |
| `-separator` | Response separator                 | "newline"     |
| `-csv`       | Path to input CSV file             | ""            |
| `-output`    | Path to output CSV file            | "results.csv" |

## Usage Examples

```bash
# Process a single phone number
dnclistsearch -pretty 1234567890

# Process CSV file
dnclistsearch -csv input.csv -output results.csv

# Run with debug logging
dnclistsearch -log debug -pretty 1234567890
```

## Output

Output Formats

- Default: Raw output
- Pretty Print: Formatted with found/not found messages
- CSV: Results written to specified output file
