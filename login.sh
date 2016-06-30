#!/usr/bin/expect

send_user "docker login index.tenxcloud.com\n"
 
set username $env(USERNAME)
set password $env(PASSWORD)
set srcRepo $env(SRC_REPOSITORY)
set dstRepo $env(DST_REPOSITORY)

send_user "username: $username, password: $password, src: $srcRepo, dst: $dstRepo"

spawn docker login index.tenxcloud.com
expect {
        "Username:" {
                send "$username\n"
                exp_continue
        }
        "Password:" {
                send "$password\n"
                exp_continue
        }
        "Email:" {
                send "$username@tenxcloud.com\n"
                exp_continue
        }
        eof {
                send_user "\n"
        }
}