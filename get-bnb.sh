#!/bin/bash 

#proxy="%E2%9C%88%EF%B8%8F%20%E6%89%8B%E5%8A%A8%E5%88%87%E6%8D%A2"
#filters=(广港 广台 广日)

#proxy='sp-ss'
#filters=(A B)

proxy="%E8%AF%9D%E5%95%A6%E5%95%A6-%E7%BD%91%E7%BB%9C%E5%8A%A0%E9%80%9F%E4%BC%98%E8%B4%A8%E6%9C%8D%E5%8A%A1%E5%95%86"
filters=(华 A B C)

proxy_url="http://127.0.0.1:9090/proxies/"${proxy}

#filter="美"

for filter in ${filters[@]};
do 
for i in `curl -s -X GET $proxy_url | jq '.all | .[] | select(contains("'${filter}'"))' `; 
do 
  echo $i
  echo curl --noproxy -s -X PUT -H "content-type: application/json"  --data \{\"name\":${i}\} $proxy_url
  curl --noproxy -s -X PUT -H "content-type: application/json"  --data \{\"name\":${i}\} $proxy_url
  ./give-me-bnb -proxy socks5://127.0.0.1:7890 -to 0x14c75FC7aE1e566f57893435F34c7A488CBEf2e1
done 
done 





