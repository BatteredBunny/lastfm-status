{ buildGoModule }:
buildGoModule {
  src = ./.;

  name = "lastfm-status";
  vendorHash = "sha256-oloJwUqtvpilMrSq/cLEyH1RV1ggvNozBENWuvGJGts=";

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
