{ pkgs ? import <nixpkgs> {} }:

pkgs.mkShell {
  buildInputs = [
    pkgs.go_1_22
  ];

  shellHook = ''
  '';

  MY_ENVIRONMENT_VARIABLE = "world";
}
