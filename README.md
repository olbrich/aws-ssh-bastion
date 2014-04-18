aws-ssh-bastion
===============

tools for sshing into amazon instances through a bastion server

The 'aws' go tool lets you lookup an instance's private or public ip address by Name.

This can be incorporated into your ~/.ssh/config so that you can seemlessly ssh to a server after passing through the appropriate bastion server.

compile the go tool and make sure it's in your path.  It looks for the usual AWS environment variables to authenticate.

Examples
========

~/.ssh/config

    # don't do any relaying when trying to hit an instance directly
    Host *.amazonaws.com
      ProxyCommand none
    # authenticates as 'user' -- this part might need some tweaking depending on your setup
    Host *-bastion
      User user
      HostbasedAuthentication yes
      ProxyCommand ssh -A -l 'user' -q -p '%p' $(aws -p %h) -W $(aws %h):%p
    # pull name of environment from the hostname and use that bastion as a gateway
    # example: staging-web would try to ssh through the 'staging-bastion' to the instance named 'staging-web' 
    Host *-*
      User user
      HostbasedAuthentication yes
      ProxyCommand ssh -A -l 'user' -q -p '%p' $(aws -p $(echo %h | cut -f1 -d-)-bastion) -W $(aws %h):%p
