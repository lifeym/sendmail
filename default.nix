{ pkgs ? import <nixpkgs> { } }:
with pkgs;
buildGo121Module {
  pname = "she";
  version = "0.1.0";

  src = ./.;
  vendorHash = "sha256-J/WOPdLykqxDa5GEs7rN/3DtDVzQfjch2tgIuu9xnWA=";

  # nativeBuildInputs = [
  # ];

  # env
  CGO_ENABLED = 0;

  # See https://github.com/NixOS/nixpkgs/blob/master/pkgs/build-support/go/module.nix
  # buildGoModule set GOPROXY to 'off',
  # so we should set it with hook if needed.
  preBuild  = ''
    export GOPROXY="https://goproxy.io,direct"
  '';
}

