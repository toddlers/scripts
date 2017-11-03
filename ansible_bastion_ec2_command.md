#### Using Bastion aka jump host in EC2 to provision machines inside the private VPC.

Please follow below steps and place the configuration accordingly.

*  Create ssh `~/.ssh/ansible_ssh_config` file like this.
  
```
Host bastion
  Hostname jump.ec2.example.io
  User ubuntu
  IdentityFile ~/.ssh/id_rsa.pem
  PasswordAuthentication no
  ForwardAgent yes
  ServerAliveInterval 60
  TCPKeepAlive yes
  ControlMaster auto
  ControlPath ~/.ssh/ansible-%r@%h:%p
  ControlPersist 15m
  ProxyCommand none
  LogLevel QUIET

Host *
  User ubuntu
  IdentityFile ~/.ssh/id_rsa
  ServerAliveInterval 60
  TCPKeepAlive yes
  ProxyCommand ssh -q -A ubuntu@bastion nc %h %p
  LogLevel QUIET
  StrictHostKeyChecking no
```

* `ansible.cfg` should have this 
```
[ssh_connection]
ssh_args = -F ~/.ssh/ansible_ssh_config -o ControlMaster=auto -o ControlPersist=5m -o LogLevel=QUIET
control_path = ~/.ssh/ansible-%%r@%%h:%%p
```

* Add the entry for the host in `inventories/main`

* Test the connection using below command

```
git:master *%=⚡ ansible -i inventories/main test.example.com  -m ping
test.example.com | SUCCESS => {
    "changed": false,
    "ping": "pong"
}
```
