# DNC List Search

Simple golang script to run a binary search through a dnc list file. 
Currently a work in progress as it includes the need to upload files and maintain them

## instructions

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

### Running

Example
```bash
./dnclist (201)0277-xxxx (201)0201xxxx 2010208xxxx
# (201)0277-xxxx Found in the file
# (201)0201xxxx Found in the file
# 2010208*xxxx Not Found in the file
```