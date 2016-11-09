## Anonymous Ethereum Network (anon-eth-net)

Totally anonymous botnet client with an emphasis on individual zombie control, resiliency of the host machine, and ease of remote code execution. Zombies mine ethereum in their free time for fun and check in at pre-specified intervals to report on their CPU, memory, disk, and network utilization. Can be combined with go-dos-yourself to enable remote network performence testing, fuzzing, spoofing, and attacking. Use responsibly.

###Currently supported platforms:
- macOS El Capitan 10.11.6
- Ubuntu 16.04.1

###Supported platforms wishlist:
- macOS >= 10.8.X
- Ubuntu >= 14.04.5  
- Windows >= 7 SP2

##Mac Setup:

####Required Packages:
1. TBA
2. TBA
3. TBA

####Required Setup:
1. Update your sudoers file with the following: <your-user-name-here> ALL=(ALL) NOPASSWD:/usr/sbin/lsof
2. TBA
3. TBA


##Linux Setup:

####Required Packages:
1. TBA
2. TBA
3. TBA

####Required Setup:
1. Update your sudoers file with the following: <your-user-name-here> ALL=(ALL) NOPASSWD:/usr/bin/netstat
2. TBA
3. TBA


##Windows Setup:


####Required Packages:
1. Git: https://github.com/git-for-windows/git/releases/download/v2.10.2.windows.1/Git-2.10.2-64-bit.exe
2. Golang: https://storage.googleapis.com/golang/go1.7.3.windows-amd64.msi

####Required Setup:
1. Use 'git clone' to download the repository
2. Install and setup Go on your local system path
3. Setup GOBIN and GOPATH system variables for your windows users. Point to the root of the anon-eth-net clone directory. For GOBIN use the same value as GOPATH but add \bin to the end
4. Create 'emaillogindetails.txt' inside src\github.com\seantcanavan\config\emaillogindetails.txt
5. Add your gmail address to line 1 and gmail password to line 2 of emaillogindetails.txt. Make sure that insecure app access is also enabled for the gmail account.
