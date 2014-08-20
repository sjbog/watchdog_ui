watchdog_ui
===========

### Install
```git clone```

#### Dependencies
```
code.google.com/p/go.crypto/bcrypt
```

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


### Screenshots

![Main view](https://raw.github.com/sjbog/watchdog_ui/master/public/img/watchdog_ui_2.png "Main view")

![Server edit view](https://raw.github.com/sjbog/watchdog_ui/master/public/img/watchdog_ui_3.png "Server edit view")

![Server status view](https://raw.github.com/sjbog/watchdog_ui/master/public/img/watchdog_ui_4.png "Server status view")

![Login view](https://raw.github.com/sjbog/watchdog_ui/master/public/img/watchdog_ui_1.png "Login view")
