## Anonymous Ethereum Network (anon-eth-net)

Totally anonymous botnet client with an emphasis on individual zombie control, resiliency of the host machine, and ease of remote code execution. Zombies mine ethereum in their free time for fun and check in at pre-specified intervals to report on their CPU, memory, disk, and network utilization.

Clients can also use [go-dos-yourself](https://github.com/seantcanavan/go-dos-yourself) to enable remote network performance testing, fuzzing, spoofing, and attacking. Use responsibly.

Feature List:
01. Automatic logging of all processes that are executed. Log files are automatically pruned to conserve disk space after a configurable amount of time.
02. Easily configurable via one simple JSON file. Up and running after setting only 5 easy variables
03. Arbitrary process loader can run processes synchronously or asynchronously. Output is automatically captured and logged to disk. Command output can also be emailed after execution finishes.
04. Built in network manager can monitor the internet connectivity of all clients at a regular interval and automatically reboot the machine as necessary in order to restore a broken link state.
05. Built in system profiler monitors all key aspects of the client machine such as disk space, CPU and memory performance, open ports, network interfaces, kernel version, etc... Reports are automatically generated at set intervals and emailed directly to the user.
06. Extremely robust and secure HTTPS REST interface. Code can be directly passed to the client in binary, command line script, or python format and directly executed. Logs are generated and sent back directly to the user after execution finishes. Configuration files can also be updated, logs can be collected / purged, and the machine can be rebooted all on command.
07. Work in progress - a robust self-updating system which ensures that anon-eth-net is always running the most up-to-date version.

###Currently supported platforms:
- macOS El Capitan 10.11.6
- Ubuntu 16.04.1
- Windows 10 (work in progress)

###Supported platforms wish list:
- macOS >= 10.8.X
- Ubuntu >= 14.04.5
- Windows >= 7 SP2

###Glossary of Terms:
01. `<clone_root_dir>` : This is the folder that you've cloned this project into or the folder one up the directory tree from `anon-eth-net` after `git clone` is executed.
02. `<username>` : This is your local username specific to your operating system and the user you are currently logged in as. Each individual operating system stores it in a different location but the name should be familiar. The folder named after your username will contain all your personal files.
  01. Mac: `/Users/` contains all the user folders
  02. Windows: `C:\Users\` contains all the user folders
  03. Ubuntu: `/home/` contains all the user folders
03. `<emaillogindetails.txt>` : This is a two line file. Line one contains a valid gmail login and line two contains the password associated with that account. You create this and no one else sees it. Insecure app access must be enabled on the gmail account: https://support.google.com/accounts/answer/6010255?hl=en. The purpose of this file is to enable automated gmail reports sent by all machines this software is installed on. The gmail address acts as a centralized database of logs for every single instance of anon-eth-net. This is how it remains 'anonymous'. You never need to log into the remote machine once it's setup and you can use a VPN to access a totally anonymous gmail address.

##Mac Setup:

####Required Packages:
01. Git: `brew install git`
02. Go: `brew install go`
03. Glide: `brew install glide`
04. TBA

####Required Setup:
01. Use `git clone` to download the repository.
02. Go is already automatically configured on your local macOS path if installed via the brew package manager.
03. Setup GOBIN and GOPATH system variables for your macOS user. Point to the root of the anon-eth-net clone directory. For GOBIN use the same value as GOPATH but add /bin to the end. GOPATH should be something like `/Users/<username>/clone_root_dir/anon-eth-net`. You can automatically set the GOPATH variable in your .bash_profile file on each terminal load / system startup:
  01. `nano ~/.bash_profile`
  02. Scroll to the bottom and paste the following WITHOUT quotes. You may wrap the folder argument after the equal sign in quotes if you have spaces in your folder names (shame on you).
  03. `export GOPATH=/Users/<username>/<clone_root_dir>/anon-eth-net`
  04. `export GOBIN=/Users/<username>/<clone_root_dir>/anon-eth-net/bin`
  05. I use a folder called 'workspace' as a former Eclipse user so my paths looks like this:
  06. `export GOPATH=/Users/seantcanavan/workspace/anon-eth-net`
  07. `export GOBIN=/Users/seantcanavan/workspace/anon-eth-net/bin`
04. Update your sudoers file with the following:
  01. `sudo visudo`
  02. `<username> ALL=(ALL) NOPASSWD:/usr/sbin/lsof`
  03. Mine looks like: `seantcanavan ALL=(ALL) NOPASSWD:/user/sbin/lsof`
  04. `<username> ALL=(ALL) NOPASSWD:/sbin/shutdown`
  05. Mine looks like: `seantcanavan ALL=(ALL) NOPASSWD:/sbin/shutdown`
05. Create `emaillogindetails.txt` inside the `assets` folder at `<clone_root_dir>/anon-eth-net/src/github.com/seantcanavan/assets/`.
06. Add your gmail address to line 1 and gmail password to line 2.
07. By default your REST commands will be encrypted over HTTPS with a test certificate and private key that are readily available from this GitHub. This is fine for testing but when you decide to deploy you'll need to create your own private key / public key / certificate combination to keep all your transmissions totally secure.
  01. Generate the private key: `openssl genrsa -out server.key 2048`
  02. Generate the certificate and public key: `openssl req -new -x509 -sha256 -key server.key -out server.pem -days 3650`
  03. Place both files in the 'assets' folder which will be at `<clone_root_dir>/anon-eth-net/src/github.com/seantcanavan/assets/`
08. Change directory to the root of the source folder: `<clone_root_dir>/anon-eth-net/src/github.com/seantcanavan/`
09. `make install`
10. TBA

##Linux Setup:

####Required Packages:
01. Git: `sudo apt-get install git`
02. Go: `sudo apt-get install golang-go`
03. iostat: `sudo apt-get install sysstat`
04. Glide:
  01. `sudo add-apt-repository ppa:masterminds/glide && sudo apt-get update`
  02. `sudo apt-get install glide`
05. TBA

####Required Setup:
01. Use `git clone` to download the repository.
02. Go is already automatically configured in your local Ubuntu path if installed via the synaptic package manager.
03. Setup GOBIN and GOPATH system variables for your Ubuntu user. Point to the root of the anon-eth-net clone directory. For GOBIN use the same value as GOPATH but add /bin to the end. GOPATH should be something like `/home/<username>/clone_root_dir/anon-eth-net`. You can automatically set the GOPATH variable in your .bashrc file on each terminal load / system startup:
  01. `nano ~/.bashrc`
  02. Scroll to the bottom and paste the following WITHOUT quotes. You may wrap the folder argument after the equal sign in quotes if you have spaces in your folder names (shame on you).
  03. `export GOPATH=/home/<username>/<clone_root_dir>/anon-eth-net`
  04. `export GOBIN=/home/<username>/<clone_root_dir>/anon-eth-net/bin`
  05. I use a folder called 'workspace' as a former Eclipse user so my paths looks like this:
  06. `export GOPATH=/home/seantcanavan/workspace/anon-eth-net`
  07. `export GOBIN=/home/seantcanavan/workspace/anon-eth-net/bin`
04. Update your sudoers file with the following:
  01. `sudo visudo`
  02. `<username> ALL=(ALL) NOPASSWD:/bin/netstat`
  03. Mine looks like: `seantcanavan ALL=(ALL) NOPASSWD:/bin/netstat`
05. Create `emaillogindetails.txt` inside the `assets` folder at `<clone_root_dir>/anon-eth-net/src/github.com/seantcanavan/assets/`.
06. Add your gmail address to line 1 and gmail password to line 2.
07. By default your REST commands will be encrypted over HTTPS with a test certificate and private key that are readily available from this GitHub. This is fine for testing but when you decide to deploy you'll need to create your own private key / public key / certificate combination to keep all your transmissions totally secure.
  01. Generate the private key: `openssl genrsa -out server.key 2048`
  02. Generate the certificate and public key: `openssl req -new -x509 -sha256 -key server.key -out server.pem -days 3650`
08. Change directory to the root of the source folder: `<clone_root_dir>/anon-eth-net/src/github.com/seantcanavan/`
09. `make install`
10. TBA

##Windows Setup:


####Required Packages:
01. Git: `https://github.com/git-for-windows/git/releases/download/v2.10.2.windows.1/Git-2.10.2-64-bit.exe`
02. Golang: `https://storage.googleapis.com/golang/go1.7.3.windows-amd64.msi`
03. TBA - glide for windows? i don't think it's a thing unfortunately. will have to revert to 'go get'

####Required Setup:
01. Use `git clone` to download the repository.
02. Configure your local system to make sure that the Go installation path is on the system path.
03. Setup GOBIN and GOPATH system variables for your windows users. Point to the root of the anon-eth-net clone directory. For GOBIN use the same value as GOPATH but add \bin to the end. GOPATH should be something like `C:\users\<username>\clone_root_dir\anon-eth-net`.
04. I don't think Windows requires elevated permissions to execute the shutdown command natively from the shell as long as the executing user is an Administrator. Fingers crossed.
05. Create `emaillogindetails.txt` inside the `assets` folder at `<clone_root_dir>\anon-eth-net\src\github.com\seantcanavan\assets\`.
06. Add your gmail address to line 1 and gmail password to line 2.
07. By default your REST commands will be encrypted over HTTPS with a test certificate and private key that are readily available from this GitHub. This is fine for testing but when you decide to deploy you'll need to create your own private key / public key / certificate combination to keep all your transmissions totally secure.
  01. Windows private key gen command here
  02. Windows certificate gen command here
08. TBA
