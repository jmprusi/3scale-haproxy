# 3scale-haproxy

PoC: 3scale HAProxy SPOA for authentication


## How to run it? 

Start the agent:

```
go mod vendor
go run main.go
```

Install HAProxy: 

```
brew install haproxy
```

Configure your credentials:
```
vi examples/haproxy.conf

(..)
    tcp-request content set-var(req.service_id) str("1234567")    # <------ Replace with your serviceID
    tcp-request content set-var(req.access_token) str("XXXXXXXXXXXXXXXXXXXXXXXXXXXX")  # <----------- Replace with your accesstoken
    tcp-request content set-var(req.system_url) str("https://yourtenant-admin.3scale.net")  # <-------- Replace with your admin portal direction

    filter spoe engine 3scale-api-authorization config spoe-api-authorization.conf

    http-request deny deny_status 403 unless { var(txn.3scale.authorized) -m int eq 1 }

    http-request set-header Host echo-api.3scale.net   #  <--------- Put your API Backend host here... 
    server echo_api1 echo-api.3scale.net:80 check      #  <--------- And here.
(..)
```

Run haproxy! 
```
cd examples
haproxy -f haproxy.conf
```

Test it:
```
curl http://localhost/  <--- Should fail..
curl http://localhost/?user_key=XXXXXXXXXXXXXXX <------- Test with the proper user_key.
```

