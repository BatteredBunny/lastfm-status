{ buildGoModule }:
buildGoModule {
  src = ./.;

  name = "lastfm-status";
  vendorHash = "sha256-saRSF2cTV0VVgYJtUVBUKErCS3CaJCVbAsdSD7NM6pU=";

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
