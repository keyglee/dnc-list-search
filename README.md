# DNC List Search

Simple golang script to run a binary search through a dnc list file. 
Currently a work in progress as it includes the need to upload files and maintain them

## Installation

### Build
```go build```

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

#### Directory example
```bash
ls-al
-rw-r--r--   1 User  root         722 Jan 31 15:38 README.md
-rwx------   1 User  root  3057397056 Jan 31 14:23 dnc.txt
-rwxr-xr-x   1 User  root     3000658 Jan 31 14:56 dnclist
-rw-r--r--   1 User  root         142 Jan 31 14:33 go.mod
-rw-r--r--   1 User  root        1390 Jan 31 14:33 go.sum
-rw-r--r--   1 User  root        3036 Jan 31 15:25 main.go
```

## Usage

Example
```bash
./dnclist (201)0277-xxxx (201)0201xxxx 2010208xxxx
# (201)0277-xxxx Found in the file
# (201)0201xxxx Found in the file
# 2010208*xxxx Not Found in the file
```


## Command Line Options

| Flag | Description | Default |
|------|-------------|---------|
| `-new` | Adds a new card from text file | "" |
| `-speedrun` | Runs program with execution timer | false |
| `-log` | Set log level (error, info, debug) | "error" |
| `-pretty` | Enable pretty print output | false |
| `-delimiter` | Response delimiter | "none" |
| `-separator` | Response separator | "newline" |
| `-csv` | Path to input CSV file | "" |
| `-output` | Path to output CSV file | "results.csv" |

## Usage Examples

```bash
# Process a single phone number
./dnclist -pretty 1234567890

# Process CSV file
./dnclist -csv input.csv -output results.csv

# Add new numbers to DNC list
./dnclist -new newlist.txt

# Run with debug logging
./dnclist -log debug -pretty 1234567890
```

## Output
Output Formats
- Default: Raw output
- Pretty Print: Formatted with found/not found messages
- CSV: Results written to specified output file