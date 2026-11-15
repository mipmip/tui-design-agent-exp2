{
  description = "tuimock - TUI mockup gallery for MCP agents";

  inputs.nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";

  outputs = { self, nixpkgs }:
    let
      systems = [
        "x86_64-linux"
        "aarch64-linux"
        "x86_64-darwin"
        "aarch64-darwin"
      ];
      forAllSystems = f: nixpkgs.lib.genAttrs systems (system: f (import nixpkgs { inherit system; }));
    in
    {
      packages = forAllSystems (pkgs: {
        default = pkgs.buildGoModule {
          pname = "tuimock";
          version = "0.1.0";
          src = ./.;
          vendorHash = "sha256-6d26F7aUEV6Tg8nF6FBEsEissUNEBdh0QatFCHPIiPE=";
          doCheck = true;
          meta = {
            description = "TUI mockup gallery for MCP agents";
            mainProgram = "tuimock";
            license = pkgs.lib.licenses.mit;
          };
        };
      });

      devShells = forAllSystems (pkgs: {
        default = pkgs.mkShell {
          packages = with pkgs; [
            go
            gopls
            golangci-lint
            jq
            jj
          ];
        };
      });

      checks = nixpkgs.lib.genAttrs systems (system: {
        default = self.packages.${system}.default;
      });

      apps = nixpkgs.lib.genAttrs systems (system: {
        default = {
          type = "app";
          program = "${self.packages.${system}.default}/bin/tuimock";
          meta.description = "TUI mockup gallery for MCP agents";
        };
      });
    };
}
