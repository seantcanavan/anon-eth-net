## Anonymous Ethereum Network (anon-eth-net)

Totally anonymous botnet client with an emphasis on individual zombie control, resiliency of the host machine, and ease of remote code execution. Zombies mine ethereum in their free time for fun and check in at pre-specified intervals to report on their CPU, memory, disk, and network utilization.

Clients can also use [go-dos-yourself](https://github.com/seantcanavan/go-dos-yourself) to enable remote network performance testing, fuzzing, spoofing, and attacking. Use responsibly.

###Feature List:
1. Automate the execution of local processes remotely via REST or locally via JSON configuration.
2. Automatically receive local processes' log output directly via email or query the server for recent logs via REST.
3. Automatically prune old logs when disk space gets low.
4. Built-in network manager regularly queries configurable internet endpoints at set intervals and reboots the machine when it's unable to reach the internet.
5. Receive automated 'checkups' from the remote machine via email which include reports about overall system health.
6. Remote machines can be fully passively and fully anonymously monitored via a steady stream of emailed logs or directly managed via REST.
7. REST supports HTTPS with TLS encryption and timestamping to prevent replay attacks and ensure anonymity.

###Currently supported platforms:
- macOS El Capitan 10.11.6
- Ubuntu 16.04.1
- Windows 10

###Supported platforms wish list:
- macOS >= 10.8.X
- Ubuntu >= 14.04.5
- Windows >= 7 SP2

##General Setup for all OS's:
1. Update the required values in assets/config.json:
  1. CheckInGmailAddress - set this to the gmail address you wish to receive system reports and process logs at.
  2. CheckInGmailPassword - set this to the password to the above gmail address.
  3. CheckInFrequencySeconds - set this to the frequency at which you'd like to receive system reports at your specified email address. value is in seconds.
  4. NetQueryFrequencySeconds - set this to the frequency at which you'd like anon-eth-net to check for internet connectivity.
  5. UpdateFrequencySeconds - set this to the frequency at which you'd like anon-eth-net to check for updates to itself.
2. Optionally update the optional values in assets/config.json:
  1. DeviceName - set this to the canonical name of the device which will be executing anon-eth-net. e.g. "main desktop", "garage pc", "sister's laptop", etc.
  2. DeviceId - if you wish to use your own method of uniquely identifying your remote devices fill in that value here otherwise anon-eth-net will generate a GUID for you automatically.
3. Set your give gmail adress to allow "insecure app access". The page to enable that is here: https://support.google.com/accounts/answer/6010255?hl=en
4. Skip to the operating system of your choice and continue to setup process!
5. Download your favorite ethereum miner from the internet and add its install location to your system PATH variable.
6. Read the manual for the miner and configure it along with all the command-line parameters required for it to operate.
7. Update assets/main_loader_<targetos>.json with the command to start up the miner. An example is already located in assets/main_loader_linux.json to copy from.
8. You're done! Run the binary!

##Mac Code Complitaion Setup:

####Required Packages:
1. Git: `brew install git`
2. Go: `brew install go`
3. Glide: `brew install glide`
4. TBA

####Required Setup:
1. Use `git clone` to download the repository.
2. Go is already automatically configured on your local macOS path if installed via the brew package manager.
3. Setup GOBIN and GOPATH system variables for your macOS user. Point to the root of the anon-eth-net clone directory. For GOBIN use the same value as GOPATH but add /bin to the end. GOPATH should be something like `/Users/<username>/clone_root_dir/anon-eth-net`. You can automatically set the GOPATH variable in your .bash_profile file on each terminal load / system startup:
  1. `nano ~/.bash_profile`
  2. Scroll to the bottom and paste the following. You may wrap the folder argument after the equal sign in quotes if you have spaces in your folder names (shame on you).
  3. `export GOPATH=/Users/<username>/<clone_root_dir>/anon-eth-net`
  4. `export GOBIN=/Users/<username>/<clone_root_dir>/anon-eth-net/bin`
  5. I use a folder called 'workspace' as a former Eclipse user so my paths looks like this:
  6. `export GOPATH=/Users/seantcanavan/workspace/anon-eth-net`
  7. `export GOBIN=/Users/seantcanavan/workspace/anon-eth-net/bin`
4. Update your sudoers file with the following:
  1. `sudo visudo`
  2. `<username> ALL=(ALL) NOPASSWD:/usr/sbin/lsof`
  3. Mine looks like: `seantcanavan ALL=(ALL) NOPASSWD:/user/sbin/lsof`
  4. `<username> ALL=(ALL) NOPASSWD:/sbin/shutdown`
  5. Mine looks like: `seantcanavan ALL=(ALL) NOPASSWD:/sbin/shutdown`
5. By default your REST commands will be encrypted over HTTPS with a test certificate and private key that are readily available from this GitHub. This is fine for testing but when you decide to deploy you'll need to create your own private key / public key / certificate combination to keep all your transmissions totally secure.
  1. Generate the private key: `openssl genrsa -out server.key 2048`
  2. Generate the certificate and public key: `openssl req -new -x509 -sha256 -key server.key -out server.pem -days 3650`
  3. Place both files in the 'assets' folder which will be at `<clone_root_dir>/anon-eth-net/src/github.com/seantcanavan/assets/`
6. Change directory to the root of the source folder: `<clone_root_dir>/anon-eth-net/src/github.com/seantcanavan/`
7. `make install`
8. TBA

##Linux Code Compliation Setup:

####Required Packages:
1. Git: `sudo apt-get install git`
2. Go: `sudo apt-get install golang-go`
3. iostat: `sudo apt-get install sysstat`
4. Glide:
  1. `sudo add-apt-repository ppa:masterminds/glide && sudo apt-get update`
  2. `sudo apt-get install glide`
5. TBA

####Required Setup:
1. Use `git clone` to download the repository.
2. Go is already automatically configured in your local Ubuntu path if installed via the synaptic package manager.
3. Setup GOBIN and GOPATH system variables for your Ubuntu user. Point to the root of the anon-eth-net clone directory. For GOBIN use the same value as GOPATH but add /bin to the end. GOPATH should be something like `/home/<username>/clone_root_dir/anon-eth-net`. You can automatically set the GOPATH variable in your .bashrc file on each terminal load / system startup:
  1. `nano ~/.bashrc`
  2. Scroll to the bottom and paste the following. You may wrap the folder argument after the equal sign in quotes if you have spaces in your folder names (shame on you).
  3. `export GOPATH=/home/<username>/<clone_root_dir>/anon-eth-net`
  4. `export GOBIN=/home/<username>/<clone_root_dir>/anon-eth-net/bin`
  5. I use a folder called 'workspace' as a former Eclipse user so my paths looks like this:
  6. `export GOPATH=/home/seantcanavan/workspace/anon-eth-net`
  7. `export GOBIN=/home/seantcanavan/workspace/anon-eth-net/bin`
4. Update your sudoers file with the following:
  1. `sudo visudo`
  2. `<username> ALL=(ALL) NOPASSWD:/bin/netstat`
  3. Mine looks like: `seantcanavan ALL=(ALL) NOPASSWD:/bin/netstat`
5. By default your REST commands will be encrypted over HTTPS with a test certificate and private key that are readily available from this GitHub. This is fine for testing but when you decide to deploy you'll need to create your own private key / public key / certificate combination to keep all your transmissions totally secure.
  1. Generate the private key: `openssl genrsa -out server.key 2048`
  2. Generate the certificate and public key: `openssl req -new -x509 -sha256 -key server.key -out server.pem -days 3650`
6. Change directory to the root of the source folder: `<clone_root_dir>/anon-eth-net/src/github.com/seantcanavan/`
7. `make install`
8. TBA

##Windows Code Compliation Setup:


####Required Packages:
1. Git: `https://github.com/git-for-windows/git/releases/download/v2.10.2.windows.1/Git-2.10.2-64-bit.exe`
2. Golang: `https://storage.googleapis.com/golang/go1.7.3.windows-amd64.msi`
3. TBA - glide for windows? i don't think it's a thing unfortunately. will have to revert to 'go get'

####Required Setup:
1. Use `git clone` to download the repository.
2. Configure your local system to make sure that the Go installation path is on the system path.
3. Setup GOBIN and GOPATH system variables for your windows users. Point to the root of the anon-eth-net clone directory. For GOBIN use the same value as GOPATH but add \bin to the end. GOPATH should be something like `C:\users\<username>\clone_root_dir\anon-eth-net`.
4. I don't think Windows requires elevated permissions to execute the shutdown command natively from the shell as long as the executing user is an Administrator. Fingers crossed.
5. By default your REST commands will be encrypted over HTTPS with a test certificate and private key that are readily available from this GitHub. This is fine for testing but when you decide to deploy you'll need to create your own private key / public key / certificate combination to keep all your transmissions totally secure.
  1. Windows private key gen command here
  2. Windows certificate gen command here
6. TBA
