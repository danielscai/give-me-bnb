
import requests

import json 
import subprocess

from urllib.parse import quote as url_encode


#proxy_group_name = "sp-ss"
#proxy_group_name = "话啦啦-网络加速优质服务商"
proxy_group_name = "✈️ 手动切换"
address = "0x8D45fe25F186F2362766B58cb38064237BfdC47c"


give_me_bnb = f"./give-me-bnb -proxy socks5://127.0.0.1:7890 -to {address}"
clash_api_proxies="http://127.0.0.1:9090/proxies"
max_delay = 1000

# https://clash.gitbook.io/doc/restful-api/proxies

proxies = {
  "http": None,
  "https": None,
}

proxy_group_name_encode=url_encode(proxy_group_name)
proxy_url = f"{clash_api_proxies}/{proxy_group_name_encode}"

proxy_result = requests.get(proxy_url,proxies=proxies)
all_proxies = proxy_result.json()['all']


for proxy in all_proxies:
    proxy_encode = url_encode(proxy)
    
    clash_api_proxy_info_url = f"{clash_api_proxies}/{proxy_encode}"
    
    proxy_delay = requests.get(clash_api_proxy_info_url,proxies=proxies)
    if proxy_delay.status_code == 200:
        # print(proxy,clash_api_proxy_info_url)
        # print(proxy_delay.json())
        delay_json = proxy_delay.json()['history']
        if len(delay_json) == 0:
            print(f'[WARN] skip empty proxy {proxy}')
        else:
            delay = delay_json[0]['delay']
            if delay == 0:
                print(f'[WARN] skip delay 0 {proxy}')
                continue
            
            if delay >= max_delay:
                print(f'[WARN] skip {proxy} delay geater than {max_delay}')
                continue
            
            # print(delay,proxy,proxy_url)
            request_data = {"name":proxy} 
            # print(request_data)
            req = requests.put(proxy_url,data=json.dumps(request_data),proxies=proxies)
            
            if req.status_code == 204:
                print(f'process {proxy}',end=' ',flush=True)
                subprocess.call(give_me_bnb,shell=True)
            else:
                print('[ERROR] swith clash proxy failed')
