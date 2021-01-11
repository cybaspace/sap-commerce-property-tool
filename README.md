# SAP Commerce Properties Tool

This tool was created to make it easier for you to work with the SAP Commerce Property files in several environments.

Intention of the tool was to work in more separate property files instead of one big `local.properties` file.  
Besides that split up for more clarity you can use CSV files in which you can define different values for different environments.  
In the event you are using a test and a preproduction environment, you can have a properties-csv-file like this:
```
systems         ; test         ; preproduction
first.property  ; test value 1 ; preproduction value 1
second.property ; test value 2 ; preproduction value 2
```
After generating the `local.properties` file with this csv file for the environment preproduction you will see this:
```
first.property=preproduction value 1
second.property=preproduction value 1
```
This will significantly reduce your workload while working with a bunch of environments.

## Usage
`yprops <command> [<args>]`

Commands are  
 * `generate`  
   Generates a `local.properties` file from multiple property and csv files
 * `get`  
   Read one property value from a `local.properties` file or from multiple property and csv files
 * `diff`  
   Calculates the difference between two `local.properties` files
 * `list`
   List all defined systems in a csv property file

## generate
### Usage
`yprops [-path=<path>] [-files=<filelist>] [-system=<systemname>] generate`

### Parameter
 * `path`  
   specify where the property and csv files are stored - there will also the generated `local.properties` will be stored.  
   If `path` is not specified the current working directory will be used.  
   If a `local.properties` file in the specified path already exists it will be overwritten without a question.
 * `filelist`  
   Comma separated list of property and csv files, e.g.  
   `application.properties,sap.properties.machine.properties.csv`  
   If `filelist` is not specified there must be a file `property-files` in the `path` which lists the files - each per line.
 * `systemname`  
   The name of the system as defined in the first line of the csv files.  
   If there are more than one csv file a column for the system **must** exist in all csv files.
   If no csv files are used, this parameter can be omitted.

This command generates a `local.properties` file out of a list of property and csv files.

### Description

Property files **must** have the prefix `.properties` and csv files the prefix `.csv`.

The first line in a csv file must define the system names.  
As the first column contains the property names, this column in the first line will be omitted.  
Each system name must exist in each csv file, if more than one is used. The order of the systems can differ between different csv files.

## generate
### Usage
`yprops [-path=<path>] [-propfiles=<filelist>] [-system=<systemname>] get <property>`

### Parameter
 * `path`  
   specify where the property and csv files are stored - there will also the generated `local.properties` will be stored.  
   If `path` is not specified the current working directory will be used.  
   If a `local.properties` file in the specified path already exists it will be overwritten without a question.
 * `filelist`  
   Comma separated list of property and csv files, e.g.  
   `application.properties,sap.properties.machine.properties.csv`  
   If `filelist` is not specified there must be a file `property-files` in the `path` which lists the files - each per line.
 * `systemname`  
   The name of the system as defined in the first line of the csv files.  
   If there are more than one csv file a column for the system **must** exist in all csv files.
   If no csv files are used, this parameter can be omitted.

## generate
### Usage
`yprops [-path=<path>] [-files=<filelist>] [-system=<systemname>] diff`

### Parameter
 * `path`  
   specify where the property and csv files are stored - there will also the generated `local.properties` will be stored.  
   If `path` is not specified the current working directory will be used.  
   If a `local.properties` file in the specified path already exists it will be overwritten without a question.
 * `filelist`  
   Comma separated list of property and csv files, e.g.  
   `application.properties,sap.properties.machine.properties.csv`  
   If `filelist` is not specified there must be a file `property-files` in the `path` which lists the files - each per line.
 * `systemname`  
   The name of the system as defined in the first line of the csv files.  
   If there are more than one csv file a column for the system **must** exist in all csv files.
   If no csv files are used, this parameter can be omitted.

## generate
### Usage
`yprops [-path=<path>] [-files=<filelist>] list`

### Parameter
 * `path`  
   specify where the property and csv files are stored - there will also the generated `local.properties` will be stored.  
   If `path` is not specified the current working directory will be used.  
   If a `local.properties` file in the specified path already exists it will be overwritten without a question.
 * `filelist`  
   Comma separated list of property and csv files, e.g.  
   `application.properties,sap.properties.machine.properties.csv`  
   If `filelist` is not specified there must be a file `property-files` in the `path` which lists the files - each per line.
 * `systemname`  
   The name of the system as defined in the first line of the csv files.  
   If there are more than one csv file a column for the system **must** exist in all csv files.
   If no csv files are used, this parameter can be omitted.

### Example
 * `yprops list -path hybris/config`  
   und
 * `yprops list -propfiles hybris/config/machine.properties.csv`  
   list all systems from specified csv file