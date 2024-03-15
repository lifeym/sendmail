{ pkgs ? import <nixpkgs> { } }:
with pkgs;
mkShell {
  buildInputs = [
    go_1_21
    ko
    gnumake 
    go-task
    nil
    gopls
  ];
  
  GOPROXY = "https://goproxy.io,direct";
}
