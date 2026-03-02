{ buildGoModule }:
buildGoModule {
  src = ./.;

  name = "lastfm-status";
  vendorHash = "sha256-lu5hnkl1PJxyUq+aGdqyg5qd3CBAOC+eJ+EQy4rgxJs=";

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
