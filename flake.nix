{
  description = "sops development";

  inputs = {
    flake-utils = {
      url = "github:numtide/flake-utils";
    };

    nixpkgs = {
      url = "github:nixos/nixpkgs/nixpkgs-unstable";
    };
  };

  outputs =
    {
      self,
      nixpkgs,
      flake-utils,
    }:
    flake-utils.lib.eachDefaultSystem (
      system:
      let
        pkgs = import nixpkgs {
          inherit system;

          config = {
            allowUnfree = true;
          };
        };
      in
      {
        devShells = {
          default = pkgs.mkShell {
            # https://nixos.wiki/wiki/Go#Using_cgo_on_NixOS
            hardeningDisable = [ "fortify" ];

            packages = with pkgs; [
              delve
              go
            ];

            shellHook = ''
							CGO_ENABLED = 0;
						'';
          };
        };

        packages = {
          default = pkgs.buildGoModule {
            pname = "sops";
            version = "3.13.1";

            src = ./.;

            subPackages = [ "cmd/sops" ];

            vendorHash = "sha256-b94pcUopemj+kXj2AacZTQ0BYaTMqXAgHUEVz6x3+Lg=";

            meta = {
              mainProgram = "sops";
            };
          };
        };

        apps = {
          default = {
            type = "app";
            program = "${self.packages.${system}.default}/bin/sops";
          };
        };
      }
    );
}