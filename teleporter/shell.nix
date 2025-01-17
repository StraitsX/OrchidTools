with import <nixpkgs> {};

pkgs.mkShell {
  name = "teleporter";
  buildInputs = [
    go_1_22
  ];

  shellHook = ''
    export APP_ENV="test"
    
    go mod download
  '';
}
