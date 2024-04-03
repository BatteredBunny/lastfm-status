{
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = {
    self,
    nixpkgs,
    flake-utils,
    ...
  } @ inputs:
    flake-utils.lib.eachDefaultSystem (
      system: let
        pkgs = import nixpkgs {
          inherit system;
        };
      in
        with pkgs; {
          nixosModules.default = import ./module.nix inputs;

          devShells.default = mkShell {
            buildInputs = with pkgs; [
              go
            ];
          };

          packages.default = pkgs.callPackage ./build.nix {};
        }
    );
}
