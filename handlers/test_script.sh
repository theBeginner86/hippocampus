#! /bin/bash

# a script that outputs following into database.aof file
# *3
# $3
# SET
# $5
# mykey
# $7
# myvalue

# loop 1000 times
for i in {1..10}
do
    echo "*3" >> database-10.aof
    echo "\$3" >> database-10.aof
    echo "SET" >> database-10.aof
    echo "\$5" >> database-10.aof
    echo "mykey" >> database-10.aof
    echo "\$7" >> database-10.aof
    echo "myvalue" >> database-10.aof
done