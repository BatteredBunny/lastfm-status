# lastfm-status

Simple api for showing currently scrobbling song on website

![Example iframe usage](example.png)

## Iframe usage

```html
<iframe
    src="http://localhost:8080/status?username=YourLastfmUsername"
></iframe>
```


## Dependencies
- golang

## Running manually
```
git clone https://github.com/ayes-web/lastfm-status
cd lastfm-status
go run . --port 8080
```

## Running with nix
```
nix run github:ayes-web/lastfm-status
```