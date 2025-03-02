{
  description = "Project flake";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs?ref=nixos-unstable";
  };

  outputs = {
    self,
    nixpkgs,
  }: let
    # System types to support.
    supportedSystems = ["x86_64-linux" "x86_64-darwin" "aarch64-linux" "aarch64-darwin"];

    # Helper function to generate an attrset '{ x86_64-linux = f "x86_64-linux"; ... }'.
    forAllSystems = nixpkgs.lib.genAttrs supportedSystems;

    # Nixpkgs instantiated for supported system types.
    nixpkgsFor = forAllSystems (system: import nixpkgs {inherit system;});
  in {
    packages = forAllSystems (system: let
      pkgs = nixpkgsFor.${system};
    in {
      default = pkgs.writeShellApplication {
        runtimeInputs = [
          # Frontend
          pkgs.elmPackages.elm
          pkgs.elmPackages.elm-live

          # Backend
          pkgs.go-task
          pkgs.go
          pkgs.air

          # Dev
          pkgs.process-compose
        ];

        text = ''
          process-compose -f ./process-compose.yaml
        '';
      };
    });

    devShells = forAllSystems (system: let
      pkgs = nixpkgsFor.${system};
    in {
      default = pkgs.mkShell {
        packages = [
          # Frontend
          pkgs.elmPackages.elm
          pkgs.elmPackages.elm-live

          # Backend
          pkgs.go-task
          pkgs.go
          pkgs.air

          # Dev
          pkgs.process-compose
        ];
      };
    });
  };
}
