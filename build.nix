{buildGoModule}:
buildGoModule {
  src = ./.;

  name = "lastfm-status";
  vendorHash = "sha256-OfHiQkZBQDkXajBKw/22rB6WZINxe+Jr7qrNsQukp/E=";

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
