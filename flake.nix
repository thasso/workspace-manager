{
  description = "wsm - workspace manager for multi-repo projects";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = nixpkgs.legacyPackages.${system};

        version =
          if self ? rev
          then "git-${builtins.substring 0 7 self.rev}"
          else "dev-dirty";

        wsm = pkgs.buildGoModule {
          pname = "wsm";
          inherit version;
          src = self;
          vendorHash = "sha256-T8wu/E60gJ5a3RqSzKiBuvFG5RxizXsvRxUpLxlFO3A=";
          nativeCheckInputs = [ pkgs.git ];
          ldflags = [
            "-X github.com/thasso/wsm/internal/cli.Version=${version}"
          ];
          meta = {
            description = "Workspace manager for multi-repo projects";
            mainProgram = "wsm";
          };
        };
      in {
        packages = {
          default = wsm;
          wsm = wsm;
        };

        devShells.default = pkgs.mkShell {
          buildInputs = with pkgs; [
            go
            golangci-lint
            gopls
          ];
        };
      }
    );
}
