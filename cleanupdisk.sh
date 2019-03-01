#!/bin/bash
time=$(date "+%Y-%m-%d %H:%M:%S")
workdir=("/boot/logs" "/home/logs" "/home/test")
  	for wdir in ${workdir[@]}
 	 do
          	echo  filepath is $wdir
          	if [ $wdir =  ${workdir[0]} ] ;then
             		fileStr=`find  $wdir -type d`
             		echo dir is $fileStr
         	 else
             		fileStr=`find  $wdir -type d`
             		echo dir is $fileStr
          	fi
         	for dir in $fileStr
          	do
                	echo dir name is $dir         
                	find $dir -name "*.png" -and -mtime +1  -type f | xargs rm	
               		find $dir -name "*.jpg" -and -mtime +1  -type f | xargs rm 
                        if [ $? -eq 0 ];then
                        	echo "${time}" delete $dir success!             
               		 else
                        	echo "${time}" delete $dir FAILD!            
               		 fi
          	done
  	done
