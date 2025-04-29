{ buildGoModule }:
buildGoModule {
  src = ./.;

  name = "lastfm-status";
  vendorHash = "sha256-5+NtQzHHQGFDXKzLmQnE6On4mFQRQhoyIH4CyE8XsO8=";

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
