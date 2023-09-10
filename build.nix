{ pkgs, buildGoModule, lib }: buildGoModule rec {
    src = ./.;

    name = "lastfm-status";
    vendorSha256 = "sha256-HD5lkUt//BD2guQgw/Q9q3XRhEMflkunfSSuJIhReok=";

    ldflags = [
        "-s"
        "-w"
    ];
}