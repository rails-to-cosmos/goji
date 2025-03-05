{
  description = "Go Development Environment with Nix Flakes";

  inputs.nixpkgs.url = "github:NixOS/nixpkgs";

  outputs = { nixpkgs, ... }:
  let
    system = "x86_64-linux";
    pkgs = nixpkgs.legacyPackages.${system};
  in {
    devShells.${system}.default = pkgs.mkShell {
      buildInputs = [
        pkgs.go
        pkgs.gopls
        pkgs.gotools
        pkgs.gofumpt
        pkgs.golangci-lint
        pkgs.delve
      ];

      shellHook = ''
        export GOPATH="$PWD/go"
        export PATH="$GOPATH/bin:$PATH"
        echo "Go development environment is ready!"
        go version
      '';
    };
  };
}
