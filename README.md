watchdog_ui
===========

### Install
```git clone```

Default account: ```root:1234```

### Servers config file

Examle of conf/servers.conf

```
[Server1]
host	= "s1.hostname.com"
port	= "22"

username	= "root"
password	= "root"
private_key	= ""
# Auth via password or private key

#query_interval	= 60
# Every N seconds server will be queried

[Server1/commands]
# Label = "cmd"

nginx	= "service nginx status"
php-fpm	= "service php5-fpm status"

```