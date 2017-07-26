# IDEA

## How it works

### User get a token

User register an account and get a token
With a form, email verification
Free for some data
Pay plan for bigger

### User can send data

With this token, user can call api to send data
We store that data for time / size
Like :
```
curl -X POST -H "X-BigData4All-Token: Mysecret" \
  https://api.bigdata4a.ll/namespace/group/data1/type1 \
  -d '{"data": "Big Json data }'
```

### User can retrieve data

User can request stats/boards based on data he uploaded
