{ buildGoLatestModule }:
buildGoLatestModule {
  src = ./.;

  name = "lastfm-status";
  vendorHash = "sha256-4+mlwx6z+4HT4aXO18BSwME4oNbUBIzD+s4cUvF7tMM=";

  ldflags = [
    "-s"
    "-w"
  ];

  meta = {
    description = "Simple api for showing currently scrobbling song on website ";
    homepage = "https://github.com/BatteredBunny/lastfm-status";
    mainProgram = "lastfm-status";
  };
}
