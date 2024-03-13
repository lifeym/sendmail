{ pkgs ? import <nixpkgs> { } }:
with pkgs;
mkShell rec {
  buildInputs = [
    go_1_21
    ko
    gnumake 
    go-task
  ];
  
  GOPROXY = "https://goproxy.io,direct";
}
