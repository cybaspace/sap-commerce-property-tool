# SAP Commerce Properties Tool

This tool was created to make it easier for you to work with the SAP Commerce property files in several environments.

## Keyfeatures

  * split one big `local.properties` in several smaller files
  * handle different values for properties for several environments in csv files
  * Compare two files or environments for differences
  * read a single value or a set of values from the property files

## Overview

Intention of the tool was to work in more separate property files instead of one big `local.properties` file:  
A typical SAP Commerce `local.properties` file contains a few hundred up to more than thousand lines.  
With this tool you can easily split up the file in a lot of property files, e.g. one file for each connected system.

Besides that split up for more clarity you can use CSV files in which you can define different values for different environments.  
In the event you are using a test and a preproduction environment, you can have a properties csv file like this:
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
`yprops [<args>] <command>`

Commands are  
 * `generate`  
   Generates a `local.properties` file from multiple property and csv files
 * `get`  
   Read one property value or a set of values from a `local.properties` file or from multiple property and csv files
 * `diff`  
   Calculates the difference between two `local.properties` files or two defined environments
 * `list`
   List all defined systems in a csv property file

Arguments are
 * `-path=<filepath>`  
   Path there the property files are stored
 * `-files=<filelist>`  
   A comma separated list of files
 * `-system=<systemname>`  
   Name of the system/environment from the csv file(s)

## generate

This command generates a `local.properties` file out of a list of property and csv files.

### Usage
`yprops [-path=<filepath>] [-files=<filelist>] [-system=<systemname>] generate`

### Parameter
 * `path`  
   specify where the property and csv files are stored - the generated `local.properties` will also be stored there, e.g.  
   `-path=hybris/config`  

   If `path` is not specified the current working directory will be used.  
   **Caution:** If a `local.properties` file in the specified path already exists it will be overwritten without a question.
 * `filelist`  
   Comma separated list of property and csv files, e.g.  
   `-files=application.properties,sap.properties.csv`  

   If `filelist` is not specified there must be a file `property-files` in the `path` which lists the files - each per line.
 * `systemname`  
   The name of the system as defined in the first line of the csv files, e.g.  
   `-systemname=preproduction`  

   If there are more than one csv file a column for the system **must** exist in all csv files.
   If no csv files are used, this parameter can be omitted.

### Description

Property files **must** have the prefix `.properties` and csv files the prefix `.csv`.

The first line in a csv file must define the system names.  
As the first column contains the property names, this column in the first line will be omitted.  
Each system name must exist in each csv file, if more than one is used. The order of the systems can differ between different csv files.

## get

Gets a single value or a set of values from a `local.properties` file or a set of files.

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

## diff

Calculates the difference between two property files or two systems in a csv file or list of files.

### Usage
`yprops [-path=<path>] [-files=<filelist>] [-system=<systemnames>] diff`

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

## list

List all systems from specified csv file(s).

### Usage
`yprops [-path=<path>] [-files=<filelist>] list`

### Parameter
 * `path`  
   specify where csv file(s) are stored.  
   If `path` is not specified, the current working directory will be used.
 * `filelist`  
   Comma separated list of one or more csv files, e.g.  
   `-files=sap.properties.csv`  
   If `filelist` is not specified there must be a file `property-files` in the `path` which lists the files - each per line. At least one of the files must be a csv file.

### Example
 * `yprops list`
 * `yprops list -path hybris/config`  
 * `yprops list -files hybris/config/sap.properties.csv`
 * `yprops list -path hybris/config -files
 