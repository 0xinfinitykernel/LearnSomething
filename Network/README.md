### OPEN FTP
```
ftp server enable
aaa
local-user admin service-type telnet terminal ssh ftp http
local-user admin ftp-directory flash:/
ftp server-source  all-interface
return
```

### ATX-1800
dd if=/dev/mtd8 of=/tmp/art_bak bs=128k count=1 2>/dev/null
echo "COUNTRY:US" |dd of=/tmp/art_bak bs=1 count=10 seek=$((0x90)) conv=notrunc 2>/dev/null
mtd write /tmp/art_bak /dev/mtd8 2>/dev/null

### SET IMEI
```
# old
AT+EGMR=1,7,"868371051160162"

# new
AT+EGMR=1,7,"865025071687702"

# reboot
AT+CFUN=1,1
```
