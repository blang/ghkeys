ghkeys 
======

ghkeys is a tiny utility to get the public ssh keys of a github user.

Usage
-----
```bash
$ go get github.com/blang/ghkeys

$ ghkeys blang
ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQD8uS2zFn...
ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDa22h62I...

// Now you can give a user access to your machine 
$ ghkeys blang >> ~/.ssh/authorized_keys
```

Contribution
-----

Feel free to make a pull request. For bigger changes create a issue first to discuss about it.


License
-----

See [LICENSE](LICENSE) file.
