 #!/bin/sh

 FILE=./data.txt

 while read LINE
 do
     mysql --defaults-extra-file=./access.cnf -D test -e "select table_name from information_schema.tables;" >> "./data.txt"
 done < $FILE

 echo "done"
