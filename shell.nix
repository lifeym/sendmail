{ pkgs ? import <nixpkgs> { } }:
with pkgs;
mkShell {
  buildInputs = [
    go_1_22
    ko
    gnumake 
    go-task
    nil
    gopls
  ];
  
  GOPROXY = "https://goproxy.io,direct";
}
