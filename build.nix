{ pkgs, buildGoModule, lib }: buildGoModule rec {
    src = ./.;

    name = "lastfm-status";
    vendorSha256 = "sha256-HD5lkUt//BD2guQgw/Q9q3XRhEMflkunfSSuJIhReok=";

    ldflags = [
        "-s"
        "-w"
    ];

    meta = {
        description = "Simple api for showing currently scrobbling song on website ";
        homepage = "https://github.com/ayes-web/lastfm-status¬";
        mainProgram = "lastfm-status";
    };
}