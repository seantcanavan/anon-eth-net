## Anonymous Ethereum Network (anon-eth-net)

Totally anonymous botnet client with an emphasis on individual zombie control, resiliency of the host machine, and ease of remote code execution. Zombies mine ethereum in their free time for fun and check in at pre-specified intervals to report on their CPU, memory, disk, and network utilization. Can be combined with go-dos-yourself to enable remote network performance testing, fuzzing, spoofing, and attacking. Use responsibly.

###Currently supported platforms:
- macOS El Capitan 10.11.6
- Ubuntu 16.04.1
- Windows 10 (work in progress)

###Supported platforms wish list:
- macOS >= 10.8.X
- Ubuntu >= 14.04.5
- Windows >= 7 SP2

##Mac Setup:

####Required Packages:
1. Git: `brew install git`
2. Go: `brew install go`
3. TBA

####Required Setup:
1. Use `git clone` to download the repository.
2. Go is already automatically configured on your local macOS path if installed via the brew package manager.
3. Setup GOBIN and GOPATH system variables for your macOS user. Point to the root of the anon-eth-net clone directory. For GOBIN use the same value as GOPATH but add /bin to the end. GOPATH should be something like `/Users/<username>/clone_root_dir/anon-eth-net`. You can automatically set the GOPATH variable in your .bash_profile file on each terminal load / system startup:
  1. `nano ~/.bash_profile`
  2. Scroll to the bottom and paste the following WITHOUT quotes. You may wrap the folder argument after the equal sign in quotes if you have spaces in your folder names (shame on you).
  3. `export GOPATH=/Users/<username>/<clone_root_dir>/anon-eth-net`
  4. `export GOBIN=/Users/<username>/<clone_root_dir>/anon-eth-net/bin`
  5. I use a folder called 'workspace' as a former Eclipse user so my paths looks like this:
  6. `export GOPATH=/Users/seantcanavan/workspace/anon-eth-net`
  7. `export GOBIN=/Users/seantcanavan/workspace/anon-eth-net/bin`
4. Update your sudoers file with the following:
  1. `sudo visudo`
  2. `<your-user-name-here> ALL=(ALL) NOPASSWD:/usr/sbin/lsof`
  3. `Mine looks like: seantcanavan ALL=(ALL) NOPASSWD:/user/sbin/lsof`
5. Create `emaillogindetails.txt` inside `<clone_root_dir>/anon-eth-net/src/github.com/seantcanavan/assets/`.
6. Add your gmail address to line 1 and gmail password to line 2 of `emaillogindetails.txt`. Make sure that insecure app access is also enabled for the gmail account.
7. TBA


##Linux Setup:

####Required Packages:
1. Git: `sudo apt-get install git`
2. Go: `sudo apt-get install golang-go`
3. TBA - probably glide for better package management

####Required Setup:
1. Use `git clone` to download the repository.
2. Go is already automatically configured in your local Ubuntu path if installed via the synaptic package manager.
3. Setup GOBIN and GOPATH system variables for your Ubuntu user. Point to the root of the anon-eth-net clone directory. For GOBIN use the same value as GOPATH but add /bin to the end. GOPATH should be something like `/home/<username>/clone_root_dir/anon-eth-net`. You can automatically set the GOPATH variable in your .bashrc file on each terminal load / system startup:
  1. `nano ~/.bashrc`
  2. Scroll to the bottom and paste the following WITHOUT quotes. You may wrap the folder argument after the equal sign in quotes if you have spaces in your folder names (shame on you).
  3. `export GOPATH=/home/<username>/<clone_root_dir>/anon-eth-net`
  4. `export GOBIN=/home/<username>/<clone_root_dir>/anon-eth-net/bin`
  5. I use a folder called 'workspace' as a former Eclipse user so my paths looks like this:
  6. `export GOPATH=/home/seantcanavan/workspace/anon-eth-net`
  7. `export GOBIN=/home/seantcanavan/workspace/anon-eth-net/bin`
4. Update your sudoers file with the following:
  1. `sudo visudo`
  2. `<your-user-name-here> ALL=(ALL) NOPASSWD:/usr/bin/netstat`
  3. `Mine looks like: seantcanavan ALL=(ALL) NOPASSWD:/usr/bin/netstat`
5. Create `emaillogindetails.txt` inside `<clone_root_dir>/anon-eth-net/src/github.com/seantcanavan/assets/`.
6. Add your gmail address to line 1 and gmail password to line 2 of `emaillogindetails.txt`. Make sure that insecure app access is also enabled for the gmail account.
7. TBA

##Windows Setup:


####Required Packages:
1. Git: https://github.com/git-for-windows/git/releases/download/v2.10.2.windows.1/Git-2.10.2-64-bit.exe
2. Golang: https://storage.googleapis.com/golang/go1.7.3.windows-amd64.msi
3. TBA - probably glide for better package management

####Required Setup:
1. Use `git clone` to download the repository.
2. Configure your local system to make sure that the Go installation path is on the system path.
3. Setup GOBIN and GOPATH system variables for your windows users. Point to the root of the anon-eth-net clone directory. For GOBIN use the same value as GOPATH but add \bin to the end. GOPATH should be something like `C:\users\<username>\clone_root_dir\anon-eth-net`.
4. Create `emaillogindetails.txt` inside `<clone_root_dir>\anon-eth-net\src\github.com\seantcanavan\assets\`.
5. Add your gmail address to line 1 and gmail password to line 2 of `emaillogindetails.txt`. Make sure that insecure app access is also enabled for the gmail account.
6. TBA
7. TBA
