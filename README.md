## cb-log
[![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/cloud-barista/cb-log?label=go.mod)](https://github.com/cloud-barista/cb-log/blob/master/go.mod)
[![GoDoc](http://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)](https://pkg.go.dev/github.com/cloud-barista/cb-log@master)&nbsp;&nbsp;&nbsp;
[![Release Version](https://img.shields.io/github/v/release/cloud-barista/cb-log)](https://github.com/cloud-barista/cb-log/releases)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://github.com/cloud-barista/cb-log/blob/master/LICENSE)

CB-Log is the logger library for the Cloud-Barista Multi-Cloud Framework.


```
[NOTE]
cb-log is currently under development. (the latest version is 0.3.0 espresso)
So, we do not recommend using the current release in production.
Please note that the functionalities of cb-log are not stable and secure yet.
If you have any difficulties in using cb-log, please let us know.
(Open an issue or Join the cloud-barista Slack)
```
***


## How to use CB-Log library in a project WITHOUT `go module`
### 1. install CB-Log library pkg
- $ go get github.com/cloud-barista/cb-log  
- export CBLOG_ROOT=$GOPATH/src/github.com/cloud-barista/cb-log
    
### 2. example
- https://github.com/cloud-barista/cb-log/blob/master/test/sample.go

### 3. test example
- $ cd $CBLOG_ROOT/test  
- $ go run sample.go   `# loglevel: debug in $CBLOG_ROOT/conf/log_conf.yaml`
  
      ```
      [CB-SPIDER].[INFO]: 2019-08-16 23:22:51 sample.go:25, main.main() - start.........
      [CB-SPIDER].[INFO]: 2019-08-16 23:22:51 sample.go:45, main.createUser1() - start creating user.
      [CB-SPIDER].[INFO]: 2019-08-16 23:22:51 sample.go:59, main.createUser1() - finish creating user.
      [CB-SPIDER].[INFO]: 2019-08-16 23:22:51 sample.go:64, main.createUser2() - start creating user.
      [CB-SPIDER].[ERROR]: 2019-08-16 23:22:51 sample.go:69, main.createUser2() - DBMS Session is closed!!
      [CB-SPIDER].[INFO]: 2019-08-16 23:22:51 sample.go:78, main.createUser2() - finish creating user.
      [CB-SPIDER].[INFO]: 2019-08-16 23:22:51 sample.go:37, main.main() - end.........

      [CB-SPIDER].[INFO]: 2019-08-16 23:22:53 sample.go:25, main.main() - start.........
      [CB-SPIDER].[INFO]: 2019-08-16 23:22:53 sample.go:45, main.createUser1() - start creating user.
      [CB-SPIDER].[INFO]: 2019-08-16 23:22:53 sample.go:59, main.createUser1() - finish creating user.
      [CB-SPIDER].[INFO]: 2019-08-16 23:22:53 sample.go:64, main.createUser2() - start creating user.
      [CB-SPIDER].[ERROR]: 2019-08-16 23:22:53 sample.go:69, main.createUser2() - DBMS Session is closed!!
      [CB-SPIDER].[INFO]: 2019-08-16 23:22:53 sample.go:78, main.createUser2() - finish creating user.
      [CB-SPIDER].[INFO]: 2019-08-16 23:22:53 sample.go:37, main.main() - end.........
      ```
      

- set Log Level: `error`
  -	$ vi $CBLOG_ROOT/conf/log_conf.yaml
      <br>`loglevel: debug` => `loglevel: error`  
    
      ```
      [CB-SPIDER].[ERROR]: 2019-08-16 23:22:57 sample.go:69, main.createUser2() - DBMS Session is closed!!

      [CB-SPIDER].[ERROR]: 2019-08-16 23:22:59 sample.go:69, main.createUser2() - DBMS Session is closed!!
      ```

## How to use CB-Log library in a project WITH `go module`
You would not need to install CB-Log by `go get github.com/cloud-barista/cb-log` because of `go module`.
### 1. Setup log_conf.yaml
- Make a directory for log_conf.yaml (if necessary)
  
  - e.g.) ```mkdir $YOUR_PROJECT_DIRECTORY/configs```
  
- Create `log_conf.yaml` below

  ```yaml
  #### Config for CB-Log Lib. ####

  cblog:
    ## true | false
    loopcheck: true # This temp method for development is busy wait. cf) cblogger.go:levelSetupLoop().

    ## debug | info | warn | error
    loglevel: debug # If loopcheck is true, You can set this online.

    ## true | false
    logfile: false 

  ## Config for File Output ##
  logfileinfo:
    filename: ./log/cblogs.log
    # filename: $CBLOG_ROOT/log/cblogs.log
    maxsize: 10 # megabytes
    maxbackups: 50
    maxage: 31 # days
  ```
  
- Set and input config path

  ```go
  var cblogger *logrus.Logger

  func init() {
    // cblog is a global variable.
    filePath := filepath.Join("..", "conf", "log_conf.yaml")
    cblogger = cblog.GetLogger("CB-SPIDER", filePath)
  }
  ```
  
### 2. Example
- https://github.com/cloud-barista/cb-log/blob/master/test/sample-with-config-path.go

### 3. Test and result
- $ cd $CBLOG_ROOT/test
  
- $ go run sample-with-config-path.go
  
      ```
      [CB-SPIDER ..\conf\log_conf.yaml]
      [CB-SPIDER].[INFO]: 2020-12-23 17:46:09 sample-with-config-path.go:27, main.main() - start.........
      [CB-SPIDER].[INFO]: 2020-12-23 17:46:09 sample-with-config-path.go:48, main.createUser3() - start creating user.
      [CB-SPIDER].[DEBUG]: 2020-12-23 17:46:09 sample-with-config-path.go:58, main.createUser3() - msg for debugging msg!!
      [CB-SPIDER].[INFO]: 2020-12-23 17:46:09 sample-with-config-path.go:63, main.createUser3() - finish creating user.
      [CB-SPIDER].[DEBUG]: 2020-12-23 17:46:09 sample-with-config-path.go:30, main.main() - msg for debugging msg!!
      [CB-SPIDER].[INFO]: 2020-12-23 17:46:09 sample-with-config-path.go:68, main.createUser4() - start creating user.
      [CB-SPIDER].[ERROR]: 2020-12-23 17:46:09 sample-with-config-path.go:73, main.createUser4() - DBMS Session is closed!!
      [CB-SPIDER].[INFO]: 2020-12-23 17:46:09 sample-with-config-path.go:82, main.createUser4() - finish creating user.
      [CB-SPIDER].[INFO]: 2020-12-23 17:46:09 sample-with-config-path.go:40, main.main() - end.........
      ```
      
- set Log Level: `debug` => `error`   
  - $ vi ../conf/log_conf.yaml
    
      ```
      [CB-SPIDER].[ERROR]: 2020-12-23 18:08:12 sample-with-config-path.go:73, main.createUser4() - DBMS Session is closed!!

      [CB-SPIDER].[ERROR]: 2020-12-23 18:08:14 sample-with-config-path.go:73, main.createUser4() - DBMS Session is closed!!
      ```

