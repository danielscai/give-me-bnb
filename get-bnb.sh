#!/bin/bash 

proxy_url="http://127.0.0.1:9090/proxies/%E2%9C%88%EF%B8%8F%20%E6%89%8B%E5%8A%A8%E5%88%87%E6%8D%A2"

#filter="美"
array=(广港 广台 广日)

for filter in ${array[@]};
do 
for i in `curl -s -X GET $proxy_url | jq '.all | .[] | select(contains("'${filter}'"))' `; 
do 
  echo $i
  curl --noproxy -s -X PUT -H "content-type: application/json"  --data \{\"name\":${i}\} $proxy_url
  ./give-me-bnb -proxy socks5://127.0.0.1:7890 -to 0x14c75FC7aE1e566f57893435F34c7A488CBEf2e1
done 
done 





