for misskey only

https://moonchan.xyz/missakujo/

一時期のnotesを削除する (something like twitter delete)

used for delete notes in one period of time.

(fill the first blank and the second and the third will be filled automaticly.)

----

you can find your userId here
`https://[instanceDomain]/settings/account-info`

you should create api token first. \
Settings -> other settings -> API -> create token \
choose `writting or deleting notes` \
then create token, you will get your token (it's a string in random charaters)

----

~~if you don't know where to find your `userId`, please open this URL~~

~~`https://[instanceDomain]/.well-known/webfinger?resource=acct:[yourUsername]@[instanceDomain]`~~

~~you will find `userId` here~~

~~`application/activity+json","href":"https://[instanceDomain]/users/[userId]`~~


[Download .exe file](https://github.com/Hana-ame/missakujo/releases/tag/v0.0.0)

# test

```sh
go build .
./missakujo

cd frontend 
npm run dev
```


