{ pkgs ? import <nixpkgs> {} }:

pkgs.mkShell {
  buildInputs = [
    pkgs.go_1_21
  ];

  shellHook = ''
  '';

  MY_ENVIRONMENT_VARIABLE = "world";
}
