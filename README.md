# lastfm-status

Simple api for showing currently scrobbling song on website

![Example iframe usage](example.png)

## Embedding it into your site

```html
<iframe
    src="https://lastfm.catnip.ee/status?username=YourLastfmUsername&theme=dark"
></iframe>
```


## Dependencies
- golang

## Running manually
```
git clone https://github.com/BatteredBunny/lastfm-status
cd lastfm-status
go run . --port 8080
```

## Running with nix
```
nix run github:BatteredBunny/lastfm-status
```

## Running as service on nixos
```nix
# flake.nix
inputs = {
    lastfm-status.url = "github:BatteredBunny/lastfm-status";
};
```

```nix
# configuration.nix
imports = [
    inputs.lastfm-status.nixosModules.default
];

services.lastfm-status = {
    enable = true;
    port = 8080;

    # Optional parameters
    package = inputs.lastfm-status.packages.${builtins.currentSystem}.default;
    cacheLength = "1m";
    monthlyCacheLength = "1h";
    
    enableRatelimiting = true;
    reverseProxy = false;
    trustedProxy = "";
};
```