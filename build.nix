{ pkgs, buildGoModule, lib }: buildGoModule rec {
    src = ./.;

    name = "lastfm-status";
    vendorSha256 = "sha256-UvQ+Q7op30ox34rO2BXcCE6vLDFsF27pVytMAVpCi2U=";

    ldflags = [
        "-s"
        "-w"
    ];
}